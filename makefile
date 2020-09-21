#!make
# DON'T EXTEND THE SCOPE OF THIS FILE FOR DEVELOPMENT CONCERNS
# USE THIS FILE FOR DEPLOYMENT ONLY

include .env

build: 
	@echo "\n[ building production image ]"
	docker build --target prod --tag devpies/client-api ./api

login: 
	@echo "\n[ logging into private registry ]"
	cat ./secrets/registry_pass | docker login --username `cat ./secrets/registry_user` --password-stdin

publish:
	@echo "\n[ publishing production grade image ]"
	docker push devpies/client-api

deploy:
	@echo "\n[ deploying production stack ]"
	@cat ./startup
	@docker stack deploy -c docker-stack.yml --with-registry-auth devpie

metrics: 
	@echo "\n[ enabling docker engine metrics ]"
	./deploy/enable-monitoring.sh

secrets: 
	@echo "\n[ creating swarm secrets ]"
	./deploy/create-secrets.sh

server:
	@echo "\n[ creating server ]"
	./deploy/create-server.sh

server-d:
	@echo "\n[ destroying server ]"
	./deploy/destroy-server.sh

swarm:
	@echo "\n[ create single node swarm ]"
	./deploy/create-swarm.sh

.PHONY: build 
.PHONY: login
.PHONY: publish
.PHONY: deploy
.PHONY: metrics
.PHONY: secrets
.PHONY: servers
.PHONY: servers-d
.PHONY: swarm
