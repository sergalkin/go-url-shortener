package interfaces

type URLService interface {
	ShortenURL(url string) (string, error)
	ExpandURL(key string) (string, error)
}
