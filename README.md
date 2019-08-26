# go-url-shortener
[![Build Status](https://travis-ci.org/allisson/go-url-shortener.svg)](https://travis-ci.org/allisson/go-url-shortener) [![Go Report Card](https://goreportcard.com/badge/github.com/allisson/go-url-shortener)](https://goreportcard.com/report/github.com/allisson/go-url-shortener)

Golang url shortener using clean code.

Features:

* Multiple storage engine: Redis and MongoDB.
* Multiple message format: JSON and Msgpack.

## Run Tests

1. Set envvars

```bash
export MONGODB_URL=mongodb://localhost/shortener
export MONGODB_TIMEOUT=30
export MONGODB_DATABASE=shortener
export REDIS_URL=redis://localhost:6379
```

2. Run tests

```bash
make test
```

## Run the server

1. Set envvars

For redis storage:

```bash
export PORT=3000
export STORAGE_ENGINE=redis
export REDIS_URL=redis://localhost:6379
```

For mongo storage:

```bash
export PORT=3000
export STORAGE_ENGINE=mongo
export MONGODB_URL=mongodb://localhost/shortener
export MONGODB_TIMEOUT=30
export MONGODB_DATABASE=shortener
```

2. Run the server

With binary:

```bash
make build
./shortener
```

Without binary:

```bash
make run
```

## How to use

Create new short url:

```bash
curl -X POST http://localhost:3000 \
  -H 'Content-Type: application/json' \
  -d '{
	"url": "https://github.com/allisson"
}'
```

The response:

```json
{
    "code": "YzaFuvFWg",
    "url": "https://github.com/allisson",
    "created_at": 1566827456
}
```

Get redirected:

```bash
curl -X GET http://localhost:3000/YzaFuvFWg
```

The response:

```bash
<a href="https://github.com/allisson">Moved Permanently</a>.
```

If you want to use msgpack format, change `Content-Type` header to `application/x-msgpack`.
