package baiaStructs

type Message struct {
	Response   string `json:"response"`
	AfterOrder bool   `json:"afterOrder,omitempty"`
	IsImage    bool   `json:"isImage"`
}

type Platillo struct {
	ID           int     `json:"id"`
	ServiceName  string  `json:"serviceName"`
	UnitaryPrice float64 `json:"unitaryPrice"`
	Quantity     int     `json:"quantity"`
}

type Order struct {
	Order []Platillo `json:"orden"`
}

type GPTUnformattedResponse struct {
	Messages []Message  `json:"messages"`
	Order    []Platillo `json:"orden"`
}

type OutputMessage struct {
	Response string `json:"response"`
	IsImage  bool   `json:"isImage"`
}

type FinalGPTResponse struct {
	Messages []OutputMessage `json:"messages"`
}
