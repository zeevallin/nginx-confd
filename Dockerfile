FROM debian:wheezy

MAINTAINER Philip Vieira <philip@chatspry.com>

# add nginx wheezy apt server
RUN apt-key adv --keyserver pgp.mit.edu --recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62 \
  && echo "deb http://nginx.org/packages/mainline/debian/ wheezy nginx" >> /etc/apt/sources.list \
  && echo "deb http://http.debian.net/debian wheezy-backports main" >> /etc/apt/sources.list

# ser versions for nginx and confd
ENV JQ_VERSION 1.4-1~bpo70+1
ENV NGINX_VERSION 1.7.10-1~wheezy
ENV CONFD_VERSION 0.7.1

# upgrade operating system and install dependencies with apt
ENV DEBIAN_FRONTEND noninteractive
RUN apt-get update -y \
  && apt-get upgrade -y \
  && apt-get install -y --force-yes \
    nginx=${NGINX_VERSION} \
    jq=${JQ_VERSION} \
    curl \
  && apt-get clean -y \
  && rm -f \
    /etc/nginx/sites-enabled/default \
    /etc/nginx/nginx.conf
ENV DEBIAN_FRONTEND newt

# install confd from the github repository
RUN curl -s -L https://github.com/kelseyhightower/confd/releases/download/v${CONFD_VERSION}/confd-${CONFD_VERSION}-linux-amd64 > confd \
  && mv confd /usr/local/bin/confd \
  && chmod +x /usr/local/bin/confd

COPY html /etc/nginx/html
COPY templates /etc/confd/templates
COPY conf.d /etc/confd/conf.d

# set up entrypoint for docker
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

# set up default scripts
COPY run.sh /run.sh
COPY certs.sh /certs.sh
RUN chmod +x /run.sh /certs.sh

WORKDIR /
USER root
CMD ["/run.sh"]
