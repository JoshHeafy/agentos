package gemini

type Payload struct {
	Content Content `json:"content"`
}

type Content struct {
	Parts Parts `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type Parts []Part
