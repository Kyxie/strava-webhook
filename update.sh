#!/bin/bash
set -e

# Absolute path to docker compose binary
DOCKER_COMPOSE="/usr/bin/docker compose"

# Helper function: run command inside the "app" container
exec_in_app() {
  $DOCKER_COMPOSE exec -T app "$@"
}

echo "[Update.sh] Importing latest Strava activities..."
exec_in_app bin/console app:strava:import-data

echo "[Update.sh] Building files..."
exec_in_app bin/console app:strava:build-files

echo "[Update.sh] Update complete."
