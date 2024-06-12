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
$ docker run -it --rm -p 8080:8080 simply-ws-dev
```

### Try the service

```
$ curl -X POST http://localhost:8080/entities \
    -H 'Content-Type: application/json' \
    -H 'Authorization: Bearer 8af4cc4fbf1eb641b14aeb7235bc7509' \
    -d '{
    "include_last_prices": true,
	"number_of_last_prices": 1,
	"page":               1,
	"entities_per_page":    2
}'
```

### Authorization

The service requires an authorization token `8af4cc4fbf1eb641b14aeb7235bc7509` for the purposes of basic authorization.
No other tokens are issued by the service.
