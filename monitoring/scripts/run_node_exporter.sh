#!/bin/bash

nohup \
  docker run --rm \
    -p 9100:9100 \
    --add-host host.docker.internal:$(ip -4 addr show docker0 | grep -Po 'inet \K[\d.]+') \
    --name node_exporter \
    --network bridge \
    prom/node-exporter &
