// +build e2e

/*
Copyright 2019 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package test

import (
	"testing"

	tb "github.com/tektoncd/pipeline/internal/builder/v1alpha1"
	knativetest "knative.dev/pkg/test"
)

const (
	epTaskName    = "ep-task"
	epTaskRunName = "ep-task-run"
)

// TestEntrypointRunningStepsInOrder is an integration test that will
// verify attempt to the get the entrypoint of a container image
// that doesn't have a cmd defined. In addition to making sure the steps
// are executed in the order specified
func TestEntrypointRunningStepsInOrder(t *testing.T) {
	c, namespace := setup(t)
	t.Parallel()

	knativetest.CleanupOnInterrupt(func() { tearDown(t, c, namespace) }, t.Logf)
	defer tearDown(t, c, namespace)

	t.Logf("Creating Task and TaskRun in namespace %s", namespace)
	task := tb.Task(epTaskName, tb.TaskSpec(
		tb.Step("ubuntu", tb.StepArgs("-c", "sleep 3 && touch foo")),
		tb.Step("ubuntu", tb.StepArgs("-c", "ls", "foo")),
	))
	if _, err := c.TaskClient.Create(task); err != nil {
		t.Fatalf("Failed to create Task: %s", err)
	}
	taskRun := tb.TaskRun(epTaskRunName, tb.TaskRunSpec(
		tb.TaskRunTaskRef(epTaskName), tb.TaskRunServiceAccountName("default"),
	))
	if _, err := c.TaskRunClient.Create(taskRun); err != nil {
		t.Fatalf("Failed to create TaskRun: %s", err)
	}

	t.Logf("Waiting for TaskRun in namespace %s to finish successfully", namespace)
	if err := WaitForTaskRunState(c, epTaskRunName, TaskRunSucceed(epTaskRunName), "TaskRunSuccess"); err != nil {
		t.Errorf("Error waiting for TaskRun to finish successfully: %s", err)
	}

}
