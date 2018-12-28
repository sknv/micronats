ENVFILE=$(CURDIR)/.env

all:

gen:
	go generate ./...

run:
	set -a; [ -f ${ENVFILE} ] && . ${ENVFILE}; set +a; go run main.go
