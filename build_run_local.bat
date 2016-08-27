::Build statically-linked binary for alpine linux
set GOOS=windows
set GOARCH=amd64
for /F "tokens=*" %%A in (.env) do set %%A
go build -o windows_binary.exe .

windows_binary.exe
