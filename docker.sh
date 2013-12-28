#! /usr/bin/env bash

SSH_PORT=2222
POSTGRES_PORT=5555
APP_PORT=9000
DEV_ARGS="-p $APP_PORT:$APP_PORT -p $SSH_PORT:22 -p $POSTGRES_PORT:5432"
PROD_ARGS="-p $APP_PORT:$APP_PORT"
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
		mkdir $DATA/app
		CWD=`pwd`
		cd $DATA/app
		revel package dasa.cc/food
		tar -xzf ./food.tar.gz
		cd "$CWD"
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
	if [ -d "$DATA/app" ]; then
		read -p "Remove \"$DATA/app\" directory? " -n 1 -r
		echo
		if [[ $REPLY =~ ^[Yy]$ ]]; then
			echo "sudo rm -R \"$DATA/app\""
			sudo rm -R "$DATA/app"
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
	start "dasa_food_prod" "$PROD_ARGS"
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
