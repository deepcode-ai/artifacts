package types

type DiffMeta struct {
	Additions [][]int `json:"additions"`
	Deletions [][]int `json:"deletions"`
}

type Coordinate struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

type Position struct {
	Begin Coordinate `json:"begin"`
	End   Coordinate `json:"end"`
}

type Location struct {
	Path     string   `json:"path"`
	Position Position `json:"position"`
}

type SourceCode struct {
	Rendered []byte `json:"rendered"`
}

type ProcessedData struct {
	SourceCode SourceCode `json:"source_code,omitempty"`
}

type Issue struct {
	IssueCode     string        `json:"issue_code"`
	IssueText     string        `json:"issue_text"`
	Location      Location      `json:"location"`
	ProcessedData ProcessedData `json:"processed_data,omitempty"`
	Identifier    string        `json:"identifier"`
}

// Location of an issue
type IssueLocation struct {
	Path     string   `json:"path"`
	Position Position `json:"position"`
}

type Namespace struct {
	Key      string                 `json:"key"`
	Value    float64                `json:"value"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

type Metric struct {
	MetricCode string      `json:"metric_code"`
	Namespaces []Namespace `json:"namespaces"`
}

type AnalysisError struct {
	HMessage string `json:"hmessage"`
	Level    int    `json:"level"`
}

type CommitDiffMeta map[string]DiffMeta

type FileMeta struct {
	IfAll      bool                      `json:"if_all"`
	Deleted    []string                  `json:"deleted"`
	Renamed    []string                  `json:"renamed"`
	Modified   []string                  `json:"modified"`
	Added      []string                  `json:"added"`
	DiffMeta   map[string]DiffMeta       `json:"diff_meta,omitempty"`
	PRDiffMeta map[string]CommitDiffMeta `json:"pr_diff_meta,omitempty"`
}

type AnalysisReport struct {
	Issues    []Issue         `json:"issues"`
	Metrics   []Metric        `json:"metrics,omitempty"`
	IsPassed  bool            `json:"is_passed"`
	Errors    []AnalysisError `json:"errors"`
	FileMeta  FileMeta        `json:"file_meta"`
	ExtraData interface{}     `json:"extra_data"`
}

type AnalysisResult struct {
	RunID    string         `json:"run_id"`
	Status   Status         `json:"status"`
	CheckSeq string         `json:"check_seq"`
	Report   AnalysisReport `json:"report"`
}

type AnalysisResultCeleryTask struct {
	ID      string         `json:"id"`
	Task    string         `json:"task"`
	KWArgs  AnalysisResult `json:"kwargs"`
	Retries int            `json:"retries"`
}

type CancelCheckResult struct {
	RunID  string `json:"run_id"`
	Status Status `json:"status"`
}

type CancelCheckResultCeleryTask struct {
	ID      string            `json:"id"`
	Task    string            `json:"task"`
	KWArgs  CancelCheckResult `json:"kwargs"`
	Retries int               `json:"retries"`
}

//proteus:generate
type MarvinAnalysisConfig struct {
	RunID                      string           `toml:"runID"`
	RunSerial                  string           `toml:"runSerial"`
	CheckSeq                   string           `toml:"checkSeq"`
	AnalyzerShortcode          string           `toml:"analyzerShortcode"`
	AnalyzerCommand            string           `toml:"analyzerCommand"`
	AnalyzerType               string           `toml:"analyzerType"`
	BaseOID                    string           `toml:"baseOID"`
	CheckoutOID                string           `toml:"checkoutOID"`
	Repository                 string           `toml:"repository"`
	IsFullRun                  bool             `toml:"is_full_run"`
	IsForDefaultAnalysisBranch bool             `toml:"isForDefaultAnalysisBranch"`
	DSConfigUpdated            bool             `toml:"dsConfigUpdated"`
	Processors                 []string         `toml:"processors"`
	DiffMetaCommits            []DiffMetaCommit `toml:"diffMetaCommits"`
}

//proteus:generate
type AnalysisStateInfo struct {
	IfAllFiles bool `json:"if_all_files"`
}
