compile:
	echo "Compiling for Darwin and Linux"
	GOOS=linux GOARCH=arm64 go build -o release/c.linux
	go build -o release/c.darwin