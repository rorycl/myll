# Docker setup for Postgres

21 February 2025

With reference to the Docker postgres official image docs
[here](https://www.docker.com/blog/how-to-use-the-postgres-docker-official-image/).

Pull the latest v17 version of postgres for debian bullseye.
```
docker pull postgres:17-bullseye
```

```
# example
docker run --name mll -d postgres
```

export PG_ADMIN_PW="adminpw"

```
# access container
docker exec -it model_db_1 bash
```
