FROM debian:wheezy

MAINTAINER Philip Vieira <philip@chatspry.com>

RUN \
  apt-key adv --keyserver pgp.mit.edu --recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62 && \
  echo "deb http://nginx.org/packages/mainline/debian/ wheezy nginx" >> /etc/apt/sources.list

ENV NGINX_VERSION 1.7.6-1~wheezy

RUN apt-get update -y && \
  apt-get install -y --force-yes nginx=${NGINX_VERSION} wget curl && \
  apt-get clean -y

RUN \
  ln -sf /dev/stdout /var/log/nginx/access.log && \
  ln -sf /dev/stderr /var/log/nginx/error.log && \
  rm -f /etc/nginx/sites-enabled/default

RUN wget -O confd https://github.com/kelseyhightower/confd/releases/download/v0.6.3/confd-0.6.3-linux-amd64
RUN \
  mv confd /usr/local/bin/confd && \
  chmod +x /usr/local/bin/confd

RUN mkdir -p /etc/confd/{conf.d,templates}

ADD ./nginx.conf.tmpl /etc/confd/templates/nginx.conf.tmpl
ADD ./nginx.toml /etc/confd/conf.d/nginx.toml
ADD ./run.sh /opt/run.sh

RUN chmod +x /opt/run.sh

EXPOSE 80 443

USER root

CMD /opt/run.sh