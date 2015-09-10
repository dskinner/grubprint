# postgres

database server component for grubprint.io

```bash
##
# initialize postgres
##

cd $GOPATH/src/grubprint.io/postgres
docker build -t grubprint/db .
docker run -d -p 5432:5432 --name grubprint_db grubprint/db

# the server is now running and initializing database.
# this may take a few minutes; check status
docker logs grubprint_db

##
# manage container
##

# start and stop container
docker stop grubprint_db
docker start grubprint_db

# open shell in container
docker exec -it grubprint_db bash

# remove container and image
docker rm grubprint_db
docker rmi grubprint/db
```
