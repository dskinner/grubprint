## Postgres

```bash
$ cd postgres
$ docker build -t food/db .
$ docker run -d -p 5432:5432 --name food_db food/db
```

May take a few minutes to initialize database. Check status with `docker logs food_db`.

## Application

```bash
$ cd cmd/food
$ go run main.go
```

Only api works. To test, visit http://localhost:8080/foods/cheese%20blue
