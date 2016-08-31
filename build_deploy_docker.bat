::Build statically-linked binary for alpine linux
set GOOS=linux
set GOARCH=amd64
set GIN_MODE=debug
go build -o alpine_binary .

deploy_docker.bat
