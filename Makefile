all:	platform-all	upx-all	compress-all	msi

# opt

GO_FLAGS += "-ldflags=-s -w"

platform-all:
	@$(MAKE)	--no-print-directory	\
			platform-windows	platform-windows-32	platform-windows-arm64	\
			platform-darwin	platform-darwin-arm64	\
			platform-linux	platform-linux-32	platform-linux-arm64

platform-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build	${GO_FLAGS}	-o build/grm-windows-64/grm.exe	./cmd/grm

platform-windows-32:
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build	${GO_FLAGS} -o build/grm-windows-32/grm.exe	./cmd/grm

platform-windows-arm64:
	CGO_ENABLED=0 GOOS=windows GOARCH=arm64 go build	${GO_FLAGS} -o	build/grm-windows-arm64/grm.exe	./cmd/grm

platform-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build	${GO_FLAGS} -o build/grm-darwin-64/grm	./cmd/grm

platform-darwin-arm64:
	CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build	${GO_FLAGS}	-o build/grm-darwin-arm64/grm	./cmd/grm

platform-linux:
	CGO_ENABLED=0	GOOS=linux	GOARCH=amd64	go	build	${GO_FLAGS}	-o	build/grm-linux-64/grm	./cmd/grm

platform-linux-32:
	CGO_ENABLED=0	GOOS=linux	GOARCH=386	go	build	${GO_FLAGS}	-o	build/grm-linux-32/grm	./cmd/grm

platform-linux-arm64:
	CGO_ENABLED=0	GOOS=linux	GOARCH=arm	go build	${GO_FLAGS}	-o	build/grm-linux-arm64/grm	./cmd/grm

upx:
	upx	--best	--lzma	build/grm-windows-64/grm.exe
	upx	--best	--lzma	build/grm-windows-32/grm.exe
	upx	--best	--lzma	build/grm-darwin-64/grm
	upx	--best	--lzma	build/grm-linux-64/grm
	upx --best	--lzma	build/grm-linux-32/grm


upx-all:	upx

compress-all:
	@$(MAKE)	--no-print-directory	\
		compress-windows	compress-windows-32	compress-windows-arm64	\
		compress-darwin	compress-darwin-arm64	\
		compress-linux	compress-linux-32	compress-linux-arm64

compress-windows:
	tar -zcvf build/grm-windows-64.tar.gz	-Cbuild/grm-windows-64/	grm.exe
	rm	-rf	build/grm-windows-64

compress-windows-32:
	tar -zcvf build/grm-windows-32.tar.gz	-Cbuild/grm-windows-32/	grm.exe
	rm -rf	build/grm-windows-32

compress-windows-arm64:
	tar -zcvf build/grm-windows-arm64.tar.gz	-Cbuild/grm-windows-arm64/	grm.exe
	rm -rf	build/grm-windows-arm64

compress-darwin:
	tar -zcvf build/grm-darwin-64.tar.gz	-Cbuild/grm-darwin-64/	grm
	rm	-rf	build/grm-darwin-64

compress-darwin-arm64:
	tar -zcvf build/grm-darwin-arm64.tar.gz	-Cbuild/grm-darwin-arm64/	grm
	rm	-rf	build/grm-darwin-arm64

compress-linux:
	tar -zcvf build/grm-linux-64.tar.gz	-Cbuild/grm-linux-64/	grm
	rm	-rf	build/grm-linux-64

compress-linux-32:
	tar -zcvf build/grm-linux-32.tar.gz	-Cbuild/grm-linux-32/	grm
	rm	-rf	build/grm-linux-32

compress-linux-arm64:
	tar -zcvf build/grm-linux-arm64.tar.gz	-Cbuild/grm-linux-arm64/	grm
	rm	-rf	build/grm-linux-arm64


msi:
ifeq	($(OS),Windows_NT)
	cd scripts
	node scripts/windows-msi.js
endif
