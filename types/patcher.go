package types

// PatcherResultCeleryTask contains the structure of the response
// sent by the patcher job to asgard once it has completed.
type PatcherResultCeleryTask struct {
	ID      string        `json:"id"`
	Task    string        `json:"task"`
	KWArgs  PatcherResult `json:"kwargs"`
	Retries int           `json:"retries"`
}

// PatcherResult represents the patcher run result containing the run id and status of the
// patcher run.
type PatcherResult struct {
	RunID     string `json:"run_id"`
	CommitSHA string `json:"commit_sha"`
	Status    Status `json:"status"`
}
