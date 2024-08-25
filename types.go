package pico

type Connect struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Auth struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}

type ConnectAck struct {
	Result bool  `json:"result"`
	Auth   *Auth `json:"auth,omitempty"`
}

type Disconnect struct {
	Reason string `json:"reason,omitempty"`
}

type Publish struct {
	Topic   string `json:"topic"`
	Message any    `json:"message"`
}

type Message struct {
	Topic   string `json:"topic"`
	Message any    `json:"message"`
}

type PublishAck struct {
	Topics map[string]bool `json:"topics"`
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
