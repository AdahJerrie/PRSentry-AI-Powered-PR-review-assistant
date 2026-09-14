package review

type FileDiff struct {
	Path string `json:"path"`
	Diff string `json:"diff"`
}

type ReviewRequest struct {
	PRID  int        `json:"pr_id"`
	Repo  string     `json:"repo"`
	Files []FileDiff `json:"files"`
}
