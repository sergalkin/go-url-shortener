package interfaces

type URLExpand interface {
	ExpandURL(key string) (string, error)
}
