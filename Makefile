all:

# -----------------------------------------------------------------------------
# dep
# -----------------------------------------------------------------------------

dep:
	dep ensure -v

# -----------------------------------------------------------------------------
# go generate section
# -----------------------------------------------------------------------------

gen-service:
	$(MAKE) -C app/${SERVICE} gen

gen-auth:
	SERVICE=auth $(MAKE) gen-service

gen:
	go generate ./...

# -----------------------------------------------------------------------------
# running section
# -----------------------------------------------------------------------------

run-service:
	$(MAKE) -C app/${SERVICE} run

run-gate:
	SERVICE=gate $(MAKE) run-service

run-auth:
	SERVICE=auth $(MAKE) run-service

# -----------------------------------------------------------------------------
# docker section
# -----------------------------------------------------------------------------

start-container:
	docker-compose up -d ${CONTAINER}

stop-container:
	docker-compose stop ${CONTAINER}

start-nats:
	CONTAINER=nats $(MAKE) start-container

stop-nats:
	CONTAINER=nats $(MAKE) stop-container
