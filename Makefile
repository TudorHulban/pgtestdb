DB_User ?= postgres
DB_Password ?= password

pg-local:
	@docker run -d --name=co-postgres -p 5471:5432 \
	-e POSTGRES_USER=$(DB_User) \
	-e POSTGRES_PASSWORD=$(DB_Password) \
	postgres:17.2-bookworm