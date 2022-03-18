package interfaces

type URLService interface {
	ShortenURL(url string) string
	ExpandURL(key string) (string, error)
}
