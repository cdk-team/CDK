package util

type K8sPod struct {
	APIVersion string        `yaml:"apiVersion"`
	Kind       string        `yaml:"kind"`
	Metadata   K8sObjectMeta `yaml:"metadata"`
	Spec       K8sPodSpec    `yaml:"spec"`
}

type K8sObjectMeta struct {
	Name      string            `yaml:"name"`
	Namespace string            `yaml:"namespace,omitempty"`
	Labels    map[string]string `yaml:"labels,omitempty"`
}

type K8sPodSpec struct {
	Containers []K8sContainer `yaml:"containers"`
}

type K8sContainer struct {
	Name    string             `yaml:"name"`
	Image   string             `yaml:"image"`
	Ports   []K8sContainerPort `yaml:"ports,omitempty"`
	Command []string           `yaml:"command,omitempty"`
}

type K8sContainerPort struct {
	ContainerPort int `yaml:"containerPort"`
}
