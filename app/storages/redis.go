package storages

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"
)

//Redis struct {
type Redis struct {
	conn redis.Conn
}

//Init () {
func (s *Redis) Init() error {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	return nil
}

//Code () string {
func (s *Redis) Code() string {

	return "dummy"
}

//Save (url string) string {
func (s *Redis) Save(url string) string {

	return "dummy"
}

//Load (code string) (string, error) {
func (s *Redis) Load(code string) (string, error) {

	return string("dummyfdgfdgdfgfdgfdgfdgfd"), nil
}
