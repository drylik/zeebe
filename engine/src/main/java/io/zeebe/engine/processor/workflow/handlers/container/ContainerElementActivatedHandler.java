/*
 * Copyright Camunda Services GmbH and/or licensed to Camunda Services GmbH under
 * one or more contributor license agreements. See the NOTICE file distributed
 * with this work for additional information regarding copyright ownership.
 * Licensed under the Zeebe Community License 1.0. You may not use this file
 * except in compliance with the Zeebe Community License 1.0.
 */
package io.zeebe.engine.processor.workflow.handlers.container;

import io.zeebe.engine.processor.workflow.BpmnStepContext;
import io.zeebe.engine.processor.workflow.deployment.model.element.ExecutableCatchEventElement;
import io.zeebe.engine.processor.workflow.deployment.model.element.ExecutableFlowElementContainer;
import io.zeebe.engine.processor.workflow.deployment.model.element.ExecutableStartEvent;
import io.zeebe.engine.processor.workflow.handlers.element.ElementActivatedHandler;
import io.zeebe.engine.state.deployment.WorkflowState;
import io.zeebe.engine.state.instance.IndexedRecord;
import io.zeebe.engine.state.instance.StoredRecord.Purpose;
import io.zeebe.protocol.impl.record.value.workflowinstance.WorkflowInstanceRecord;
import io.zeebe.protocol.record.intent.WorkflowInstanceIntent;
import java.util.List;

public class ContainerElementActivatedHandler<T extends ExecutableFlowElementContainer>
    extends ElementActivatedHandler<T> {
  private final WorkflowState workflowState;

  public ContainerElementActivatedHandler(final WorkflowState workflowState) {
    this(null, workflowState);
  }

  public ContainerElementActivatedHandler(
      final WorkflowInstanceIntent nextState, final WorkflowState workflowState) {
    super(nextState);
    this.workflowState = workflowState;
  }

  @Override
  protected boolean handleState(final BpmnStepContext<T> context) {
    if (!super.handleState(context)) {
      return false;
    }

    final ExecutableFlowElementContainer element = context.getElement();
    final ExecutableStartEvent firstStartEvent = element.getStartEvents().get(0);

    // workflows with none start event only have a single none start event and no other types of
    // start events; note that sub-processes only have a single none start event, so
    // publishing a deferred record only applies to processes
    if (firstStartEvent.isNone() || firstStartEvent.getEventSubProcess() != null) {
      activateNoneStartEvent(context, firstStartEvent);
    } else {
      publishDeferredRecord(context);
    }

    context
        .getStateDb()
        .getElementInstanceState()
        .spawnToken(context.getElementInstance().getKey());
    return true;
  }

  private void publishDeferredRecord(final BpmnStepContext<T> context) {
    final IndexedRecord deferredRecord = getDeferredRecord(context);
    context
        .getOutput()
        .appendFollowUpEvent(
            deferredRecord.getKey(), deferredRecord.getState(), deferredRecord.getValue());
  }

  private void activateNoneStartEvent(
      final BpmnStepContext<T> context, final ExecutableCatchEventElement firstStartEvent) {
    final WorkflowInstanceRecord value = context.getValue();

    value.setElementId(firstStartEvent.getId());
    value.setBpmnElementType(firstStartEvent.getElementType());
    value.setFlowScopeKey(context.getKey());
    context.getOutput().appendNewEvent(WorkflowInstanceIntent.ELEMENT_ACTIVATING, value);
  }

  private IndexedRecord getDeferredRecord(final BpmnStepContext<T> context) {
    final long scopeKey = context.getKey();
    final List<IndexedRecord> deferredRecords =
        context.getElementInstanceState().getDeferredRecords(scopeKey);

    if (deferredRecords.isEmpty()) {
      throw new IllegalStateException(
          "Expected process with no none start events to have a deferred record, but nothing was found");
    }

    assert deferredRecords.size() == 1
        : "should only have one deferred start event per workflow instance";

    final IndexedRecord deferredRecord = deferredRecords.get(0);
    workflowState
        .getElementInstanceState()
        .removeStoredRecord(scopeKey, deferredRecord.getKey(), Purpose.DEFERRED);
    return deferredRecord;
  }
}
