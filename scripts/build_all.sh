#!/bin/bash

docker build -t cubepanion-api-base . -f docker/Dockerfile.base

docker build -t chests-service . -f docker/Dockerfile.chest-service &
pids+=($!)

docker build -t cubesocket . -f docker/Dockerfile.cubesocket &
pids+=($!)

docker build -t games-service . -f docker/Dockerfile.games-service &
pids+=($!)

docker build -t leaderboard-service . -f docker/Dockerfile.leaderboard-service &
pids+=($!)

docker build -t maps-service . -f docker/Dockerfile.maps-service &
pids+=($!)

# Wait for all background processes and check for errors
for pid in "${pids[@]}"; do
    wait "$pid"
    if [ $? -ne 0 ]; then
        echo "Error: Process $pid failed."
        exit 1
    fi
done

echo "All Docker builds completed successfully."
