package testrunner

import (
	"testing"
)

func TestGetPersona_SRE(t *testing.T) {
	p := GetPersona("sre")
	if p == nil {
		t.Fatal("GetPersona(sre) returned nil")
	}
	if p.Name != "SRE" {
		t.Errorf("Name = %q, want SRE", p.Name)
	}
}

func TestGetPersona_Developer(t *testing.T) {
	p := GetPersona("developer")
	if p == nil {
		t.Fatal("GetPersona(developer) returned nil")
	}
	if p.Name != "Developer" {
		t.Errorf("Name = %q, want Developer", p.Name)
	}
}

func TestGetPersona_Security(t *testing.T) {
	p := GetPersona("security")
	if p == nil {
		t.Fatal("GetPersona(security) returned nil")
	}
	if p.Name != "Security Analyst" {
		t.Errorf("Name = %q, want Security Analyst", p.Name)
	}
}

func TestGetPersona_Unknown(t *testing.T) {
	p := GetPersona("unknown")
	if p != nil {
		t.Error("GetPersona(unknown) should return nil")
	}
}

func TestGetAllPersonas(t *testing.T) {
	personas := GetAllPersonas()
	if len(personas) < 3 {
		t.Errorf("len(personas) = %d, want >= 3", len(personas))
	}
}

func TestPersona_HasCriteria(t *testing.T) {
	p := GetPersona("sre")
	if len(p.Criteria) == 0 {
		t.Error("SRE persona should have criteria")
	}
}

func TestPersona_HasDescription(t *testing.T) {
	p := GetPersona("sre")
	if p.Description == "" {
		t.Error("SRE persona should have description")
	}
}

func TestListPersonaNames(t *testing.T) {
	names := ListPersonaNames()
	if len(names) < 3 {
		t.Errorf("len(names) = %d, want >= 3", len(names))
	}

	// Check that expected names are present
	found := make(map[string]bool)
	for _, name := range names {
		found[name] = true
	}

	expected := []string{"sre", "developer", "security"}
	for _, exp := range expected {
		if !found[exp] {
			t.Errorf("expected persona %q not found", exp)
		}
	}
}
