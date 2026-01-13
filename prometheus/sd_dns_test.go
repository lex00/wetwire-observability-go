package prometheus

import (
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestDNSSD_Serialize_SRV(t *testing.T) {
	sd := &DNSSD{
		Names: []string{"_prometheus._tcp.example.com"},
		Type:  DNSSDTypeSRV,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"names:",
		"- _prometheus._tcp.example.com",
		"type: SRV",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestDNSSD_Serialize_A(t *testing.T) {
	sd := &DNSSD{
		Names: []string{"web.example.com"},
		Type:  DNSSDTypeA,
		Port:  9100,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"names:",
		"- web.example.com",
		"type: A",
		"port: 9100",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestDNSSD_Serialize_AAAA(t *testing.T) {
	sd := &DNSSD{
		Names: []string{"ipv6.example.com"},
		Type:  DNSSDTypeAAAA,
		Port:  9100,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "type: AAAA") {
		t.Errorf("yaml.Marshal() missing type: AAAA\nGot:\n%s", yamlStr)
	}
}

func TestDNSSD_Serialize_MultipleNames(t *testing.T) {
	sd := &DNSSD{
		Names: []string{
			"_prometheus._tcp.dc1.example.com",
			"_prometheus._tcp.dc2.example.com",
		},
		Type: DNSSDTypeSRV,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"names:",
		"- _prometheus._tcp.dc1.example.com",
		"- _prometheus._tcp.dc2.example.com",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestDNSSD_Serialize_WithRefreshInterval(t *testing.T) {
	sd := &DNSSD{
		Names:           []string{"service.example.com"},
		Type:            DNSSDTypeSRV,
		RefreshInterval: Duration(30 * time.Second),
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "refresh_interval: 30s") {
		t.Errorf("yaml.Marshal() missing refresh_interval\nGot:\n%s", yamlStr)
	}
}

func TestNewDNSSD(t *testing.T) {
	sd := NewDNSSD()
	if sd == nil {
		t.Error("NewDNSSD() returned nil")
	}
}

func TestDNSSD_FluentAPI(t *testing.T) {
	sd := NewDNSSD().
		WithNames("_prometheus._tcp.example.com", "_metrics._tcp.example.com").
		WithType(DNSSDTypeSRV).
		WithRefreshInterval(Duration(60 * time.Second))

	if len(sd.Names) != 2 {
		t.Errorf("len(Names) = %d, want 2", len(sd.Names))
	}
	if sd.Names[0] != "_prometheus._tcp.example.com" {
		t.Errorf("Names[0] = %v, want _prometheus._tcp.example.com", sd.Names[0])
	}
	if sd.Type != DNSSDTypeSRV {
		t.Errorf("Type = %v, want SRV", sd.Type)
	}
	if sd.RefreshInterval != Duration(60*time.Second) {
		t.Errorf("RefreshInterval = %v, want 60s", sd.RefreshInterval)
	}
}

func TestDNSSD_WithNames(t *testing.T) {
	sd := NewDNSSD().WithNames("service.example.com")
	if len(sd.Names) != 1 {
		t.Errorf("len(Names) = %d, want 1", len(sd.Names))
	}
	if sd.Names[0] != "service.example.com" {
		t.Errorf("Names[0] = %v, want service.example.com", sd.Names[0])
	}
}

func TestDNSSD_WithType_SRV(t *testing.T) {
	sd := NewDNSSD().WithType(DNSSDTypeSRV)
	if sd.Type != DNSSDTypeSRV {
		t.Errorf("Type = %v, want SRV", sd.Type)
	}
}

func TestDNSSD_WithType_A(t *testing.T) {
	sd := NewDNSSD().WithType(DNSSDTypeA)
	if sd.Type != DNSSDTypeA {
		t.Errorf("Type = %v, want A", sd.Type)
	}
}

func TestDNSSD_WithType_AAAA(t *testing.T) {
	sd := NewDNSSD().WithType(DNSSDTypeAAAA)
	if sd.Type != DNSSDTypeAAAA {
		t.Errorf("Type = %v, want AAAA", sd.Type)
	}
}

func TestDNSSD_WithType_MX(t *testing.T) {
	sd := NewDNSSD().WithType(DNSSDTypeMX)
	if sd.Type != DNSSDTypeMX {
		t.Errorf("Type = %v, want MX", sd.Type)
	}
}

func TestDNSSD_WithType_NS(t *testing.T) {
	sd := NewDNSSD().WithType(DNSSDTypeNS)
	if sd.Type != DNSSDTypeNS {
		t.Errorf("Type = %v, want NS", sd.Type)
	}
}

func TestDNSSD_WithPort(t *testing.T) {
	sd := NewDNSSD().WithPort(9100)
	if sd.Port != 9100 {
		t.Errorf("Port = %v, want 9100", sd.Port)
	}
}

func TestDNSSD_WithRefreshInterval(t *testing.T) {
	sd := NewDNSSD().WithRefreshInterval(Duration(45 * time.Second))
	if sd.RefreshInterval != Duration(45*time.Second) {
		t.Errorf("RefreshInterval = %v, want 45s", sd.RefreshInterval)
	}
}

func TestDNSSD_TypeConstants(t *testing.T) {
	// Verify that constants have expected string values
	if DNSSDTypeSRV != "SRV" {
		t.Errorf("DNSSDTypeSRV = %v, want SRV", DNSSDTypeSRV)
	}
	if DNSSDTypeA != "A" {
		t.Errorf("DNSSDTypeA = %v, want A", DNSSDTypeA)
	}
	if DNSSDTypeAAAA != "AAAA" {
		t.Errorf("DNSSDTypeAAAA = %v, want AAAA", DNSSDTypeAAAA)
	}
	if DNSSDTypeMX != "MX" {
		t.Errorf("DNSSDTypeMX = %v, want MX", DNSSDTypeMX)
	}
	if DNSSDTypeNS != "NS" {
		t.Errorf("DNSSDTypeNS = %v, want NS", DNSSDTypeNS)
	}
}

func TestScrapeConfig_WithDNSSD(t *testing.T) {
	sc := NewScrapeConfig("dns-discovery").
		WithDNSSD(NewDNSSD().
			WithNames("_prometheus._tcp.example.com").
			WithType(DNSSDTypeSRV).
			WithRefreshInterval(Duration(30 * time.Second)))

	if len(sc.DNSSDConfigs) != 1 {
		t.Errorf("len(DNSSDConfigs) = %d, want 1", len(sc.DNSSDConfigs))
	}

	data, err := yaml.Marshal(sc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"job_name: dns-discovery",
		"dns_sd_configs:",
		"names:",
		"- _prometheus._tcp.example.com",
		"type: SRV",
		"refresh_interval: 30s",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestDNSSD_Unmarshal_SRV(t *testing.T) {
	input := `
names:
  - _prometheus._tcp.example.com
  - _metrics._tcp.example.com
type: SRV
refresh_interval: 30s
`
	var sd DNSSD
	if err := yaml.Unmarshal([]byte(input), &sd); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if len(sd.Names) != 2 {
		t.Errorf("len(Names) = %d, want 2", len(sd.Names))
	}
	if sd.Names[0] != "_prometheus._tcp.example.com" {
		t.Errorf("Names[0] = %v, want _prometheus._tcp.example.com", sd.Names[0])
	}
	if sd.Type != DNSSDTypeSRV {
		t.Errorf("Type = %v, want SRV", sd.Type)
	}
	if sd.RefreshInterval != Duration(30*time.Second) {
		t.Errorf("RefreshInterval = %v, want 30s", sd.RefreshInterval)
	}
}

func TestDNSSD_Unmarshal_A(t *testing.T) {
	input := `
names:
  - web.example.com
type: A
port: 9100
refresh_interval: 1m
`
	var sd DNSSD
	if err := yaml.Unmarshal([]byte(input), &sd); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if sd.Type != DNSSDTypeA {
		t.Errorf("Type = %v, want A", sd.Type)
	}
	if sd.Port != 9100 {
		t.Errorf("Port = %v, want 9100", sd.Port)
	}
	if sd.RefreshInterval != Duration(60*time.Second) {
		t.Errorf("RefreshInterval = %v, want 1m", sd.RefreshInterval)
	}
}

func TestScrapeConfig_MultipleDNSSD(t *testing.T) {
	sc := NewScrapeConfig("multi-dns").
		WithDNSSD(NewDNSSD().
			WithNames("_prometheus._tcp.dc1.example.com").
			WithType(DNSSDTypeSRV)).
		WithDNSSD(NewDNSSD().
			WithNames("_prometheus._tcp.dc2.example.com").
			WithType(DNSSDTypeSRV))

	if len(sc.DNSSDConfigs) != 2 {
		t.Errorf("len(DNSSDConfigs) = %d, want 2", len(sc.DNSSDConfigs))
	}

	data, err := yaml.Marshal(sc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "dns_sd_configs:") {
		t.Errorf("yaml.Marshal() missing dns_sd_configs\nGot:\n%s", yamlStr)
	}
}
