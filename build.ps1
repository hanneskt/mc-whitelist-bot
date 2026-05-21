$ErrorActionPreference = "Stop"

$BinaryName = "whitelistbot.bin"

$env:CGO_ENABLED = "0"
$env:GOOS = "linux"
$env:GOARCH = "amd64"

go build -ldflags="-s -w" -o $BinaryName main.go
