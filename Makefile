DB_Name ?= my_database
DB_User ?= postgres
DB_Password ?= password

pg-local:
    @docker run -d --name=co-postgres -p 5471:5432 \
    -e POSTGRES_USER=$(DB_User) \
    -e POSTGRES_PASSWORD=$(DB_Password) \
    -e POSTGRES_DB=$(DB_Name) \
    postgres:17.3-bookworm