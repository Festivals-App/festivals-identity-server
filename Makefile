# Makefile for festivals-identity-server

VERSION=development
DATE=$(shell date +"%d-%m-%Y-%H-%M")
REF=refs/tags/development
export

build:
	go build -ldflags="-X 'github.com/Festivals-App/festivals-identity-server/server/status.ServerVersion=$(VERSION)' -X 'github.com/Festivals-App/festivals-identity-server/server/status.BuildTime=$(DATE)' -X 'github.com/Festivals-App/festivals-identity-server/server/status.GitRef=$(REF)'" -o festivals-identity-server main.go

install:
	cp festivals-identity-server /usr/local/bin/festivals-identity-server
	cp config_template.toml /etc/festivals-identity-server.conf
	cp operation/service_template.service /etc/systemd/system/festivals-identity-server.service

update:
	systemctl stop festivals-identity-server
	cp festivals-identity-server /usr/local/bin/festivals-identity-server
	systemctl start festivals-identity-server

uninstall:
	rm /usr/local/bin/festivals-identity-server
	rm /etc/festivals-identity-server.conf
	rm /etc/systemd/system/festivals-identity-server.service

run:
	./festivals-identity-server

stop: 
	killall festivals-identity-server

clean:
	rm -r festivals-identity-server
