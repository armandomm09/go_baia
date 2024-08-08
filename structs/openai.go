package baiaStructs

type Message struct {
	Response   string `json:"response"`
	AfterOrder bool   `json:"afterOrder,omitempty"`
	IsImage    bool   `json:"isImage"`
}

type Platillo struct {
	ID               int     `json:"id"`
	NombrePlatillo   string  `json:"nombre_platillo"`
	PrecioPorCadaUno float64 `json:"precio_por_cada_uno"`
	Cantidad         int     `json:"cantidad"`
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
