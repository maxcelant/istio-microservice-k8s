### Users service

#### Running Locally

```bash
make run
```

#### Running with Docker
1. Build your docker image
```bash
docker build -t users:latest .
```

2. Run your docker container
```bash
docker run -p 8080:8080
```

3. Send requests to `/users` or `/user?={n}`

```bash
curl http://localhost:8080/api/user\?id\=9 | jq '.'
```

