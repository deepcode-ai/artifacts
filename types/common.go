package types

type File struct {
	Path int8 `json:"path"`
	Q1   int8 `json:"q1"`
	Q3   int8 `json:"q2"`
	Q7   int8 `json:"q7"`
	Q15  int8 `json:"q15"`
	Q30  int8 `json:"q30"`
	Q60  int8 `json:"q60"`
	Q180 int8 `json:"a180"`
}

type BeaconResult struct {
	Files []File `json:"files"`
}

type BeaconStatusMsg struct {
	RunID     string       `json:"run_id"`
	ProjectID string       `json:"project_id"`
	Status    Status       `json:"status"`
	Result    BeaconResult `json:"result"`
}

type BeaconResultCeleryTask struct {
	ID      string          `json:"id"`
	Task    string          `json:"task"`
	KWArgs  BeaconStatusMsg `json:"kwargs"`
	Retries int             `json:"retries"`
}

type Error struct {
	HMessage string `json:"hmessage"`
	Level    int    `json:"level"`
}

type Status struct {
	Code     int    `json:"code"`
	HMessage string `json:"hmessage"`
	Err      string `json:"err"`
}
