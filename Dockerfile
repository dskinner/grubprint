FROM debian:wheezy

MAINTAINER Daniel Skinner <daniel@dasa.cc>

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get -y update
RUN apt-get -y upgrade

ADD ./food /root/
ADD ./public /root/public
ADD ./templates /root/templates

WORKDIR /root/
CMD ./food
