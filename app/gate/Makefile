ENVFILE=$(CURDIR)/.env

all:

run:
	set -a; [ -f ${ENVFILE} ] && . ${ENVFILE}; set +a; go run main.go
