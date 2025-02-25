// Copyright © 2018 Camunda Services GmbH (info@camunda.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/zeebe-io/zeebe/clients/go/mock_pb"
	"github.com/zeebe-io/zeebe/clients/go/pb"
	"github.com/zeebe-io/zeebe/clients/go/utils"
	"testing"
)

type DataType struct {
	Foo string `json:"foo,omitempty"`
}

func (cmd DataType) String() string {
	return fmt.Sprintf("{\"foo\":\"%s\"}", cmd.Foo)
}

func TestCreateWorkflowInstanceCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	request := &pb.CreateWorkflowInstanceRequest{
		WorkflowKey: 123,
	}
	stub := &pb.CreateWorkflowInstanceResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
	}

	client.EXPECT().CreateWorkflowInstance(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	response, err := command.WorkflowKey(123).Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceCommandByBpmnProcessId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	request := &pb.CreateWorkflowInstanceRequest{
		BpmnProcessId: "foo",
		Version:       LatestVersion,
	}
	stub := &pb.CreateWorkflowInstanceResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
	}

	client.EXPECT().CreateWorkflowInstance(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	response, err := command.BPMNProcessId("foo").LatestVersion().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceCommandByBpmnProcessIdAndVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	request := &pb.CreateWorkflowInstanceRequest{
		BpmnProcessId: "foo",
		Version:       56,
	}
	stub := &pb.CreateWorkflowInstanceResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
	}

	client.EXPECT().CreateWorkflowInstance(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	response, err := command.BPMNProcessId("foo").Version(56).Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceCommandWithVariablesFromString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"bar\"}"

	request := &pb.CreateWorkflowInstanceRequest{
		WorkflowKey: 123,
		Variables:   variables,
	}
	stub := &pb.CreateWorkflowInstanceResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
	}

	client.EXPECT().CreateWorkflowInstance(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromString(variables)
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceCommandWithVariablesFromStringer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"bar\"}"

	request := &pb.CreateWorkflowInstanceRequest{
		WorkflowKey: 123,
		Variables:   variables,
	}
	stub := &pb.CreateWorkflowInstanceResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
	}

	client.EXPECT().CreateWorkflowInstance(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromStringer(DataType{Foo: "bar"})
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceCommandWithVariablesFromObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"bar\"}"

	request := &pb.CreateWorkflowInstanceRequest{
		WorkflowKey: 123,
		Variables:   variables,
	}
	stub := &pb.CreateWorkflowInstanceResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
	}

	client.EXPECT().CreateWorkflowInstance(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromObject(DataType{Foo: "bar"})
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceCommandWithVariablesFromObjectOmitempty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{}"

	request := &pb.CreateWorkflowInstanceRequest{
		WorkflowKey: 123,
		Variables:   variables,
	}
	stub := &pb.CreateWorkflowInstanceResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
	}

	client.EXPECT().CreateWorkflowInstance(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromObject(DataType{Foo: ""})
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceCommandWithVariablesFromObjectIgnoreOmitempty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"\"}"

	request := &pb.CreateWorkflowInstanceRequest{
		WorkflowKey: 123,
		Variables:   variables,
	}
	stub := &pb.CreateWorkflowInstanceResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
	}

	client.EXPECT().CreateWorkflowInstance(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromObjectIgnoreOmitempty(DataType{Foo: ""})
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceCommandWithVariablesFromMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"bar\"}"
	variablesMap := make(map[string]interface{})
	variablesMap["foo"] = "bar"

	request := &pb.CreateWorkflowInstanceRequest{
		WorkflowKey: 123,
		Variables:   variables,
	}
	stub := &pb.CreateWorkflowInstanceResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
	}

	client.EXPECT().CreateWorkflowInstance(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromMap(variablesMap)
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)
	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			WorkflowKey: 123,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           "{}",
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	response, err := command.WorkflowKey(123).WithResult().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultCommandByBpmnProcessId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			BpmnProcessId: "foo",
			Version:       LatestVersion,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           "{}",
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	response, err := command.BPMNProcessId("foo").LatestVersion().WithResult().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultCommandByBpmnProcessIdAndVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			BpmnProcessId: "foo",
			Version:       56,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           "{}",
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	response, err := command.BPMNProcessId("foo").Version(56).WithResult().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultCommandWithVariablesFromString(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"bar\"}"

	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			WorkflowKey: 123,
			Variables:   variables,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           variables,
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromString(variables)
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.WithResult().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultCommandWithVariablesFromStringer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"bar\"}"

	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			WorkflowKey: 123,
			Variables:   variables,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           variables,
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromStringer(DataType{Foo: "bar"})
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.WithResult().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultCommandWithVariablesFromObject(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"bar\"}"

	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			WorkflowKey: 123,
			Variables:   variables,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           variables,
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromObject(DataType{Foo: "bar"})
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.WithResult().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultCommandWithVariablesFromObjectOmitempty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{}"

	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			WorkflowKey: 123,
			Variables:   variables,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           variables,
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromObject(DataType{Foo: ""})
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.WithResult().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultCommandWithVariablesFromObjectIgnoreOmitempty(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"\"}"

	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			WorkflowKey: 123,
			Variables:   variables,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           variables,
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromObjectIgnoreOmitempty(DataType{Foo: ""})
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.WithResult().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultCommandWithVariablesFromMap(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	variables := "{\"foo\":\"bar\"}"
	variablesMap := make(map[string]interface{})
	variablesMap["foo"] = "bar"

	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			WorkflowKey: 123,
			Variables:   variables,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           variables,
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	variablesCommand, err := command.WorkflowKey(123).VariablesFromMap(variablesMap)
	if err != nil {
		t.Error("Failed to set variables: ", err)
	}

	response, err := variablesCommand.WithResult().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultAndFetchVariablesCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)
	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			WorkflowKey: 123,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
		FetchVariables: []string{"a", "b", "c"},
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           "{}",
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	response, err := command.WorkflowKey(123).WithResult().FetchVariables("a", "b", "c").Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestCreateWorkflowInstanceWithResultAndFetchEmptyVariablesListCommand(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)
	request := &pb.CreateWorkflowInstanceWithResultRequest{
		Request: &pb.CreateWorkflowInstanceRequest{
			WorkflowKey: 123,
		},
		RequestTimeout: utils.DefaultTestTimeoutInMs,
		FetchVariables: []string{},
	}
	stub := &pb.CreateWorkflowInstanceWithResultResponse{
		WorkflowKey:         123,
		BpmnProcessId:       "foo",
		Version:             4545,
		WorkflowInstanceKey: 5632,
		Variables:           "{}",
	}

	client.EXPECT().CreateWorkflowInstanceWithResult(gomock.Any(), &utils.RpcTestMsg{Msg: request}).Return(stub, nil)

	command := NewCreateInstanceCommand(client, utils.DefaultTestTimeout, func(error) bool { return false })

	response, err := command.WorkflowKey(123).WithResult().FetchVariables().Send()

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}
