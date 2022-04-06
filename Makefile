all:	platform-all	compress-all

platform-all:	platform-windows	platform-windows-32	platform-windows-arm64	platform-darwin	platform-darwin-arm64

platform-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o pkg/grm-windows-64/grm.exe

platform-windows-32:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -o pkg/grm-windows-32/grm.exe

platform-windows-arm64:
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build -o	pkg/grm-windows-arm64/grm.exe

platform-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  -o pkg/grm-darwin-64/grm

platform-darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build  -o pkg/grm-darwin-arm64/grm

compress-all:	compress-windows	compress-windows-32	compress-windows-arm64	compress-darwin	compress-darwin-arm64

compress-windows:	platform-windows
	tar -zcvf pkg/grm-windows-64.tar.gz	pkg/grm-windows-64/grm.exe
	rm	-rf	pkg/grm-windows-64

compress-windows-32:	platform-windows-32
	tar -zcvf pkg/grm-windows-32.tar.gz	pkg/grm-windows-32/grm.exe
	rm -rf	pkg/grm-windows-32

compress-windows-arm64:	platform-windows-arm64
	tar -zcvf pkg/grm-windows-arm64.tar.gz	pkg/grm-windows-arm64/grm.exe
	rm -rf	pkg/grm-windows-arm64

compress-darwin:	platform-darwin
	tar -zcvf pkg/grm-darwin-64.tar.gz	pkg/grm-darwin-64/grm
	rm	-rf	pkg/grm-darwin-64

compress-darwin-arm64:	platform-darwin-arm64
	tar -zcvf pkg/grm-darwin-arm64.tar.gz	pkg/grm-darwin-arm64/grm
	rm	-rf	pkg/grm-darwin-arm64