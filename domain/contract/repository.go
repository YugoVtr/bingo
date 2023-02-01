package contract

type BingoRepository interface {
	AddHistoric(int)
	ListenHistoric(chan<- int)
}
