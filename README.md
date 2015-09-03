## Choose new name

a core goal of the project is meal planning and cost analysis

here's some words that could be combined:

* agent
* diet - plays well with marketing
* plan
* meal
* food
* geek - ... i mean ... cost analysis ... so ...

## Postgres

```bash
$ cd postgres
$ docker build -t food/db .
$ docker run -d -p 5432:5432 --name food_db food/db
```

May take a few minutes to initialize database. Check status with `docker logs food_db`.

## Application

```bash
$ go get github.com/cespare/reflex
$ reflex -c reflex.conf
```

To test, visit http://localhost:8080/foods/cheese%20blue
