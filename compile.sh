CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -race
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -C cli -a -race
