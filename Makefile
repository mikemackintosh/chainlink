all: test build

test:
	go test -v ./...

build:
	go build -o bin/chainlink cmd/chainlink/*.go

start:
	go run cmd/chainlink/*.go


config:
	sudo networksetup -setdnsservers Wi-Fi 127.0.0.1
	sudo networksetup -setdnsservers Ethernet 127.0.0.1
	sudo chflags schg /Library/Preferences/SystemConfiguration/preferences.plist
