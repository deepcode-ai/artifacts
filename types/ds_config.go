package types

// DSConfig is the struct for .deepcode.toml file
type Analyzer struct {
	Name                string      `toml:"name" json:"name"`
	RuntimeVersion      string      `toml:"runtime_version,omitempty" json:"runtime_version,omitempty"`
	Enabled             bool        `toml:"enabled" json:"enabled"`
	DependencyFilePaths []string    `toml:"dependency_file_paths,omitempty" json:"dependency_file_paths,omitempty"`
	Meta                interface{} `toml:"meta,omitempty" json:"meta,omitempty"`
	Thresholds          interface{} `toml:"thresholds,omitempty" json:"thresholds,omitempty"`
}

type Transformer struct {
	Name    string `toml:"name" json:"name"`
	Enabled bool   `toml:"enabled" json:"enabled"`
}

type DSConfig struct {
	Version         int           `toml:"version" json:"version"`
	ExcludePatterns []string      `toml:"exclude_patterns,omitempty" json:"exclude_patterns,omitempty"`
	TestPatterns    []string      `toml:"test_patterns,omitempty" json:"test_patterns,omitempty"`
	Analyzers       []Analyzer    `toml:"analyzers,omitempty" json:"analyzers,omitempty"`
	Transformers    []Transformer `toml:"transformers,omitempty" json:"transformers,omitempty"`
}

//proteus:generate
type AnalysisConfig struct {
	Files           []string    `json:"files"`
	ExcludePatterns []string    `json:"exclude_patterns"`
	ExcludeFiles    []string    `json:"exclude_files"`
	TestFiles       []string    `json:"test_files"`
	TestPatterns    []string    `json:"test_patterns"`
	AnalyzerMeta    interface{} `json:"analyzer_meta"`
}

//proteus:generate
type IDEConfig struct {
	IsIDE bool `json:"is_ide"`
}
