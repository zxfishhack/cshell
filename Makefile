all: darwin-arm64 darwin-amd64
	rm -rf cshell.app

common:
	go generate main.go
	rm -rf cshell.app
	mkdir -p cshell.app/Contents/MacOS
	mkdir -p cshell.app/Contents/Resources
	cp Info.plist cshell.app/Contents
	cp icon.icns cshell.app/Contents/Resources

darwin-amd64: common
	CGO_ENABLED=1 GOARCH=amd64 go build -o cshell.app/Contents/MacOS/cshell .
	zip -r cshell.amd64.zip cshell.app

darwin-arm64: common
	GOARCH=arm64 go build -o cshell.app/Contents/MacOS/cshell .
	zip -r cshell.arm64.zip cshell.app