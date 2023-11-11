# 🚗 ReCo web backend 🌐 ☁️

## 📖 Prerequisites

###  GOLANG

[Download](https://go.dev/learn/)

### GO AIR Package

Used for live reload while developing.

```bash
go install github.com/cosmtrek/air@latest
```

### Delve

Used for debugging.

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

## First time setup

### 📦 Dependencies

To install all dependencies run the command below.

```bash
go build .
```

## 📅 Daily business

Use the command below to start the project with live reload.

```bash
air
```

###  Development without live reload

If you want to develop without live reload you can use the command below. However, this is not recommended.

```bash
go run .
```

### Development outside docker

Warning: This method is not recommended.

The Database need to be installed. We are currently using `timescaledb-ha:pg14-latest` as database. The database can be installed with the command below.

```bash
docker pull timescale/timescaledb-ha:pg14-latest
```

The database can be started with the command below.

```bash
docker run -d --name timescaledb -p 5555:5432 -e POSTGRES_PASSWORD=postgres timescale/timescaledb-ha:pg14-latest
```


## 🐞 Debugging

Debugging is done with Delve using VSCode

- Open the project in VSCode
- Start the project with `air`
- Open the `Run and Debug` tab in VSCode
- Select `Attach to Air` from the dropdown
- You are ready to debug 🐞


## Database migrations

We are using [golang-migrate/migrate](https://github.com/golang-migrate/migrate) for database migrations. Currently the migration tools is integrated to the project. Migrations are run automatically when the project is started. Migration are written in SQL and can be found in the `migrations` folder. Naming convention for migrations is `<version>_<name>.up.sql` and `<version>_<name>.down.sql`. The version number is used to determine the order of the migrations. The name is used to describe the migration. The `up` file is used to migrate the database forward and the `down` file is used to rollback the migration.

Migrations can be run manually with the command below, but the cli tools needs to be installed first. Read more about CLI migrations in the [documentation](https://github.com/golang-migrate/migrate).



