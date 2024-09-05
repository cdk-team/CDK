
/*
Copyright 2022 The Authors of https://github.com/CDK-TEAM/CDK .

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

package util

import (
	"time"
)

type PodList struct {
	Kind       string   `json:"kind"`
	APIVersion string   `json:"apiVersion"`
	Metadata   Metadata `json:"metadata"`
	Items      []Pod    `json:"items"`
}

type Metadata struct {
	Name              string            `json:"name,omitempty"`
	GenerateName      string            `json:"generateName,omitempty"`
	Namespace         string            `json:"namespace,omitempty"`
	UID               string            `json:"uid,omitempty"`
	ResourceVersion   string            `json:"resourceVersion,omitempty"`
	CreationTimestamp time.Time         `json:"creationTimestamp,omitempty"`
	Labels            map[string]string `json:"labels,omitempty"`
	Annotations       map[string]string `json:"annotations,omitempty"`
}

type OwnerReference struct {
	APIVersion         string `json:"apiVersion"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	UID                string `json:"uid"`
	Controller         bool   `json:"controller"`
	BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
}

type ManagedField struct {
	Manager     string    `json:"manager"`
	Operation   string    `json:"operation"`
	APIVersion  string    `json:"apiVersion"`
	Time        time.Time `json:"time"`
	FieldsType  string    `json:"fieldsType"`
	FieldsV1    FieldsV1  `json:"fieldsV1"`
}

type FieldsV1 struct {
	Metadata struct {
		GenerateName string            `json:"generateName"`
		Labels       map[string]string `json:"labels"`
		OwnerRefs    map[string]string `json:"ownerReferences"`
	} `json:"metadata"`
	Spec struct {
		Affinity      map[string]interface{} `json:"affinity"`
		Containers    map[string]interface{} `json:"containers"`
		RestartPolicy string                 `json:"restartPolicy"`
		SchedulerName string                 `json:"schedulerName"`
		SecurityContext struct {
			Sysctls []Sysctl `json:"sysctls"`
		} `json:"securityContext"`
	} `json:"spec"`
}

type Sysctl struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Pod struct {
	Metadata     Metadata       `json:"metadata"`
	Spec         PodSpec        `json:"spec"`
	Status       PodStatus      `json:"status"`
	ManagedFields []ManagedField `json:"managedFields"`
}

type PodSpec struct {
	Containers                   []Container   `json:"containers"`
	RestartPolicy                string        `json:"restartPolicy"`
	TerminationGracePeriodSeconds int           `json:"terminationGracePeriodSeconds"`
	DNSPolicy                    string        `json:"dnsPolicy"`
	ServiceAccountName           string        `json:"serviceAccountName"`
	NodeName                     string        `json:"nodeName"`
	SecurityContext              SecurityContext `json:"securityContext"`
	Affinity                     Affinity      `json:"affinity"`
	Tolerations                  []Toleration  `json:"tolerations"`
	SchedulerName                string        `json:"schedulerName"`
}

type SecurityContext struct {
	Sysctls []Sysctl `json:"sysctls"`
}

type Affinity struct {
	NodeAffinity NodeAffinity `json:"nodeAffinity"`
}

type NodeAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution struct {
		NodeSelectorTerms []NodeSelectorTerm `json:"nodeSelectorTerms"`
	} `json:"requiredDuringSchedulingIgnoredDuringExecution"`
}

type NodeSelectorTerm struct {
	MatchFields []MatchField `json:"matchFields"`
}

type MatchField struct {
	Key      string   `json:"key"`
	Operator string   `json:"operator"`
	Values   []string `json:"values"`
}

type Container struct {
	Name            string          `json:"name"`
	Image           string          `json:"image"`
	Ports           []Port          `json:"ports"`
	Env             []EnvVar        `json:"env"`
	Resources       Resources       `json:"resources"`
	SecurityContext SecurityContext `json:"securityContext"`
}

type Port struct {
	Name          string `json:"name"`
	HostPort      int    `json:"hostPort"`
	ContainerPort int    `json:"containerPort"`
	Protocol      string `json:"protocol"`
}

type EnvVar struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Resources struct{}

type PodStatus struct {
	Phase             string             `json:"phase"`
	Conditions        []Condition        `json:"conditions"`
	HostIP            string             `json:"hostIP"`
	PodIP             string             `json:"podIP"`
	PodIPs            []PodIP            `json:"podIPs"`
	StartTime         time.Time          `json:"startTime"`
	ContainerStatuses []ContainerStatus  `json:"containerStatuses"`
	QOSClass          string             `json:"qosClass"`
}

type Condition struct {
	Type               string    `json:"type"`
	Status             string    `json:"status"`
	LastProbeTime      time.Time `json:"lastProbeTime"`
	LastTransitionTime time.Time `json:"lastTransitionTime"`
}

type PodIP struct {
	IP string `json:"ip"`
}

type ContainerStatus struct {
	Name        string      `json:"name"`
	State       ContainerState `json:"state"`
	LastState   ContainerState `json:"lastState"`
	Ready       bool        `json:"ready"`
	RestartCount int         `json:"restartCount"`
	Image       string      `json:"image"`
	ImageID     string      `json:"imageID"`
	ContainerID string      `json:"containerID"`
	Started     bool        `json:"started"`
}

type ContainerState struct {
	Running    *ContainerStateRunning `json:"running,omitempty"`
	Terminated *ContainerStateTerminated `json:"terminated,omitempty"`
	Waiting    *ContainerStateWaiting `json:"waiting,omitempty"`
}

type ContainerStateRunning struct {
	StartedAt time.Time `json:"startedAt"`
}

type ContainerStateTerminated struct {
	ExitCode int    `json:"exitCode"`
	Reason   string `json:"reason"`
}

type ContainerStateWaiting struct {
	Reason string `json:"reason"`
}

type Toleration struct {
	Key               string `json:"key"`
	Operator          string `json:"operator"`
	Effect            string `json:"effect"`
	TolerationSeconds *int64 `json:"tolerationSeconds,omitempty"`
}


