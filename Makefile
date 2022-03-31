platform-all: platform-windows platform-windows-32 platform-windows-arm64 platform-darwin platform-darwin-arm64

platform-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o pkg/gonrm-windows-64/gonrm.exe

platform-windows-32:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o pkg/gonrm-windows-32/gonrm.exe

platform-windows-arm64:
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o pkg/gonrm-windows-arm64/gonrm.exe

platform-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  -o pkg/gonrm-darwin-64/gonrm

platform-darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build  -o pkg/gonrm-darwin-arm64/gonrm
