package games

type Game interface {
	Name() string
	Play()
}
