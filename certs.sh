#!/bin/bash

set -eo pipefail

mkdir -p /www/certs
chown root:www-data /www/certs

echo "---> pulling certs from etcd ($ETCD_HOST)"
certs=$(curl -s -L "http://$ETCD_HOST/v2/keys/www/certs" | jq ".node.nodes[].key" | xargs)
for cert in $certs; do
  echo "---> pulling cert $cert"
  curl -s -L "http://$ETCD_HOST/v2/keys/$cert" | jq ".node.value" | xargs echo -e > $cert
  chown root:www-data $cert
  chmod 0640 $cert
done