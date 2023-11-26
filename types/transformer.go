package types

// Transformers types.
type TransformerReport struct {
	CodeDir       string   `json:"code_directory,omitempty"`
	ChangedFiles  []string `json:"changed_files,omitempty"`
	CommitSHA     string   `json:"commit_sha"`
	CommitCreated bool     `json:"commit_created"`
	Errors        []Error  `json:"errors"`
	Patches       []Patch  `json:"patches"`
}

type TransformerResult struct {
	RunID  string            `json:"run_id"`
	Status Status            `json:"status"`
	Report TransformerReport `json:"report"`
}

type TransformerResultCeleryTask struct {
	ID      string            `json:"id"`
	Task    string            `json:"task"`
	KWArgs  TransformerResult `json:"kwargs"`
	Retries int               `json:"retries"`
}

type MarvinTransformerConfig struct {
	RunID              string      `toml:"runID"`
	BaseOID            string      `toml:"baseOID"`
	CheckoutOID        string      `toml:"checkoutOID"`
	TransformerCommand string      `toml:"transformerCommand"`
	TransformerTools   []string    `toml:"transformerTools"`
	DSConfigUpdated    bool        `toml:"dsConfigUpdated"`
	PatchCommit        PatchCommit `toml:"patch_commit"`
}

type TransformerConfig struct {
	ExcludePatterns []string `json:"exclude_patterns"`
	ExcludeFiles    []string `json:"exclude_files"`
	Files           []string `json:"files"`
	Tools           []string `json:"tools"`
}
