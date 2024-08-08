package baiaStructs

type DBMessage struct {
	Content   []OutputMessage `json:"content" bson:"content"`
	Role      string          `json:"role" bson:"role"`
	Timestamp int64           `json:"timestamp" bson:"timestamp"`
}
type Conversation struct {
	UserID   string      `json:"userID" bson:"userID"`
	ID       string      `json:"ID" bson:"ID"`
	IsActive bool        `json:"isActive" bson:"isActive"`
	Messages []DBMessage `json:"messages" bson:"messages"`
}
