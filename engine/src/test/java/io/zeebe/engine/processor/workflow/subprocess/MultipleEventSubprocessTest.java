/*
 * Copyright Camunda Services GmbH and/or licensed to Camunda Services GmbH under
 * one or more contributor license agreements. See the NOTICE file distributed
 * with this work for additional information regarding copyright ownership.
 * Licensed under the Zeebe Community License 1.0. You may not use this file
 * except in compliance with the Zeebe Community License 1.0.
 */
package io.zeebe.engine.processor.workflow.subprocess;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.groups.Tuple.tuple;

import io.zeebe.engine.util.EngineRule;
import io.zeebe.model.bpmn.Bpmn;
import io.zeebe.model.bpmn.BpmnModelInstance;
import io.zeebe.model.bpmn.builder.EmbeddedSubProcessBuilder;
import io.zeebe.model.bpmn.builder.ProcessBuilder;
import io.zeebe.model.bpmn.builder.StartEventBuilder;
import io.zeebe.protocol.record.intent.JobIntent;
import io.zeebe.protocol.record.intent.MessageSubscriptionIntent;
import io.zeebe.protocol.record.intent.WorkflowInstanceIntent;
import io.zeebe.protocol.record.value.BpmnElementType;
import io.zeebe.test.util.BrokerClassRuleHelper;
import io.zeebe.test.util.record.RecordingExporter;
import java.time.Duration;
import org.junit.ClassRule;
import org.junit.Rule;
import org.junit.Test;

public class MultipleEventSubprocessTest {

  @ClassRule public static final EngineRule ENGINE = EngineRule.singlePartition();

  @Rule public final BrokerClassRuleHelper helper = new BrokerClassRuleHelper();

  @Test
  public void shouldTriggerMultipleEventSubprocesses() {
    final BpmnModelInstance model = twoEventSubprocModel(false, false, helper.getMessageName());
    ENGINE.deployment().withXmlResource(model).deploy();

    // when
    final long wfInstanceKey =
        ENGINE.workflowInstance().ofBpmnProcessId("proc").withVariable("key", "123").create();
    triggerTimerStart(wfInstanceKey);
    triggerMessageStart(wfInstanceKey, helper.getMessageName());

    // then
    assertThat(
            RecordingExporter.workflowInstanceRecords()
                .withWorkflowInstanceKey(wfInstanceKey)
                .withElementType(BpmnElementType.SUB_PROCESS)
                .limit(8))
        .extracting(r -> tuple(r.getValue().getElementId(), r.getIntent()))
        .containsSubsequence(
            tuple("event_sub_proc_timer", WorkflowInstanceIntent.ELEMENT_ACTIVATED),
            tuple("event_sub_proc_timer", WorkflowInstanceIntent.ELEMENT_COMPLETED))
        .containsSubsequence(
            tuple("event_sub_proc_msg", WorkflowInstanceIntent.ELEMENT_ACTIVATED),
            tuple("event_sub_proc_msg", WorkflowInstanceIntent.ELEMENT_COMPLETED));
  }

  @Test
  public void shouldInterruptOtherActiveEventSubprocess() {
    // given
    final BpmnModelInstance model =
        twoEventSubprocWithTasksModel(false, true, helper.getMessageName());
    ENGINE.deployment().withXmlResource(model).deploy();

    // when
    final long wfInstanceKey =
        ENGINE.workflowInstance().ofBpmnProcessId("proc").withVariable("key", "123").create();

    triggerTimerStart(wfInstanceKey);
    assertThat(
            RecordingExporter.workflowInstanceRecords(WorkflowInstanceIntent.ELEMENT_ACTIVATED)
                .withWorkflowInstanceKey(wfInstanceKey)
                .withElementId("event_sub_task_timer")
                .exists())
        .describedAs("Expected service task after timer start event to exist")
        .isTrue();
    triggerMessageStart(wfInstanceKey, helper.getMessageName());

    // then
    assertThat(
            RecordingExporter.workflowInstanceRecords()
                .withWorkflowInstanceKey(wfInstanceKey)
                .limitToWorkflowInstanceCompleted())
        .extracting(r -> tuple(r.getValue().getElementId(), r.getIntent()))
        .containsSubsequence(
            tuple("event_sub_proc_timer", WorkflowInstanceIntent.ELEMENT_ACTIVATED),
            tuple("event_sub_task_timer", WorkflowInstanceIntent.ELEMENT_ACTIVATED),
            tuple("event_sub_start_msg", WorkflowInstanceIntent.EVENT_OCCURRED),
            tuple("event_sub_proc_timer", WorkflowInstanceIntent.ELEMENT_TERMINATING),
            tuple("event_sub_task_timer", WorkflowInstanceIntent.ELEMENT_TERMINATED),
            tuple("event_sub_proc_timer", WorkflowInstanceIntent.ELEMENT_TERMINATED),
            tuple("event_sub_proc_msg", WorkflowInstanceIntent.ELEMENT_ACTIVATED),
            tuple("event_sub_proc_msg", WorkflowInstanceIntent.ELEMENT_COMPLETED),
            tuple("proc", WorkflowInstanceIntent.ELEMENT_COMPLETED));
  }

