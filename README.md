# nginx-confd
Service discovery for NGINX configured through confd

## Test locally
You can test on your local machine by tweaking `docker-compose.yml` to match config for your local host.
By default you need to have a host record `docker` pointing at your docker host machine for this to work.
Then you can scale the `app` service to whatever number you like.

```
$ docker-compose up -d
$ docker-compose scale app=5
```