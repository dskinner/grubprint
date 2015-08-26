FROM debian:jessie

MAINTAINER Daniel Skinner <daniel@dasa.cc>

ENV DEBIAN_FRONTEND noninteractive

RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -sf /bin/true /sbin/initctl

RUN apt-get -y update
RUN apt-get -y upgrade

ADD ./food /root/
ADD ./public /root/public
ADD ./templates /root/templates

WORKDIR /root/
CMD ./food
