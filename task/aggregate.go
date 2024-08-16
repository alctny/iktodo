package task

type AggregateResult struct {
	Total    uint `json:"total"`
	Finished uint `json:"finished"`
	Unfinish uint `json:"unfinish"`
}
