package types

type Change struct {
	BeforeHTML string `json:"before_html"`
	AfterHTML  string `json:"after_html"`
	Changeset  string `json:"changeset"`
	Identifier string `json:"identifier"`
}

type Patch struct {
	Filename string   `json:"filename"`
	Changes  []Change `json:"changes"`
	Action   string   `json:"action"`
}

type FixedIssue struct {
	IssueCode string `json:"issue_code"`
	Count     int    `json:"count"`
}

type IssuesFixed struct {
	Filename    string       `json:"filename"`
	FixedIssues []FixedIssue `json:"fixed_issues"`
}

type AutofixReport struct {
	CodeDir      string        `json:"code_directory,omitempty"`
	ChangedFiles []string      `json:"changed_files,omitempty"`
	IssuesFixed  []IssuesFixed `json:"issues_fixed"`
	Metrics      []Metric      `json:"metrics,omitempty"`
	Patches      []Patch       `json:"patches"`
	Errors       []Error       `json:"errors"`
	ExtraData    interface{}   `json:"extra_data"`
}

type AutofixResult struct {
	RunID    string        `json:"run_id"`
	Status   Status        `json:"status"`
	CheckSeq string        `json:"check_seq"`
	Report   AutofixReport `json:"report"`
}

type AutofixResultCeleryTask struct {
	ID      string        `json:"id"`
	Task    string        `json:"task"`
	KWArgs  AutofixResult `json:"kwargs"`
	Retries int           `json:"retries"`
}

// Issues to be autofixed
//
//proteus:generate
type AutofixIssue struct {
	IssueCode   string          `json:"issue_code"`
	Occurrences []IssueLocation `json:"occurrences"`
}

//proteus:generate
type AutofixConfig struct {
	Issues []AutofixIssue `json:"issues"`
	Meta   interface{}    `json:"meta"`
}

//proteus:generate
type MarvinAutofixConfig struct {
	RunID             string `toml:"runID"`
	AnalyzerShortcode string `toml:"analyzerShortcode"`
	AutofixerCommand  string `toml:"autofixerCommand"`
	CheckoutOID       string `toml:"checkoutOID"`
	AutofixIssues     string `toml:"autofix_issues"`
}
