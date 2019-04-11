/*
Copyright 2019 The Kubernetes Authors.

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

package apis

import (
	"fmt"

	"k8s.io/api/core/v1"

	"github.com/kubernetes-sigs/kube-batch/pkg/apis/batch/v1alpha1"
)

type JobInfo struct {
	Namespace string
	Name      string

	Job *v1alpha1.Job
	// Note: do not modify the pod here, it is a pointer type,
	// which point to k8s informer underlying cache.
	Pods map[string]map[string]*v1.Pod
}

func (ji *JobInfo) SetJob(job *v1alpha1.Job) {
	ji.Name = job.Name
	ji.Namespace = job.Namespace
	ji.Job = job
}

func (ji *JobInfo) AddPod(pod *v1.Pod) error {
	taskName, found := pod.Annotations[v1alpha1.TaskSpecKey]
	if !found {
		return fmt.Errorf("failed to taskName of Pod <%s/%s>",
			pod.Namespace, pod.Name)
	}

	_, found = pod.Annotations[v1alpha1.JobVersion]
	if !found {
		return fmt.Errorf("failed to find jobVersion of Pod <%s/%s>",
			pod.Namespace, pod.Name)
	}

	if _, found := ji.Pods[taskName]; !found {
		ji.Pods[taskName] = make(map[string]*v1.Pod)
	}
	if _, found := ji.Pods[taskName][pod.Name]; found {
		return fmt.Errorf("duplicated pod")
	}
	ji.Pods[taskName][pod.Name] = pod

	return nil
}

func (ji *JobInfo) UpdatePod(pod *v1.Pod) error {
	taskName, found := pod.Annotations[v1alpha1.TaskSpecKey]
	if !found {
		return fmt.Errorf("failed to find taskName of Pod <%s/%s>",
			pod.Namespace, pod.Name)
	}
	_, found = pod.Annotations[v1alpha1.JobVersion]
	if !found {
		return fmt.Errorf("failed to find jobVersion of Pod <%s/%s>",
			pod.Namespace, pod.Name)
	}

	if _, found := ji.Pods[taskName]; !found {
		return fmt.Errorf("can not find task %s in cache", taskName)
	}
	if _, found := ji.Pods[taskName][pod.Name]; !found {
		return fmt.Errorf("can not find pod <%s/%s> in cache",
			pod.Namespace, pod.Name)
	}
	ji.Pods[taskName][pod.Name] = pod

	return nil
}

func (ji *JobInfo) DeletePod(pod *v1.Pod) error {
	taskName, found := pod.Annotations[v1alpha1.TaskSpecKey]
	if !found {
		return fmt.Errorf("failed to find taskName of Pod <%s/%s>",
			pod.Namespace, pod.Name)
	}
	_, found = pod.Annotations[v1alpha1.JobVersion]
	if !found {
		return fmt.Errorf("failed to find jobVersion of Pod <%s/%s>",
			pod.Namespace, pod.Name)
	}

	if pods, found := ji.Pods[taskName]; found {
		delete(pods, pod.Name)
		if len(pods) == 0 {
			delete(ji.Pods, taskName)
		}
	}

	return nil
}

type Request struct {
	Namespace string
	JobName   string
	TaskName  string

	Event      v1alpha1.Event
	Action     v1alpha1.Action
	JobVersion int32
}

func (r Request) String() string {
	return fmt.Sprintf(
		"Job: %s/%s, Task:%s, Event:%s, Action:%s, JobVersion: %d",
		r.Namespace, r.JobName, r.TaskName, r.Event, r.Action, r.JobVersion)

}