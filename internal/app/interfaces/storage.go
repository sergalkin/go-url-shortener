package interfaces

type Storage interface {
	Store(key string, url string)
	Get(key string) (string, bool)
}
