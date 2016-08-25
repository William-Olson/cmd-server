all: lb-pull server
	@echo done

server: FORCE
	docker build -t willko/version-server:latest ./server

lb-pull: FORCE
	docker pull dockercloud/haproxy:latest

FORCE:

