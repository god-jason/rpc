package pico

type Connect struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Id       string `json:"id,omitempty"`
	Token    string `json:"token,omitempty"`
}

type ConnectAck struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

type Disconnect struct {
	Reason string `json:"reason,omitempty"`
}

type Publish struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}

type PublishAck struct {
}

type Subscribe struct {
	Filters []string `json:"filters"`
}

type SubscribeAck struct {
	Filters map[string]bool `json:"filters"`
}

type Unsubscribe struct {
	Filters []string `json:"filters"`
}

type UnsubscribeAck struct {
	Filters map[string]bool `json:"filters"`
}
