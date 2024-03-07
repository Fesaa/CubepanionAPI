#!/bin/bash

docker build -t chests-service . -f Dockerfile.chest-service &
pids+=($!)

docker build -t cubesocket . -f Dockerfile.cubesocket &
pids+=($!)

docker build -t games-service . -f Dockerfile.games-service &
pids+=($!)

docker build -t leaderboard-service . -f Dockerfile.leaderboard-service &
pids+=($!)

docker build -t maps-service . -f Dockerfile.maps-service &
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
