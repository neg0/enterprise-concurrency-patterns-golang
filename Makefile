.PHONY: help     # Generate list of targets with descriptions
help:
	@echo "\n"
	@grep '^.PHONY: .* #' Makefile | sed 's/\.PHONY: \(.*\) # \(.*\)/\1 \2/' | expand -t20

.PHONY: up       # Creates container for Golang
up:
	docker-compose -f docker-compose.yml up -d

.PHONY: down     # It shuts down the running containers
down:
	docker-compose -f docker-compose.yml down

.PHONY: go      # Enters the Golang Container
go:
	docker-compose -f docker-compose.yml exec golang bash