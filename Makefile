all: kill lb-pull server
	docker-compose up -d
	@echo done

server: FORCE
	docker build -t willko/version-server:latest ./server

lb-pull: FORCE
	docker pull dockercloud/haproxy:latest

kill:
	docker-compose kill || true
	docker-compose rm -f || true

FORCE:

