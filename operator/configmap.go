package operator

import (
	"encoding/json"

	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-observability-go/grafana"
)

// K8sConfigMap represents a Kubernetes ConfigMap.
type K8sConfigMap struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   ObjectMeta        `yaml:"metadata"`
	Data       map[string]string `yaml:"data,omitempty"`
	BinaryData map[string][]byte `yaml:"binaryData,omitempty"`

	// Convenience fields (not serialized)
	Name      string            `yaml:"-"`
	Namespace string            `yaml:"-"`
	Labels    map[string]string `yaml:"-"`
}

// ConfigMap creates a new K8sConfigMap.
func ConfigMap(name, namespace string) *K8sConfigMap {
	return &K8sConfigMap{
		APIVersion: "v1",
		Kind:       "ConfigMap",
		Metadata: ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Name:      name,
		Namespace: namespace,
		Data:      make(map[string]string),
	}
}

// WithLabels sets the ConfigMap labels.
func (cm *K8sConfigMap) WithLabels(labels map[string]string) *K8sConfigMap {
	cm.Labels = labels
	cm.Metadata.Labels = labels
	return cm
}

// AddLabel adds a label to the ConfigMap.
func (cm *K8sConfigMap) AddLabel(key, value string) *K8sConfigMap {
	if cm.Labels == nil {
		cm.Labels = make(map[string]string)
	}
	cm.Labels[key] = value
	cm.Metadata.Labels = cm.Labels
	return cm
}

// WithAnnotations sets the ConfigMap annotations.
func (cm *K8sConfigMap) WithAnnotations(annotations map[string]string) *K8sConfigMap {
	cm.Metadata.Annotations = annotations
	return cm
}

// AddAnnotation adds an annotation to the ConfigMap.
func (cm *K8sConfigMap) AddAnnotation(key, value string) *K8sConfigMap {
	if cm.Metadata.Annotations == nil {
		cm.Metadata.Annotations = make(map[string]string)
	}
	cm.Metadata.Annotations[key] = value
	return cm
}

// WithData sets the ConfigMap data.
func (cm *K8sConfigMap) WithData(data map[string]string) *K8sConfigMap {
	cm.Data = data
	return cm
}

// AddData adds a data entry to the ConfigMap.
func (cm *K8sConfigMap) AddData(key, value string) *K8sConfigMap {
	if cm.Data == nil {
		cm.Data = make(map[string]string)
	}
	cm.Data[key] = value
	return cm
}

// WithBinaryData sets the ConfigMap binary data.
func (cm *K8sConfigMap) WithBinaryData(data map[string][]byte) *K8sConfigMap {
	cm.BinaryData = data
	return cm
}

// AddBinaryData adds a binary data entry to the ConfigMap.
func (cm *K8sConfigMap) AddBinaryData(key string, value []byte) *K8sConfigMap {
	if cm.BinaryData == nil {
		cm.BinaryData = make(map[string][]byte)
	}
	cm.BinaryData[key] = value
	return cm
}

// ForGrafanaSidecar adds the standard Grafana sidecar label.
// This makes the ConfigMap discoverable by the Grafana sidecar.
func (cm *K8sConfigMap) ForGrafanaSidecar() *K8sConfigMap {
	return cm.AddLabel("grafana_dashboard", "1")
}

// WithFolder sets the Grafana folder annotation.
func (cm *K8sConfigMap) WithFolder(folder string) *K8sConfigMap {
	return cm.AddAnnotation("grafana_folder", folder)
}

// WithGrafanaLabel sets a custom Grafana label.
// Use this when your Grafana sidecar uses a different label selector.
func (cm *K8sConfigMap) WithGrafanaLabel(key, value string) *K8sConfigMap {
	return cm.AddLabel(key, value)
}

// Serialize converts the ConfigMap to YAML bytes.
func (cm *K8sConfigMap) Serialize() ([]byte, error) {
	return yaml.Marshal(cm)
}

// MustSerialize converts the ConfigMap to YAML bytes, panicking on error.
func (cm *K8sConfigMap) MustSerialize() []byte {
	data, err := cm.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}

// DashboardConfigMap creates a ConfigMap containing a Grafana dashboard.
// The dashboard is serialized to JSON and stored with a .json extension.
// The ConfigMap is labeled for Grafana sidecar discovery.
func DashboardConfigMap(name, namespace string, dashboard *grafana.Dashboard) *K8sConfigMap {
	cm := ConfigMap(name, namespace).
		ForGrafanaSidecar()

	// Serialize dashboard to JSON
	dashJSON, err := json.MarshalIndent(dashboard, "", "  ")
	if err != nil {
		// Return empty ConfigMap on error - caller can check Data
		return cm
	}

	cm.AddData(name+".json", string(dashJSON))
	return cm
}

// DashboardConfigMapJSON creates a ConfigMap from raw dashboard JSON.
// Use this when you have pre-serialized dashboard JSON.
func DashboardConfigMapJSON(name, namespace string, dashboardJSON string) *K8sConfigMap {
	return ConfigMap(name, namespace).
		ForGrafanaSidecar().
		AddData(name+".json", dashboardJSON)
}
