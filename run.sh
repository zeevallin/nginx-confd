#!/bin/bash

set -eo pipefail

echo "---> booting container with etcd ($ETCD_HOST)"

until confd -onetime -node $ETCD_HOST -config-file /etc/confd/conf.d/certs.toml; do
  echo "---> waiting for confd to update certificates"
  sleep 5
done

until confd -onetime -node $ETCD_HOST -config-file /etc/confd/conf.d/nginx.toml; do
  echo "---> waiting for confd to update nginx config"
  sleep 5
done

confd -interval 10 -node $ETCD_HOST -config-file /etc/confd/conf.d/certs.toml &
confd -interval 10 -node $ETCD_HOST -config-file /etc/confd/conf.d/nginx.toml &
echo "---> confd is listening for changes on etcd..."

# Start nginx
echo "---> starting nginx service..."
nginx