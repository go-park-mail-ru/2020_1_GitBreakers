#!/bin/bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

nohup \
    docker run \
      --add-host host.docker.internal:"$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+')" \
      -p 9090:9090 \
      --name prometheus \
      --network bridge \
      -v $CURRENT_DIR/prometheus:/etc/prometheus \
      prom/prometheus \
      --config.file=/etc/prometheus/prometheus.yml &
