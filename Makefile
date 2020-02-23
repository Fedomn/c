compile:
	echo "Compiling for Darwin and Linux"
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o release/c.linux
	go build -ldflags="-s -w" -o release/c.darwin

upx: compile
	upx release/c.linux
	upx release/c.darwin

test:
	ginkgo -v