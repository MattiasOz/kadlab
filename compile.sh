CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -C cli -a