  @Test
  public void shouldCloseEventSubscriptionWhenScopeCloses() {
    final BpmnModelInstance model = nestedMsgModel(helper.getMessageName());
    ENGINE.deployment().withXmlResource(model).deploy();

    // when
    final long wfInstanceKey =
        ENGINE.workflowInstance().ofBpmnProcessId("proc").withVariable("key", "123").create();

    assertThat(
            RecordingExporter.messageSubscriptionRecords(MessageSubscriptionIntent.OPENED)
                .withWorkflowInstanceKey(wfInstanceKey)
                .exists())
        .describedAs("Expected event subprocess message start subscription to be opened.")
        .isTrue();
    completeJob(wfInstanceKey, "sub_proc_type");

    // then
    assertThat(
            RecordingExporter.messageSubscriptionRecords(MessageSubscriptionIntent.CLOSED)
                .withWorkflowInstanceKey(wfInstanceKey)
                .withMessageName(helper.getMessageName())
                .exists())
        .describedAs("Expected event subprocess start message subscription to be closed.")
        .isTrue();
  }

  @Test
  public void shouldCorrelateTwoMessagesIfNonInterrupting() {
    // given
    final BpmnModelInstance model =
        twoEventSubprocWithTasksModel(false, false, helper.getMessageName());
    ENGINE.deployment().withXmlResource(model).deploy();

    // when
    final long wfInstanceKey =
        ENGINE.workflowInstance().ofBpmnProcessId("proc").withVariable("key", "123").create();
    triggerMessageStart(wfInstanceKey, helper.getMessageName());
    triggerMessageStart(wfInstanceKey, helper.getMessageName());

    // then
    assertThat(
            RecordingExporter.messageSubscriptionRecords(MessageSubscriptionIntent.CORRELATED)
                .withWorkflowInstanceKey(wfInstanceKey)
                .limit(2)
                .count())
        .isEqualTo(2);
  }

  @Test
  public void shouldKeepWorkflowInstanceActive() {
    // given
    final BpmnModelInstance model =
        twoEventSubprocWithTasksModel(false, false, helper.getMessageName());
    ENGINE.deployment().withXmlResource(model).deploy();

    final long wfInstanceKey =
        ENGINE.workflowInstance().ofBpmnProcessId("proc").withVariable("key", "123").create();
    triggerTimerStart(wfInstanceKey);

    RecordingExporter.jobRecords(JobIntent.CREATED)
        .withWorkflowInstanceKey(wfInstanceKey)
        .withType("timerTask")
        .await();

    // when
    completeJob(wfInstanceKey, "type");

    // then
    RecordingExporter.workflowInstanceRecords(WorkflowInstanceIntent.ELEMENT_COMPLETED)
        .withWorkflowInstanceKey(wfInstanceKey)
        .withElementType(BpmnElementType.END_EVENT)
        .withElementId("end_proc")
        .await();

    completeJob(wfInstanceKey, "timerTask");
    assertThat(RecordingExporter.workflowInstanceRecords().limitToWorkflowInstanceCompleted())
        .extracting(r -> tuple(r.getValue().getElementId(), r.getIntent()))
        .containsSubsequence(
            tuple("event_sub_task_timer", WorkflowInstanceIntent.ELEMENT_ACTIVATED),
            tuple("end_proc", WorkflowInstanceIntent.ELEMENT_COMPLETED),
            tuple("event_sub_task_timer", WorkflowInstanceIntent.ELEMENT_COMPLETED),
            tuple("proc", WorkflowInstanceIntent.ELEMENT_COMPLETED));
  }

