# nginx-confd-reg
Automatically register containers with [nginx-confd-proxy](https://github.com/chatspry/nginx-confd/tree/master/proxy)

## Basic usage
Running the registration application looks something like this.

```bash
$ docker run \
    -e HOST=$(hostname) \
    -e ETCD_URL=$(hostname):4001 \
    -e CLUSTER=local \
    -e MAPPING=website:*:8080 \
  chatspry/nginx-confd-reg
```

## Map containers to upstreams
The `MAPPING` environment variable can be used to register containers matching patterns.
When container name and internal port combinations are matched, the containers host port along with the host ip will be registered with the upstream.

- `<upstream>`:`<glob pattern>`:`<internal port>`
- `web`:`web_*`:`3000`

You can easily map multiple types of containers differently by giving the `MAPPING` variable comma separated mapping values.

```bash
MAPPING=x_site:x_site_*:3000,y_site:y_site_*:8080
```