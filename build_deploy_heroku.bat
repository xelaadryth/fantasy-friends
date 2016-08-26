::Build statically-linked binary for alpine linux
set GOOS=linux
set GOARCH=amd64
go build -o alpine_binary .

deploy_heroku.bat
