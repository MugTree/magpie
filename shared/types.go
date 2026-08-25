package shared

type Feed struct {
	Url                    string
	CSSSelectorContainer   string
	CSSSelectorStart       string
	CSSSelectorStop        string
	HTMLExtractionStrategy string
}

type PageScrapeParams struct {
	Link           string
	Container      string
	ClipStartPoint string
	ClipEndPoint   string
	Strategy       string
}
