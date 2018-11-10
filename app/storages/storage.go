// Package storages allows multiple implementation on how to store short URLs.
package storages

//IStorage interface {
type IStorage interface {
	Save(string) string
	Load(string) (string, error)
}
