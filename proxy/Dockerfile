FROM debian:wheezy

MAINTAINER ["Philip Vieira", "<philip@chatspry.com>"]

# add nginx wheezy apt server
RUN apt-key adv --keyserver pgp.mit.edu --recv-keys 573BFD6B3D8FBC641079A6ABABF5BD827BD9BF62 \
  && echo "deb http://nginx.org/packages/mainline/debian/ wheezy nginx" >> /etc/apt/sources.list \
  && echo "deb http://http.debian.net/debian wheezy-backports main" >> /etc/apt/sources.list

# ser versions for nginx and confd
ENV JQ_VERSION 1.4-1~bpo70+1
ENV NGINX_VERSION 1.9.1-1~wheezy
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
    /etc/nginx/html \
    /etc/nginx/nginx.conf
ENV DEBIAN_FRONTEND newt

# install confd from the github repository
RUN curl -s -L https://github.com/kelseyhightower/confd/releases/download/v${CONFD_VERSION}/confd-${CONFD_VERSION}-linux-amd64 > confd \
  && mv confd /usr/local/bin/confd \
  && chmod +x /usr/local/bin/confd

COPY static /etc/nginx/static/default
COPY include /etc/nginx/include/default
COPY nginx.conf /etc/confd/templates/
COPY nginx.toml /etc/confd/conf.d/

# set up entrypoint for docker
COPY docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT ["/docker-entrypoint.sh"]

# set up default scripts
COPY run.sh /run.sh
RUN chmod +x /run.sh

# add custom static files and includes on build
ONBUILD COPY static /etc/nginx/static/custom
ONBUILD COPY include /etc/nginx/include/custom

WORKDIR /
USER root
CMD ["/run.sh"]
