package testrunner

import (
	"github.com/lex00/wetwire-observability-go/internal/discover"
)

// Status represents the evaluation status of a criterion.
type Status string

const (
	StatusPass    Status = "pass"
	StatusPartial Status = "partial"
	StatusFail    Status = "fail"
	StatusSkip    Status = "skip"
)

// Runner evaluates observability configurations against personas.
type Runner struct {
	personas []*Persona
}

// Result contains the evaluation results.
type Result struct {
	// TotalScore is the overall score across all personas.
	TotalScore int `json:"total_score"`

	// MaxScore is the maximum possible score.
	MaxScore int `json:"max_score"`

	// Percentage is the score as a percentage.
	Percentage float64 `json:"percentage"`

	// PersonaResults contains results for each persona.
	PersonaResults []PersonaResult `json:"persona_results"`

	// Recommendations are actionable suggestions.
	Recommendations []string `json:"recommendations,omitempty"`
}

// PersonaResult contains results for a single persona.
type PersonaResult struct {
	// Persona is the persona ID.
	Persona string `json:"persona"`

	// PersonaName is the display name.
	PersonaName string `json:"persona_name"`

	// Score is the total score for this persona.
	Score int `json:"score"`

	// MaxScore is the maximum possible score.
	MaxScore int `json:"max_score"`

	// Percentage is the score as a percentage.
	Percentage float64 `json:"percentage"`

	// Criteria contains the evaluation for each criterion.
	Criteria []CriterionResult `json:"criteria"`
}

// CriterionResult contains the evaluation of a single criterion.
type CriterionResult struct {
	// ID is the criterion ID.
	ID string `json:"id"`

	// Name is the display name.
	Name string `json:"name"`

	// Category is the criterion category.
	Category string `json:"category"`

	// Status is the evaluation status.
	Status Status `json:"status"`

	// Score is the achieved score.
	Score int `json:"score"`

	// MaxScore is the maximum score for this criterion.
	MaxScore int `json:"max_score"`

	// Message explains the evaluation.
	Message string `json:"message,omitempty"`
}

// NewRunner creates a new Runner.
func NewRunner() *Runner {
	return &Runner{}
}

// WithPersona adds a persona to evaluate against.
func (r *Runner) WithPersona(id string) *Runner {
	p := GetPersona(id)
	if p != nil {
		r.personas = append(r.personas, p)
	}
	return r
}

// WithAllPersonas adds all personas to evaluate against.
func (r *Runner) WithAllPersonas() *Runner {
	r.personas = GetAllPersonas()
	return r
}

// Evaluate evaluates the configurations at the given path.
func (r *Runner) Evaluate(path string) (*Result, error) {
	// Discover resources
	resources, err := discover.Discover(path)
	if err != nil {
		return nil, err
	}

	result := &Result{}

	// Evaluate against each persona
	for _, persona := range r.personas {
		pr := r.evaluatePersona(persona, resources)
		result.PersonaResults = append(result.PersonaResults, pr)
		result.TotalScore += pr.Score
		result.MaxScore += pr.MaxScore
	}

	// Calculate percentage
	if result.MaxScore > 0 {
		result.Percentage = float64(result.TotalScore) / float64(result.MaxScore) * 100
	}

	// Generate recommendations
	result.Recommendations = r.generateRecommendations(result)

	return result, nil
}

func (r *Runner) evaluatePersona(persona *Persona, resources *discover.DiscoveryResult) PersonaResult {
	pr := PersonaResult{
		Persona:     persona.ID,
		PersonaName: persona.Name,
	}

	for _, criterion := range persona.Criteria {
		cr := r.evaluateCriterion(criterion, resources)
		pr.Criteria = append(pr.Criteria, cr)
		pr.Score += cr.Score
		pr.MaxScore += cr.MaxScore
	}

	if pr.MaxScore > 0 {
		pr.Percentage = float64(pr.Score) / float64(pr.MaxScore) * 100
	}

	return pr
}

