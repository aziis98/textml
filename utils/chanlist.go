package utils

// ChanList is a channel paired with a slice of all received items
type ChannelList[T any] struct {
	send    chan T
	receive chan T
	history []T
}

func NewChanList[T any](capacity ...int) *ChannelList[T] {
	cap := 0
	if len(capacity) > 0 {
		cap = capacity[0]
	}

	c := &ChannelList[T]{
		send:    make(chan T, cap),
		receive: make(chan T, cap),
		history: []T{},
	}

	go c.run()
	return c
}

func (c *ChannelList[T]) run() {
	for v := range c.send {
		c.history = append(c.history, v)
		c.receive <- v
	}
}

func (c *ChannelList[T]) Send() chan<- T {
	return c.send
}

func (c *ChannelList[T]) Receive() <-chan T {
	return c.receive
}
