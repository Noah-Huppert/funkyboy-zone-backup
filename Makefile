.PHONY: dev

DEV_SCRIPT ?= ./dev/dev.sh

# dev starts Prometheus and Prometheus Push Gateway Docker containers
dev:
	${DEV_SCRIPT}
