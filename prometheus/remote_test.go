package prometheus

import (
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestSecret_String(t *testing.T) {
	secret := Secret("my-secret-token")
	if string(secret) != "my-secret-token" {
		t.Errorf("Secret = %q, want my-secret-token", string(secret))
	}
}

func TestSecretFromEnv(t *testing.T) {
	secret := SecretFromEnv("PROMETHEUS_TOKEN")
	if string(secret) != "${PROMETHEUS_TOKEN}" {
		t.Errorf("SecretFromEnv = %q, want ${PROMETHEUS_TOKEN}", string(secret))
	}
}

func TestNewRemoteWrite(t *testing.T) {
	rw := NewRemoteWrite("http://localhost:9090/api/v1/write")
	if rw.URL != "http://localhost:9090/api/v1/write" {
		t.Errorf("URL = %v, want http://localhost:9090/api/v1/write", rw.URL)
	}
}

func TestRemoteWrite_FluentAPI(t *testing.T) {
	rw := NewRemoteWrite("http://thanos-receive:10908/api/v1/receive").
		WithName("thanos").
		WithTimeout(30 * Second).
		WithHeaders(map[string]string{"X-Tenant": "production"}).
		WithBasicAuth("user", "secret").
		WithQueueConfig(NewQueueConfig().
			WithCapacity(10000).
			WithMaxShards(50).
			WithMaxSamplesPerSend(5000)).
		WithMetadataConfig(NewMetadataConfig().
			WithSendInterval(Minute))

	if rw.Name != "thanos" {
		t.Errorf("Name = %v, want thanos", rw.Name)
	}
	if rw.RemoteTimeout != 30*Second {
		t.Errorf("RemoteTimeout = %v, want 30s", rw.RemoteTimeout)
	}
	if rw.Headers["X-Tenant"] != "production" {
		t.Errorf("Headers[X-Tenant] = %v, want production", rw.Headers["X-Tenant"])
	}
	if rw.BasicAuth == nil || rw.BasicAuth.Username != "user" {
		t.Error("BasicAuth not configured correctly")
	}
	if rw.QueueConfig == nil || rw.QueueConfig.Capacity != 10000 {
		t.Error("QueueConfig not configured correctly")
	}
	if rw.MetadataConfig == nil || rw.MetadataConfig.SendInterval != Minute {
		t.Error("MetadataConfig not configured correctly")
	}
}

func TestRemoteWrite_WithBearerToken(t *testing.T) {
	rw := NewRemoteWrite("http://cortex:9009/api/v1/push").
		WithBearerToken("my-token")

	if rw.BearerToken != "my-token" {
		t.Errorf("BearerToken = %v, want my-token", rw.BearerToken)
	}
}

func TestRemoteWrite_WithBearerTokenFile(t *testing.T) {
	rw := NewRemoteWrite("http://cortex:9009/api/v1/push").
		WithBearerTokenFile("/var/run/secrets/token")

	if rw.BearerTokenFile != "/var/run/secrets/token" {
		t.Errorf("BearerTokenFile = %v, want /var/run/secrets/token", rw.BearerTokenFile)
	}
}

func TestRemoteWrite_WithTLSConfig(t *testing.T) {
	tls := &TLSConfig{
		CAFile:   "/etc/prometheus/ca.crt",
		CertFile: "/etc/prometheus/client.crt",
		KeyFile:  "/etc/prometheus/client.key",
	}
	rw := NewRemoteWrite("https://secure-endpoint/write").WithTLSConfig(tls)

	if rw.TLSConfig != tls {
		t.Error("TLSConfig not set correctly")
	}
}

func TestRemoteWrite_WithWriteRelabelConfigs(t *testing.T) {
	rw := NewRemoteWrite("http://thanos/write").
		WithWriteRelabelConfigs(
			DropByLabel("__name__", "go_.*"),
			DropByLabel("__name__", "promhttp_.*"),
		)

	if len(rw.WriteRelabelConfigs) != 2 {
		t.Errorf("len(WriteRelabelConfigs) = %d, want 2", len(rw.WriteRelabelConfigs))
	}
}

func TestRemoteWrite_Serialize(t *testing.T) {
	rw := NewRemoteWrite("http://thanos-receive:10908/api/v1/receive").
		WithName("thanos").
		WithTimeout(30 * Second).
		WithQueueConfig(NewQueueConfig().
			WithCapacity(10000).
			WithMaxShards(50))

	data, err := yaml.Marshal(rw)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"url: http://thanos-receive:10908/api/v1/receive",
		"name: thanos",
		"remote_timeout: 30s",
		"queue_config:",
		"capacity: 10000",
		"max_shards: 50",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestNewRemoteRead(t *testing.T) {
	rr := NewRemoteRead("http://localhost:9090/api/v1/read")
	if rr.URL != "http://localhost:9090/api/v1/read" {
		t.Errorf("URL = %v, want http://localhost:9090/api/v1/read", rr.URL)
	}
}

func TestRemoteRead_FluentAPI(t *testing.T) {
	rr := NewRemoteRead("http://thanos-query:10902/api/v1/read").
		WithName("thanos-read").
		WithTimeout(Minute).
		WithReadRecent(true).
		WithRequiredMatchers(map[string]string{"cluster": "production"}).
		WithFilterExternalLabels(true)

	if rr.Name != "thanos-read" {
		t.Errorf("Name = %v, want thanos-read", rr.Name)
	}
	if rr.RemoteTimeout != Minute {
		t.Errorf("RemoteTimeout = %v, want 1m", rr.RemoteTimeout)
	}
	if !rr.ReadRecent {
		t.Error("ReadRecent not set")
	}
	if rr.RequiredMatchers["cluster"] != "production" {
		t.Errorf("RequiredMatchers[cluster] = %v, want production", rr.RequiredMatchers["cluster"])
	}
	if !rr.FilterExternalLabels {
		t.Error("FilterExternalLabels not set")
	}
}

func TestRemoteRead_Serialize(t *testing.T) {
	rr := NewRemoteRead("http://thanos-query:10902/api/v1/read").
		WithName("thanos").
		WithReadRecent(true).
		WithRequiredMatchers(map[string]string{"env": "prod"})

	data, err := yaml.Marshal(rr)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"url: http://thanos-query:10902/api/v1/read",
		"name: thanos",
		"read_recent: true",
		"required_matchers:",
		"env: prod",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestQueueConfig_FluentAPI(t *testing.T) {
	qc := NewQueueConfig().
		WithCapacity(2500).
		WithMaxShards(200).
		WithMinShards(1).
		WithMaxSamplesPerSend(500).
		WithBatchSendDeadline(5 * Second).
		WithBackoff(Duration(30*time.Millisecond), 5*Second).
		WithRetryOnRateLimit(true)

	if qc.Capacity != 2500 {
		t.Errorf("Capacity = %d, want 2500", qc.Capacity)
	}
	if qc.MaxShards != 200 {
		t.Errorf("MaxShards = %d, want 200", qc.MaxShards)
	}
	if qc.MinShards != 1 {
		t.Errorf("MinShards = %d, want 1", qc.MinShards)
	}
	if qc.MaxSamplesPerSend != 500 {
		t.Errorf("MaxSamplesPerSend = %d, want 500", qc.MaxSamplesPerSend)
	}
	if qc.BatchSendDeadline != 5*Second {
		t.Errorf("BatchSendDeadline = %v, want 5s", qc.BatchSendDeadline)
	}
	if qc.MinBackoff != Duration(30*time.Millisecond) {
		t.Errorf("MinBackoff = %v, want 30ms", qc.MinBackoff)
	}
	if qc.MaxBackoff != 5*Second {
		t.Errorf("MaxBackoff = %v, want 5s", qc.MaxBackoff)
	}
	if !qc.RetryOnRateLimit {
		t.Error("RetryOnRateLimit not set")
	}
}

func TestQueueConfig_Serialize(t *testing.T) {
	qc := NewQueueConfig().
		WithCapacity(5000).
		WithMaxShards(100).
		WithBatchSendDeadline(5 * Second)

	data, err := yaml.Marshal(qc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"capacity: 5000",
		"max_shards: 100",
		"batch_send_deadline: 5s",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestMetadataConfig_FluentAPI(t *testing.T) {
	mc := NewMetadataConfig().
		WithSendInterval(5 * Minute).
		WithMaxSamplesPerSend(2000)

	if !mc.Send {
		t.Error("Send should default to true")
	}
	if mc.SendInterval != 5*Minute {
		t.Errorf("SendInterval = %v, want 5m", mc.SendInterval)
	}
	if mc.MaxSamplesPerSend != 2000 {
		t.Errorf("MaxSamplesPerSend = %d, want 2000", mc.MaxSamplesPerSend)
	}
}

func TestMetadataConfig_Serialize(t *testing.T) {
	mc := NewMetadataConfig().
		WithSendInterval(Minute).
		WithMaxSamplesPerSend(500)

	data, err := yaml.Marshal(mc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"send: true",
		"send_interval: 1m",
		"max_samples_per_send: 500",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestPrometheusConfig_WithRemoteWrite(t *testing.T) {
	config := &PrometheusConfig{
		Global: &GlobalConfig{
			ScrapeInterval: 15 * Second,
		},
		RemoteWrite: []*RemoteWriteConfig{
			NewRemoteWrite("http://thanos:10908/api/v1/receive").
				WithName("thanos").
				WithTimeout(30 * Second).
				WithQueueConfig(NewQueueConfig().
					WithCapacity(10000)),
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"scrape_interval: 15s",
		"remote_write:",
		"url: http://thanos:10908/api/v1/receive",
		"name: thanos",
		"queue_config:",
		"capacity: 10000",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestPrometheusConfig_WithRemoteRead(t *testing.T) {
	config := &PrometheusConfig{
		Global: &GlobalConfig{
			ScrapeInterval: 15 * Second,
		},
		RemoteRead: []*RemoteReadConfig{
			NewRemoteRead("http://thanos-query:10902/api/v1/read").
				WithName("thanos").
				WithReadRecent(true).
				WithRequiredMatchers(map[string]string{"env": "prod"}),
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"remote_read:",
		"url: http://thanos-query:10902/api/v1/read",
		"name: thanos",
		"read_recent: true",
		"required_matchers:",
		"env: prod",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRemoteWrite_Unmarshal(t *testing.T) {
	input := `
url: http://cortex:9009/api/v1/push
name: cortex
remote_timeout: 30s
headers:
  X-Scope-OrgID: tenant1
basic_auth:
  username: writer
  password: secret
queue_config:
  capacity: 5000
  max_shards: 100
`
	var rw RemoteWriteConfig
	if err := yaml.Unmarshal([]byte(input), &rw); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if rw.URL != "http://cortex:9009/api/v1/push" {
		t.Errorf("URL = %v, want http://cortex:9009/api/v1/push", rw.URL)
	}
	if rw.Name != "cortex" {
		t.Errorf("Name = %v, want cortex", rw.Name)
	}
	if rw.Headers["X-Scope-OrgID"] != "tenant1" {
		t.Errorf("Headers[X-Scope-OrgID] = %v, want tenant1", rw.Headers["X-Scope-OrgID"])
	}
	if rw.BasicAuth == nil || rw.BasicAuth.Username != "writer" {
		t.Error("BasicAuth not parsed correctly")
	}
	if rw.QueueConfig == nil || rw.QueueConfig.Capacity != 5000 {
		t.Error("QueueConfig not parsed correctly")
	}
}

func TestRemoteRead_Unmarshal(t *testing.T) {
	input := `
url: http://thanos-query:10902/api/v1/read
name: thanos
remote_timeout: 1m
read_recent: true
required_matchers:
  cluster: production
  env: prod
`
	var rr RemoteReadConfig
	if err := yaml.Unmarshal([]byte(input), &rr); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if rr.URL != "http://thanos-query:10902/api/v1/read" {
		t.Errorf("URL = %v, want http://thanos-query:10902/api/v1/read", rr.URL)
	}
	if !rr.ReadRecent {
		t.Error("ReadRecent not parsed")
	}
	if rr.RequiredMatchers["cluster"] != "production" {
		t.Errorf("RequiredMatchers[cluster] = %v, want production", rr.RequiredMatchers["cluster"])
	}
}

func TestCompleteRemoteWriteExample(t *testing.T) {
	// This test verifies a complete Thanos remote write configuration
	config := &PrometheusConfig{
		Global: &GlobalConfig{
			ScrapeInterval:     15 * Second,
			EvaluationInterval: 15 * Second,
			ExternalLabels: map[string]string{
				"cluster":  "production",
				"region":   "us-west-2",
				"instance": "prometheus-0",
			},
		},
		ScrapeConfigs: []*ScrapeConfig{
			NewScrapeConfig("prometheus").
				WithStaticTargets("localhost:9090"),
		},
		RemoteWrite: []*RemoteWriteConfig{
			NewRemoteWrite("http://thanos-receive:10908/api/v1/receive").
				WithName("thanos").
				WithTimeout(30 * Second).
				WithWriteRelabelConfigs(
					// Drop high-cardinality metrics
					DropByLabel("__name__", "go_.*"),
					DropByLabel("__name__", "promhttp_.*"),
				).
				WithQueueConfig(NewQueueConfig().
					WithCapacity(10000).
					WithMaxShards(50).
					WithMaxSamplesPerSend(5000).
					WithBatchSendDeadline(5 * Second).
					WithBackoff(Duration(30*time.Millisecond), 5*Second).
					WithRetryOnRateLimit(true)).
				WithMetadataConfig(NewMetadataConfig().
					WithSendInterval(Minute)),
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	// Verify round-trip
	var restored PrometheusConfig
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if len(restored.RemoteWrite) != 1 {
		t.Errorf("len(RemoteWrite) = %d, want 1", len(restored.RemoteWrite))
	}

	rw := restored.RemoteWrite[0]
	if rw.URL != "http://thanos-receive:10908/api/v1/receive" {
		t.Errorf("URL = %v, want http://thanos-receive:10908/api/v1/receive", rw.URL)
	}
	if rw.QueueConfig == nil {
		t.Error("QueueConfig not restored")
	} else if rw.QueueConfig.Capacity != 10000 {
		t.Errorf("QueueConfig.Capacity = %d, want 10000", rw.QueueConfig.Capacity)
	}

	t.Logf("Generated config:\n%s", string(data))
}
