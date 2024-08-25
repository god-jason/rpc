package pico

type pending struct {
	c chan *Pack
}

func newPending() *pending {
	return &pending{
		c: make(chan *Pack),
	}
}
