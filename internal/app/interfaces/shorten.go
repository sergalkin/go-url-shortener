package interfaces

type URLShorten interface {
	ShortenURL(url string) (string, error)
}