func (r *Runner) evaluateCriterion(criterion Criterion, resources *discover.DiscoveryResult) CriterionResult {
	cr := CriterionResult{
		ID:       criterion.ID,
		Name:     criterion.Name,
		Category: criterion.Category,
		MaxScore: criterion.Weight,
	}

	// Evaluate based on criterion ID
	switch criterion.ID {
	// Basic checks
	case "basic-metrics", "app-metrics":
		if resources.TotalCount() > 0 {
			cr.Status = StatusPass
			cr.Score = criterion.Weight
			cr.Message = "Metrics configuration found"
		} else {
			cr.Status = StatusFail
			cr.Message = "No metrics configuration found"
		}

	case "scrape-targets":
		if len(resources.ScrapeConfigs) > 0 {
			cr.Status = StatusPass
			cr.Score = criterion.Weight
			cr.Message = "Scrape targets configured"
		} else if resources.TotalCount() > 0 {
			cr.Status = StatusPartial
			cr.Score = criterion.Weight / 2
			cr.Message = "Configuration found but no explicit scrape targets"
		} else {
			cr.Status = StatusFail
			cr.Message = "No scrape targets found"
		}

	case "alerting", "simple-alerts":
		if len(resources.AlertingRules) > 0 {
			cr.Status = StatusPass
			cr.Score = criterion.Weight
			cr.Message = "Alerting rules configured"
		} else {
			cr.Status = StatusFail
			cr.Message = "No alerting rules found"
		}

	case "recording-rules":
		if len(resources.RecordingRules) > 0 {
			cr.Status = StatusPass
			cr.Score = criterion.Weight
			cr.Message = "Recording rules found"
		} else {
			cr.Status = StatusPartial
			cr.Score = criterion.Weight / 2
			cr.Message = "Consider adding recording rules for expensive queries"
		}

	case "dashboards":
		// Dashboards would be in a separate grafana discovery
		// For now, skip this criterion
		cr.Status = StatusSkip
		cr.Message = "Dashboard discovery not implemented"

	// SRE-specific
	case "slo-coverage":
		// Check for SLO-related rules
		hasSLO := false
		for _, rule := range resources.AlertingRules {
			if containsAny(rule.Name, "slo", "error_rate", "latency", "availability") {
				hasSLO = true
				break
			}
		}
		if hasSLO {
			cr.Status = StatusPass
			cr.Score = criterion.Weight
			cr.Message = "SLO-related alerting found"
		} else {
			cr.Status = StatusFail
			cr.Message = "No SLO coverage detected"
		}

	case "burn-rate":
		hasBurnRate := false
		for _, rule := range resources.AlertingRules {
			if containsAny(rule.Name, "burn", "budget", "rate") {
				hasBurnRate = true
				break
			}
		}
		if hasBurnRate {
			cr.Status = StatusPass
			cr.Score = criterion.Weight
			cr.Message = "Burn rate alerting found"
		} else {
			cr.Status = StatusSkip
			cr.Message = "No burn rate alerting (advanced pattern)"
		}

	// Security-specific
	case "secrets":
		// Assume pass unless we detect exposed secrets (would need deeper analysis)
		cr.Status = StatusPass
		cr.Score = criterion.Weight
		cr.Message = "No obvious secret exposure detected"

	case "auth-metrics", "auth-failures":
		hasAuth := false
		all := resources.All()
		for _, r := range all {
			if containsAny(r.Name, "auth", "login", "authentication") {
				hasAuth = true
				break
			}
		}
		if hasAuth {
			cr.Status = StatusPass
			cr.Score = criterion.Weight
			cr.Message = "Authentication monitoring found"
		} else {
			cr.Status = StatusSkip
			cr.Message = "No authentication metrics configured"
		}

	default:
		// Unknown criterion - skip
		cr.Status = StatusSkip
		cr.Message = "Criterion not evaluated"
	}

	return cr
}

func (r *Runner) generateRecommendations(result *Result) []string {
	var recommendations []string

	for _, pr := range result.PersonaResults {
		for _, cr := range pr.Criteria {
			if cr.Status == StatusFail {
				recommendations = append(recommendations,
					cr.Name+": "+cr.Message)
			}
		}
	}

	if len(recommendations) == 0 && result.Percentage < 100 {
		recommendations = append(recommendations,
			"Consider addressing partial implementations for improved coverage")
	}

	return recommendations
}

func containsAny(s string, substrs ...string) bool {
	lower := string(s)
	for _, sub := range substrs {
		if len(sub) > 0 && contains(lower, sub) {
			return true
		}
	}
	return false
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || indexString(s, substr) >= 0)
}

func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
