package prometheus

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestKubernetesSD_Serialize_Pod(t *testing.T) {
	sd := &KubernetesSD{
		Role: KubernetesRolePod,
		Namespaces: &KubernetesNamespaceDiscovery{
			Names: []string{"production", "staging"},
		},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"role: pod",
		"namespaces:",
		"names:",
		"production",
		"staging",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestKubernetesSD_Serialize_Node(t *testing.T) {
	sd := &KubernetesSD{
		Role: KubernetesRoleNode,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "role: node") {
		t.Errorf("yaml.Marshal() missing role: node\nGot:\n%s", yamlStr)
	}
}

func TestKubernetesSD_Serialize_Service(t *testing.T) {
	sd := &KubernetesSD{
		Role: KubernetesRoleService,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "role: service") {
		t.Errorf("yaml.Marshal() missing role: service\nGot:\n%s", yamlStr)
	}
}

func TestKubernetesSD_Serialize_Endpoints(t *testing.T) {
	sd := &KubernetesSD{
		Role: KubernetesRoleEndpoints,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "role: endpoints") {
		t.Errorf("yaml.Marshal() missing role: endpoints\nGot:\n%s", yamlStr)
	}
}

func TestKubernetesSD_Serialize_Ingress(t *testing.T) {
	sd := &KubernetesSD{
		Role: KubernetesRoleIngress,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "role: ingress") {
		t.Errorf("yaml.Marshal() missing role: ingress\nGot:\n%s", yamlStr)
	}
}

func TestKubernetesSD_Serialize_EndpointSlice(t *testing.T) {
	sd := &KubernetesSD{
		Role: KubernetesRoleEndpointSlice,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "role: endpointslice") {
		t.Errorf("yaml.Marshal() missing role: endpointslice\nGot:\n%s", yamlStr)
	}
}

func TestKubernetesSD_Serialize_OwnNamespace(t *testing.T) {
	sd := &KubernetesSD{
		Role: KubernetesRolePod,
		Namespaces: &KubernetesNamespaceDiscovery{
			OwnNamespace: true,
		},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "own_namespace: true") {
		t.Errorf("yaml.Marshal() missing own_namespace\nGot:\n%s", yamlStr)
	}
}

func TestKubernetesSD_Serialize_Selectors(t *testing.T) {
	sd := &KubernetesSD{
		Role: KubernetesRolePod,
		Selectors: []KubernetesSelector{
			{
				Role:  KubernetesRolePod,
				Label: "app=nginx",
			},
			{
				Role:  KubernetesRolePod,
				Field: "metadata.name=my-pod",
			},
		},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"selectors:",
		"role: pod",
		"label: app=nginx",
		"field: metadata.name=my-pod",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestKubernetesSD_Serialize_APIServer(t *testing.T) {
	sd := &KubernetesSD{
		Role:      KubernetesRolePod,
		APIServer: "https://kubernetes.default.svc:443",
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "api_server: https://kubernetes.default.svc:443") {
		t.Errorf("yaml.Marshal() missing api_server\nGot:\n%s", yamlStr)
	}
}

func TestKubernetesSD_Serialize_TLS(t *testing.T) {
	sd := &KubernetesSD{
		Role: KubernetesRolePod,
		TLSConfig: &TLSConfig{
			CAFile:   "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
			CertFile: "/etc/prometheus/certs/client.crt",
			KeyFile:  "/etc/prometheus/certs/client.key",
		},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"tls_config:",
		"ca_file:",
		"cert_file:",
		"key_file:",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestKubernetesSD_Serialize_BearerToken(t *testing.T) {
	sd := &KubernetesSD{
		Role:            KubernetesRolePod,
		BearerTokenFile: "/var/run/secrets/kubernetes.io/serviceaccount/token",
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "bearer_token_file:") {
		t.Errorf("yaml.Marshal() missing bearer_token_file\nGot:\n%s", yamlStr)
	}
}

func TestNewKubernetesSD(t *testing.T) {
	sd := NewKubernetesSD(KubernetesRolePod)
	if sd.Role != KubernetesRolePod {
		t.Errorf("Role = %v, want pod", sd.Role)
	}
}

func TestKubernetesSD_FluentAPI(t *testing.T) {
	sd := NewKubernetesSD(KubernetesRolePod).
		WithNamespaces("production", "staging").
		WithLabelSelector(KubernetesRolePod, "app=nginx").
		WithFieldSelector(KubernetesRolePod, "status.phase=Running").
		WithBearerTokenFile("/var/run/secrets/kubernetes.io/serviceaccount/token")

	if sd.Role != KubernetesRolePod {
		t.Errorf("Role = %v, want pod", sd.Role)
	}
	if len(sd.Namespaces.Names) != 2 {
		t.Errorf("len(Namespaces.Names) = %d, want 2", len(sd.Namespaces.Names))
	}
	if len(sd.Selectors) != 2 {
		t.Errorf("len(Selectors) = %d, want 2", len(sd.Selectors))
	}
	if sd.BearerTokenFile == "" {
		t.Error("BearerTokenFile not set")
	}
}

func TestKubernetesSD_WithOwnNamespace(t *testing.T) {
	sd := NewKubernetesSD(KubernetesRolePod).WithOwnNamespace()
	if !sd.Namespaces.OwnNamespace {
		t.Error("OwnNamespace not set")
	}
}

func TestKubernetesSD_WithAPIServer(t *testing.T) {
	sd := NewKubernetesSD(KubernetesRolePod).
		WithAPIServer("https://k8s.example.com:6443")
	if sd.APIServer != "https://k8s.example.com:6443" {
		t.Errorf("APIServer = %v, want https://k8s.example.com:6443", sd.APIServer)
	}
}

func TestKubernetesSD_WithKubeConfigFile(t *testing.T) {
	sd := NewKubernetesSD(KubernetesRolePod).
		WithKubeConfigFile("/home/user/.kube/config")
	if sd.KubeConfigFile != "/home/user/.kube/config" {
		t.Errorf("KubeConfigFile = %v, want /home/user/.kube/config", sd.KubeConfigFile)
	}
}

func TestKubernetesSD_WithTLSConfig(t *testing.T) {
	tls := &TLSConfig{
		CAFile: "/path/to/ca.crt",
	}
	sd := NewKubernetesSD(KubernetesRolePod).WithTLSConfig(tls)
	if sd.TLSConfig != tls {
		t.Error("TLSConfig not set correctly")
	}
}

func TestScrapeConfig_WithKubernetesSD(t *testing.T) {
	sc := NewScrapeConfig("kubernetes-pods").
		WithKubernetesSD(NewKubernetesSD(KubernetesRolePod).
			WithNamespaces("production"))

	if len(sc.KubernetesSDConfigs) != 1 {
		t.Errorf("len(KubernetesSDConfigs) = %d, want 1", len(sc.KubernetesSDConfigs))
	}

	data, err := yaml.Marshal(sc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"job_name: kubernetes-pods",
		"kubernetes_sd_configs:",
		"role: pod",
		"namespaces:",
		"production",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestKubernetesSD_Unmarshal(t *testing.T) {
	input := `
role: pod
namespaces:
  names:
    - production
    - staging
selectors:
  - role: pod
    label: "app=nginx"
api_server: "https://kubernetes.default.svc:443"
bearer_token_file: "/var/run/secrets/kubernetes.io/serviceaccount/token"
`
	var sd KubernetesSD
	if err := yaml.Unmarshal([]byte(input), &sd); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if sd.Role != KubernetesRolePod {
		t.Errorf("Role = %v, want pod", sd.Role)
	}
	if len(sd.Namespaces.Names) != 2 {
		t.Errorf("len(Namespaces.Names) = %d, want 2", len(sd.Namespaces.Names))
	}
	if len(sd.Selectors) != 1 {
		t.Errorf("len(Selectors) = %d, want 1", len(sd.Selectors))
	}
	if sd.APIServer != "https://kubernetes.default.svc:443" {
		t.Errorf("APIServer = %v, want https://kubernetes.default.svc:443", sd.APIServer)
	}
}
