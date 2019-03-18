#!/usr/bin/env bash
#?
# dev.sh - Starts a Prometheus and Prometheus Push Gateway container
#
# USAGE
#
#	dev.sh
#
# BEHAVIOR
#
#	Starts a Prometheus container which is configured to scrape a Prometheus Push Gateway container. 
#
#	The script will not launch either container if it is already running.
#
#	The script will exit when you press enter. This will shutdown any containers.
#
#?

# {{{1 Exit on any error
set -e

# {{{1 Configuration
prog_dir=$(realpath $(dirname "$0"))
prom_container_tag="prom/prometheus:latest"
prom_container_name="mountain-backup-prometheus"
push_container_name="mountain-backup-push-gateway"

# {{{1 Check if docker exists on system
if ! which docker &> /dev/null; then
	echo "Error: Docker not installed on system, must be installed" >&2
	exit 1
fi

# {{{1 Launch containers
# {{{2 Cleanup function which stops any running containers
function cleanup() {
	code=0
	# Remove containers
	for name in "$prom_container_name" "$push_container_name"; do
		if docker ps | grep "$name" &> /dev/null; then
			if ! docker stop "$name"; then
				echo "Error: cleanup: Failed to stop $name container" >&2
				code=1
			fi
		fi
	done

	exit "$code"
}
trap cleanup EXIT

# {{{2 Prometheus
if ! docker ps | grep "$prom_container_name" &> /dev/null; then
	docker run \
		-t \
		--rm \
		--net host \
		--name "$prom_container_name" \
		-p 9090:9090 \
		-v "$prog_dir/prometheus-config":/etc/prometheus/ \
		"$prom_container_tag" \
			--config.file=/etc/prometheus/prometheus.yml \
			--web.listen-address="localhost:9090" \
			--web.enable-admin-api &
	
	if [[ "$?" != "0" ]]; then
		echo "Error: Failed to start Prometheus container" >&2
		exit 1
	fi
fi

# {{{2 Push Gateway
if ! docker ps | grep "$push_container_name" &> /dev/null; then
	docker run \
		-t \
		--rm \
		--net host \
		--name "$push_container_name" \
		-p 9091:9091 \
		prom/pushgateway &

	if [[ "$?" != "0" ]]; then
		echo "Error: Failed to start Prometheus Push Gateway container" >&2
		exit 1
	fi
fi

# {{{1 Exit on key press
read foo
