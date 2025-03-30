# Makefile for festivals-identity-server

SHELL := bash
VERSION=development
DATE=$(shell date +"%d-%m-%Y-%H-%M")
REF=refs/tags/development
DEV_PATH_MAC=$(shell echo ~/Library/Containers/org.festivalsapp.project)
export

build:
	go build -ldflags="-X 'github.com/Festivals-App/festivals-identity-server/server/status.ServerVersion=$(VERSION)' -X 'github.com/Festivals-App/festivals-identity-server/server/status.BuildTime=$(DATE)' -X 'github.com/Festivals-App/festivals-identity-server/server/status.GitRef=$(REF)'" -o festivals-identity-server main.go

install:
	mkdir -p $(DEV_PATH_MAC)/usr/local/bin
	mkdir -p $(DEV_PATH_MAC)/etc
	mkdir -p $(DEV_PATH_MAC)/var/log
	mkdir -p $(DEV_PATH_MAC)/usr/local/festivals-identity-server
	
	cp operation/local/ca.crt  $(DEV_PATH_MAC)/usr/local/festivals-identity-server/ca.crt
	cp operation/local/server.crt  $(DEV_PATH_MAC)/usr/local/festivals-identity-server/server.crt
	cp operation/local/server.key  $(DEV_PATH_MAC)/usr/local/festivals-identity-server/server.key
	cp operation/local/authentication.publickey.pem  $(DEV_PATH_MAC)/usr/local/festivals-identity-server/authentication.publickey.pem
	cp operation/local/authentication.privatekey.pem  $(DEV_PATH_MAC)/usr/local/festivals-identity-server/authentication.privatekey.pem
	cp festivals-identity-server $(DEV_PATH_MAC)/usr/local/bin/festivals-identity-server
	chmod +x $(DEV_PATH_MAC)/usr/local/bin/festivals-identity-server
	cp operation/local/config_template_dev.toml $(DEV_PATH_MAC)/etc/festivals-identity-server.conf

run:
	./festivals-identity-server --container="$(DEV_PATH_MAC)"
	
run-env:
	$(DEV_PATH_MAC)/usr/local/bin/festivals-gateway --container="$(DEV_PATH_MAC)" &

stop-env:
	killall festivals-gateway

clean:
	rm -r festivals-identity-server
