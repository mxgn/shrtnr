# shrtnr


### vendoring dependencies:

```sh
go get -u github.com/kardianos/govendor
```

Make sure you're in the correct project directory

```sh
govendor init
govendor add +external
```

## run only one detached service:
```sh
docker-compose up -d --no-deps pgdb
```