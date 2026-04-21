
cd /d %~dp0


mkdir build

go build -ldflags="-s -w" -trimpath -o trimpath -o build\VRC-GoWorldPage.exe main.go
