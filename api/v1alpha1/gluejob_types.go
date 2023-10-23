/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/glue@v1.63.0/types#JobCommand
type GlueJobCommand struct {
	Name string `json:"name"`
	// +kubebuilder:default=3
	// +kubebuilder:validation:Format=`^(2|3)$`
	PythonVersion int `json:"pythonVersion,omitempty"`
	// +kubebuilder:default=glueetl
	Runtime string `json:"runtime,omitempty"`
	// +kubebuilder:validation:Format=`^s3://.+\/.+$`
	ScriptLocation string `json:"scriptLocation"`
}

type GlueJobExecutionProperty struct {
	// +kubebuilder:default=1
	MaxConcurrentRuns int32 `json:"maxConcurrentRuns,omitempty"`
}

// GlueJobSpec defines the desired state of GlueJob
type GlueJobSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Name is the name of the Glue Job
	Name string `json:"name"`

	// Command is the Glue Job Command https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/service/glue@v1.63.0/types#JobCommand
	Command GlueJobCommand `json:"command"`

	// Role is the IAM role to be used by the Glue Job
	// +kubebuilder:validation:Format=`^arn:aws:iam::.*:role\/.*$`
	Role string `json:"role"`

	// Timeout is the timeout in minutes for the Glue Job, max 2 days
	// +kubebuilder:default=20
	// +kubebuilder:validation:Maximum=2880
	TimeoutInMinutes int32 `json:"timeout,omitempty"`

	// GlueVersion is the version of Glue to be used by the Glue Job
	// +kubebuilder:default="4.0"
	GlueVersion string `json:"glueVersion,omitempty"`

	// NumberOfWorkers is the number of workers to be used by the Glue Job
	// +kubebuilder:default=2
	NumberOfWorkers int32 `json:"numberOfWorkers,omitempty"`

	// WorkerType is the type of worker to be used by the Glue Job
	// +kubebuilder:default=G.1X
	WorkerType string `json:"workerType,omitempty"`

	// ExecutionProperty is the execution property to be used by the Glue Job
	// +kubebuilder:default=FLEX
	// +kubebuilder:validation:Format=`^(FLEX|STANDARD)$`
	ExecutionClass string `json:"executionClass,omitempty"`

	// ExecutionProperty is the execution property to be used by the Glue Job
	// +kubebuilder:default={maxConcurrentRuns: 1}
	ExecutionProperty *GlueJobExecutionProperty `json:"executionProperty,omitempty"`

	// MaxRetries is the max number of retries to be used by the Glue Job
	// +kubebuilder:default=0
	MaxRetries int32 `json:"maxRetries,omitempty"`

	// DefaultArguments is the default arguments to be used by the Glue Job
	DefaultArguments map[string]string `json:"defaultArguments,omitempty"`
}

// GlueJobStatus defines the observed state of GlueJob
type GlueJobStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Conditions store the status conditions of the GlueJob instances
	// +operator-sdk:csv:customresourcedefinitions:type=status
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,1,rep,name=conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// GlueJob is the Schema for the gluejobs API
type GlueJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GlueJobSpec   `json:"spec,omitempty"`
	Status GlueJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GlueJobList contains a list of GlueJob
type GlueJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GlueJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&GlueJob{}, &GlueJobList{})
}
