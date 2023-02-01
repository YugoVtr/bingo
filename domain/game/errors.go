package game

const (
	errGameStarted = BingoError("game started")
)

type BingoError string

func (err BingoError) Error() string {
	return string(err)
}
