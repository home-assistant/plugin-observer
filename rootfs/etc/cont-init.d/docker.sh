#!/usr/bin/with-contenv bashio
# ==============================================================================
# Setup docker client settings
# ==============================================================================

echo "unix:///var/run/docker.sock" > /var/run/s6/container_environment/DOCKER_HOST
