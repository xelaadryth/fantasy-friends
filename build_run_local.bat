::Build statically-linked binary for Windows
del windows_binary.exe
set GOOS=windows
set GOARCH=amd64
for /F "tokens=*" %%A in (.env) do set %%A
::Override from release to debug mode
set GIN_MODE=debug
go build -o windows_binary.exe .

windows_binary.exe
