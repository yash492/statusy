package http

type Client interface {
	GetBytes(url string) ([]byte, error)
}
