### Before we GO

Dear Simply Wall Street Devs, initially I wanted to write this service in Java, but since I am experienced in writing 
REST and gRPC services in GO, I decided to go for GO.

### Getting started

Build the Docker image
```
$ docker build -t simply-ws-dev .
```

### Start the service

```
docker run -it --rm -p 8080:8080 simply-ws-dev
```

### Try the service

```
$ curl -X POST http://localhost:8080/entities \
    -H 'Content-Type: application/json' \
    -d '{
    "include_last_prices": true,
	"number_of_last_prices": 3,
	"page":               1,
	"entities_per_page":    2
}'
```
