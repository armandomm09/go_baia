package baiaStructs

type DBAssistantMessage struct {
	Content   []OutputMessage `json:"content" bson:"content"`
	Role      string          `json:"role" bson:"role"`
	Timestamp int64           `json:"timestamp" bson:"timestamp"`
}

type Conversation struct {
	UserID   string               `json:"userID" bson:"userID"`
	ID       string               `json:"ID" bson:"ID"`
	IsActive bool                 `json:"isActive" bson:"isActive"`
	Messages []DBAssistantMessage `json:"messages" bson:"messages"`
}

type Location struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
}

type Address struct {
	Street      string `json:"street" bson:"street"`
	Number      int32  `json:"number" bson:"number"`
	Suburb      string `json:"suburb" bson:"suburb"`
	Description string `json:"description" bson:"description"`
}

type FinalOrder struct {
	ID               string   `json:"ID" bson:"ID"`
	UserID           string   `json:"userID" bson:"userID"`
	CreationDate     int64    `json:"creationDate" bson:"creationDate"`
	State            string   `json:"state" bson:"state"`
	Order            Order    `json:"order" bson:"order"`
	Total            float32  `json:"total" bson:"total"`
	DeliveryLocation Location `json:"deliveryLocation" bson:"deliveryLocation"`
	DeliveryAddress  Address  `json:"deliveryAddress" bson:"deliveryAddress"`
	DeliveryDate     int64    `json:"deliveryDate" bson:"deliveryDate"`
	Comments         string   `json:"comments" bson:"comments"`
}
