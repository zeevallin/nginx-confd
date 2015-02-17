#!/bin/bash

set -eo pipefail

export ETCD_HOST="${ETCD_PORT_4001_TCP_ADDR:-127.0.0.1}:${ETCD_PORT_4001_TCP_PORT:-4001}"

exec "$@"