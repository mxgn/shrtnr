package storages

type IStorage interface {
	AddLongUrl(string) (string, error)
	GetLongUrl(string) (string, error)
}
