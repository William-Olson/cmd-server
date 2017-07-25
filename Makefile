all: kill lb-pull server
	docker-compose up -d
	@echo done

server: version
	docker build -t willko/version-server:latest ./

lb-pull: FORCE
	docker pull dockercloud/haproxy:latest

version: FORCE
	./version

kill:
	docker-compose kill || true
	docker-compose rm -f || true

FORCE:

