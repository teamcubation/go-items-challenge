package analysis

type FileAnalysisResult struct {
	IsValid   bool
	Score     int
	Evidences []Evidence
}

type Evidence struct {
	File string `json:"file"`
	Line int    `json:"line"`
}
