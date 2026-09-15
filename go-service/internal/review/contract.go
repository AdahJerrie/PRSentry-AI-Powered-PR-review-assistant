package review

type FileDiff struct {
	Path string `json:"path"`
	Diff string `json:"diff"`
}

type ReviewRequest struct {
	PRID  int        `json:"pr_id"`
	Files []FileDiff `json:"files"`
}
