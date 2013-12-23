#! /usr/bin/env sh
docker run -p 2222:22 -p 5555:5432 -v `pwd`:/root/ -d -name dasa_food dasa/food
