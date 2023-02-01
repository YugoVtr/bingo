package repository

func NewInMemory() *gameInMemory {
	return new(gameInMemory)
}

type gameInMemory struct {
	historic []int
}

func (s *gameInMemory) AddHistoric(h int) {
	s.historic = append(s.historic, h)
}

func (s *gameInMemory) ListenHistoric(ch chan<- int) {
	go func() {
		count := len(s.historic)
		for {
			if currentCount := len(s.historic); currentCount > count {
				ch <- s.historic[currentCount-1]
				count = currentCount
			}
		}
	}()
}
