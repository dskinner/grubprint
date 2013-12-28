FROM ubuntu:precise

MAINTAINER Daniel Skinner <daniel@dasa.cc>

ENV DEBIAN_FRONTEND noninteractive

RUN dpkg-divert --local --rename --add /sbin/initctl
RUN ln -s /bin/true /sbin/initctl

RUN echo "deb http://archive.ubuntu.com/ubuntu precise main universe" > /etc/apt/sources.list
RUN apt-get -y update
RUN apt-get -y upgrade
RUN apt-get -y install vim wget openssh-server

RUN mkdir -p /var/run/sshd
RUN echo "root:root" | chpasswd

RUN locale-gen en_US.UTF-8

RUN echo "deb http://apt.postgresql.org/pub/repos/apt/ precise-pgdg main" > /etc/apt/sources.list.d/pgdg.list
RUN wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | apt-key add -
RUN apt-get -y update
RUN apt-get -y install postgresql-9.3 postgresql-contrib-9.3

CMD sh /data/run.sh
