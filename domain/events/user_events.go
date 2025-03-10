package events

type KafkaUserCreatedEvent struct {
	Event string `json:"event"`
	Type  string `json:"type"`
	Data  struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	} `json:"data"`
}

type KafkaUserFetchByIdEvent struct {
	Event string `json:"event"`
	Type  string `json:"type"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}

type KafkaUserReadAllEvent struct {
	Event string                 `json:"event"`
	Type  string                 `json:"type"`
	Data  map[string]interface{} `json:"data"` // Allows any JSON object
}
