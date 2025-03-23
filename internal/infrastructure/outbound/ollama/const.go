package ollama

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Messages []Message

type Response struct {
	Message Message `json:"message"`
}

type Options struct {
	Temperature float64 `json:"temperature"`
}
