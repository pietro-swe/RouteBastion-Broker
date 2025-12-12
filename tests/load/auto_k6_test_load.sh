#!/usr/bin/env bash
set -euo pipefail

USERS_LIST=(50 75 100 200 300)
INSTANCES_LIST=(1 2 3)
SPAWN_RATE=10
RUN_TIME="5m"
RUNS=10
HOST="http://localhost:8000"

trap cleanup_and_exit INT

cleanup_and_exit() {
  echo
  echo "Caught interrupt signal. Cleaning up..."
  docker compose -f ../../docker-compose.yml down -v
  echo "Cleanup done. Exiting."
  exit 0
}

generate_kong_targets() {
  count=$1
  out=""
  for i in $(seq 1 "$count"); do
    out+="      - target: route-bastion-broker-api-broker-${i}:8080\n        weight: 100\n"
  done
  echo -e "$out"
}

check_if_node_installed() {
  if ! command -v node &> /dev/null; then
    echo "Node.js is not installed. Please install Node.js to proceed."
    exit 1
  fi
}

install_k6_if_missing() {
  if ! command -v k6 &> /dev/null; then
    echo "k6 not found, installing..."
    npm install -g k6
  fi
}

check_if_node_installed

install_k6_if_missing

echo "=== Starting load tests ==="

for INSTANCES in "${INSTANCES_LIST[@]}"; do
  echo
  echo "=== Preparing environment for $INSTANCES instance(s) ==="

  export KONG_TARGETS="$(generate_kong_targets "$INSTANCES")"
  envsubst < kong/config.tpl.yaml > ../../docker/kong/config.yaml

  echo "Restarting Docker environment..."
  docker compose -f ../../docker-compose.yml down -v
  docker compose -f ../../docker-compose.yml up -d --scale api-broker="$INSTANCES"

  for USERS in "${USERS_LIST[@]}"; do
    for i in $(seq 1 "$RUNS"); do
      TIMESTAMP=$(date +%Y-%m-%d-%Hh%M)
      CSV_NAME="${TIMESTAMP}_run${i}_${USERS}users_${SPAWN_RATE}rate_${RUN_TIME}min_${INSTANCES}instances"

      # rampa = USERS / SPAWN_RATE
      RAMP_SEC=$(( USERS / SPAWN_RATE ))

      echo
      echo "[$(date '+%H:%M:%S')] Running test $i/$RUNS with $USERS users against $INSTANCES instance(s)..."

      k6 run test.js \
        --env HOST="$HOST" \
        --env VUS="$USERS" \
        --env RUN_TIME="$RUN_TIME" \
        --env SPAWN_RATE="$SPAWN_RATE" \
        --env RAMP="$RAMP_SEC" \
        --env FILE_NAME="results/${CSV_NAME}.csv"

      echo "[$(date '+%H:%M:%S')] Test $i/$RUNS completed."
    done
  done

  echo "[$(date '+%H:%M:%S')] Done with $INSTANCES instance(s)."
done

echo
echo "🎉 All load tests completed successfully!"