  @Test
  public void shouldTerminateEventSubprocessIfScopeTerminates() {
    // given
    final BpmnModelInstance model =
        twoEventSubprocWithTasksModel(false, false, helper.getMessageName());
    ENGINE.deployment().withXmlResource(model).deploy();

    final long wfInstanceKey =
        ENGINE.workflowInstance().ofBpmnProcessId("proc").withVariable("key", "123").create();
    triggerTimerStart(wfInstanceKey);

    RecordingExporter.jobRecords(JobIntent.CREATED)
        .withWorkflowInstanceKey(wfInstanceKey)
        .withType("timerTask")
        .await();

    // when
    ENGINE.workflowInstance().withInstanceKey(wfInstanceKey).cancel();

    // then
    assertThat(RecordingExporter.workflowInstanceRecords().limitToWorkflowInstanceTerminated())
        .extracting(r -> tuple(r.getValue().getBpmnElementType(), r.getIntent()))
        .containsSubsequence(
            tuple(BpmnElementType.SERVICE_TASK, WorkflowInstanceIntent.ELEMENT_ACTIVATED),
            tuple(BpmnElementType.PROCESS, WorkflowInstanceIntent.ELEMENT_TERMINATING),
            tuple(BpmnElementType.SUB_PROCESS, WorkflowInstanceIntent.ELEMENT_TERMINATING),
            tuple(BpmnElementType.SERVICE_TASK, WorkflowInstanceIntent.ELEMENT_TERMINATING),
            tuple(BpmnElementType.SERVICE_TASK, WorkflowInstanceIntent.ELEMENT_TERMINATED),
            tuple(BpmnElementType.SUB_PROCESS, WorkflowInstanceIntent.ELEMENT_TERMINATED),
            tuple(BpmnElementType.PROCESS, WorkflowInstanceIntent.ELEMENT_TERMINATED));
  }

  private void triggerMessageStart(long wfInstanceKey, String msgName) {
    RecordingExporter.messageSubscriptionRecords(MessageSubscriptionIntent.OPENED)
        .withWorkflowInstanceKey(wfInstanceKey)
        .await();

    ENGINE.message().withName(msgName).withCorrelationKey("123").publish();
  }

  private void triggerTimerStart(long wfInstanceKey) {
    RecordingExporter.workflowInstanceRecords(WorkflowInstanceIntent.ELEMENT_ACTIVATED)
        .withWorkflowInstanceKey(wfInstanceKey)
        .withElementType(BpmnElementType.SERVICE_TASK)
        .await();

    ENGINE.increaseTime(Duration.ofSeconds(60));
  }

  private static void completeJob(long wfInstanceKey, String taskType) {
    RecordingExporter.jobRecords(JobIntent.CREATED)
        .withWorkflowInstanceKey(wfInstanceKey)
        .withType(taskType)
        .await();

    ENGINE.job().ofInstance(wfInstanceKey).withType(taskType).complete();
  }

  private BpmnModelInstance twoEventSubprocModel(
      boolean timerInterrupt, boolean msgInterrupt, String msgName) {
    final ProcessBuilder builder = Bpmn.createExecutableProcess("proc");

    builder
        .eventSubProcess("event_sub_proc_timer")
        .startEvent("event_sub_start_timer")
        .interrupting(timerInterrupt)
        .timerWithDuration("PT60S")
        .endEvent("event_sub_end_timer");
    builder
        .eventSubProcess("event_sub_proc_msg")
        .startEvent("event_sub_start_msg")
        .interrupting(msgInterrupt)
        .message(b -> b.name(msgName).zeebeCorrelationKey("key"))
        .endEvent("event_sub_end_msg");

    return builder
        .startEvent("start_proc")
        .serviceTask("task", t -> t.zeebeTaskType("type"))
        .endEvent("end_proc")
        .done();
  }

  private BpmnModelInstance twoEventSubprocWithTasksModel(
      boolean timerInterrupt, boolean msgInterrupt, String msgName) {
    final ProcessBuilder builder = Bpmn.createExecutableProcess("proc");

    builder
        .eventSubProcess("event_sub_proc_timer")
        .startEvent("event_sub_start_timer")
        .interrupting(timerInterrupt)
        .timerWithDuration("PT60S")
        .serviceTask("event_sub_task_timer", b -> b.zeebeTaskType("timerTask"))
        .endEvent("event_sub_end_timer");
    builder
        .eventSubProcess("event_sub_proc_msg")
        .startEvent("event_sub_start_msg")
        .interrupting(msgInterrupt)
        .message(b -> b.name(msgName).zeebeCorrelationKey("key"))
        .endEvent("event_sub_end_msg");

    return builder
        .startEvent("start_proc")
        .serviceTask("task", t -> t.zeebeTaskType("type"))
        .endEvent("end_proc")
        .done();
  }

  private static BpmnModelInstance nestedMsgModel(String msgName) {
    final StartEventBuilder procBuilder =
        Bpmn.createExecutableProcess("proc").startEvent("proc_start");
    procBuilder.serviceTask("proc_task", b -> b.zeebeTaskType("proc_type")).endEvent();
    final EmbeddedSubProcessBuilder subProcBuilder =
        procBuilder.subProcess("sub_proc").embeddedSubProcess();

    subProcBuilder
        .eventSubProcess("event_sub_proc")
        .startEvent("event_sub_start")
        .interrupting(true)
        .message(b -> b.name(msgName).zeebeCorrelationKey("key"))
        .endEvent("event_sub_end");
    return subProcBuilder
        .startEvent("sub_start")
        .serviceTask("sub_proc_task", t -> t.zeebeTaskType("sub_proc_type"))
        .endEvent("sub_end")
        .done();
  }
}
