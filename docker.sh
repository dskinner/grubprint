#! /usr/bin/env bash

SSH_PORT=2222
POSTGRES_PORT=5555
DEV_ARGS="-p $SSH_PORT:22 -p $POSTGRES_PORT:5432"
DATA=`pwd`/data

usage() {
	echo "Usage of $0:"
	echo "  build"
	echo "  start-dev"
	echo "  start-prod"
	echo "  clean"
	echo "  shell"
}

start() {
	if [ "`docker ps -a | grep $1`" ]; then
		docker start $1
	else
		docker run -v `pwd`/data:/data/ -d -name $1 $2 dasa/food
	fi
}

clean() {
	if [ "`docker ps | grep $1`" ]; then
		echo "docker stop $1"
		docker stop $1
	fi

	if [ "`docker ps -a | grep $1`" ]; then
		echo "docker rm $1"
		docker rm $1
	fi
}

dropData() {
	if [ -d "$DATA/postgres" ]; then
		read -p "Remove \"$DATA/postgres\" directory? " -n 1 -r
		echo
		if [[ $REPLY =~ ^[Yy]$ ]]; then
			echo "sudo rm -R \"$DATA/postgres\""
		    sudo rm -R "$DATA/postgres"
		fi
	fi
}

case "$1" in
"build")
	docker build -t dasa/food .
	;;
"start-dev")
	start "dasa_food_dev" "$DEV_ARGS"
	;;
"start-prod")
	start "dasa_food_prod"
	;;
"clean")
	clean "dasa_food_prod"
	clean "dasa_food_dev"
	dropData
	;;
"shell")
	ssh -p $SSH_PORT -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no root@localhost
	;;
*)
	usage
	exit 1
	;;
esac
