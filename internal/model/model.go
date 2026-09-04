package model

import "time"

type EventType string

const (
	EventShiftStart    EventType = "shift_start"
	EventShiftEnd      EventType = "shift_end"
	EventModelCall     EventType = "model_call"
	EventToolCall      EventType = "tool_call"
	EventToolResult    EventType = "tool_result"
	EventCommand       EventType = "command"
	EventCommandResult EventType = "command_result"
	EventFileChange    EventType = "file_change"
	EventAcceptance    EventType = "acceptance"
	EventClaim         EventType = "claim"
	EventProgress      EventType = "progress"
	EventSupervisor    EventType = "supervisor"
	EventQueue         EventType = "queue"
	EventProvision     EventType = "provision"
	EventTeardown      EventType = "teardown"
	EventInfo          EventType = "info"
)

type Event struct {
	Version    int               `json:"version"`
	Seq        int64             `json:"seq"`
	Time       time.Time         `json:"time"`
	Type       EventType         `json:"type"`
	Actor      string            `json:"actor,omitempty"`
	Name       string            `json:"name,omitempty"`
	Target     string            `json:"target,omitempty"`
	Input      string            `json:"input,omitempty"`
	Output     string            `json:"output,omitempty"`
	Success    *bool             `json:"success,omitempty"`
	DurationMS int64             `json:"duration_ms,omitempty"`
	CostUSD    float64           `json:"cost_usd,omitempty"`
	TokensIn   int64             `json:"tokens_in,omitempty"`
	TokensOut  int64             `json:"tokens_out,omitempty"`
	Progress   bool              `json:"progress,omitempty"`
	Authorized *bool             `json:"authorized,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
	PrevHash   string            `json:"prev_hash,omitempty"`
	Hash       string            `json:"hash,omitempty"`
}

type Severity int

const (
	SeverityNuisance  Severity = 1
	SeverityWaste     Severity = 2
	SeverityJob       Severity = 3
	SeverityIntegrity Severity = 4
	SeverityBoundary  Severity = 5
)

type Finding struct {
	Attribution string   `json:"attribution,omitempty"`
	Code        string   `json:"code"`
	Severity    Severity `json:"severity"`
	Seq         int64    `json:"seq,omitempty"`
	Actor       string   `json:"actor,omitempty"`
	Summary     string   `json:"summary"`
	Evidence    string   `json:"evidence,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	RepeatCount int      `json:"repeat_count,omitempty"`
	AvoidableMS int64    `json:"avoidable_ms,omitempty"`
}

type PhaseTime struct {
	Phase      string `json:"phase"`
	DurationMS int64  `json:"duration_ms"`
}

type Profile struct {
	TotalDurationMS      int64       `json:"total_duration_ms"`
	ModelDurationMS      int64       `json:"model_duration_ms"`
	ToolDurationMS       int64       `json:"tool_duration_ms"`
	CommandDurationMS    int64       `json:"command_duration_ms"`
	QueueDurationMS      int64       `json:"queue_duration_ms"`
	ProvisionDurationMS  int64       `json:"provision_duration_ms"`
	AcceptanceDurationMS int64       `json:"acceptance_duration_ms"`
	TeardownDurationMS   int64       `json:"teardown_duration_ms"`
	FailedWorkMS         int64       `json:"failed_work_ms"`
	RepeatedWorkMS       int64       `json:"repeated_work_ms"`
	EstimatedAvoidableMS int64       `json:"estimated_avoidable_ms"`
	CostUSD              float64     `json:"cost_usd"`
	TokensIn             int64       `json:"tokens_in"`
	TokensOut            int64       `json:"tokens_out"`
	Bottlenecks          []PhaseTime `json:"bottlenecks,omitempty"`
}

type ShiftResult struct {
	ScenarioID string    `json:"scenario_id"`
	Worker     string    `json:"worker"`
	Backend    string    `json:"backend"`
	Accepted   bool      `json:"accepted"`
	Outcome    string    `json:"outcome"`
	Findings   []Finding `json:"findings"`
	Profile    Profile   `json:"profile"`
	TraceHead  string    `json:"trace_head"`
}

type Dossier struct {
	Worker               string             `json:"worker"`
	Shifts               int                `json:"shifts"`
	Accepted             int                `json:"accepted"`
	CompletionRate       float64            `json:"completion_rate"`
	CleanSuccesses       int                `json:"clean_successes"`
	RecoveredSuccesses   int                `json:"recovered_successes"`
	MessySuccesses       int                `json:"messy_successes"`
	UsefulFailures       int                `json:"useful_failures"`
	FailedShifts         int                `json:"failed_shifts"`
	ShiftsWithMistakes   int                `json:"shifts_with_mistakes"`
	MistakesPerShift     float64            `json:"mistakes_per_shift"`
	FalseSuccesses       int                `json:"false_successes"`
	TotalFindings        int                `json:"total_findings"`
	RepeatedMistakes     int                `json:"repeated_mistakes"`
	PhantomToolCalls     int                `json:"phantom_tool_calls"`
	FalseBlockers        int                `json:"false_blockers"`
	BoundaryViolations   int                `json:"boundary_violations"`
	StuckLoops           int                `json:"stuck_loops"`
	AverageCostUSD       float64            `json:"average_cost_usd"`
	AverageDurationMS    int64              `json:"average_duration_ms"`
	MistakesByCode       map[string]int     `json:"mistakes_by_code"`
	ToolCalls            map[string]int     `json:"tool_calls"`
	UnavailableToolCalls map[string]int     `json:"unavailable_tool_calls"`
	SeverityCounts       map[Severity]int   `json:"severity_counts"`
	StrengthSignals      map[string]float64 `json:"strength_signals"`
}
