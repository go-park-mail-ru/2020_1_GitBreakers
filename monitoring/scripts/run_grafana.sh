#!/bin/bash

CURRENT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null 2>&1 && pwd)"

nohup \
  docker run \
    -p 3030:3000 \
    --name grafana \
    --add-host host.docker.internal:$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+') \
    --network bridge \
    -v /etc/grafana/provisioning:$CURRENT_DIR/grafana/provisioning \
    grafana/grafana &
