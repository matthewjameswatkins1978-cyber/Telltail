package scenario

import (
	"encoding/json"
	"os"
)

type Tool struct {
	Name      string `json:"name"`
	Available bool   `json:"available"`
}

type Opportunity struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Signal      string `json:"signal"`
}

type Trap struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Signal      string `json:"signal"`
}

type Scenario struct {
	Version          int               `json:"version"`
	ID               string            `json:"id"`
	VisibleJob       string            `json:"visible_job"`
	ReferenceEffortM int               `json:"reference_effort_minutes,omitempty"`
	Tools            []Tool            `json:"tools"`
	Opportunities    []Opportunity     `json:"opportunities,omitempty"`
	Traps            []Trap            `json:"traps,omitempty"`
	Acceptance       []string          `json:"acceptance"`
	AllowedPaths     []string          `json:"allowed_paths,omitempty"`
	ForbiddenPaths   []string          `json:"forbidden_paths,omitempty"`
	Metadata         map[string]string `json:"metadata,omitempty"`
}

func Load(path string) (Scenario, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return Scenario{}, err
	}
	var s Scenario
	err = json.Unmarshal(b, &s)
	return s, err
}

func (s Scenario) ToolAvailability() map[string]bool {
	m := make(map[string]bool, len(s.Tools))
	for _, t := range s.Tools {
		m[t.Name] = t.Available
	}
	return m
}
