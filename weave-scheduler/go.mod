module weave-scheduler

go 1.23

toolchain go1.24.5

require (
	github.com/robfig/cron/v3 v3.0.1
	weave-module v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.3.1 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/redis/go-redis/v9 v9.2.1 // indirect
	github.com/streadway/amqp v1.1.0 // indirect
	golang.org/x/crypto v0.14.0 // indirect
	golang.org/x/text v0.13.0 // indirect
	gorm.io/driver/postgres v1.5.3 // indirect
	gorm.io/gorm v1.25.5 // indirect
)

replace weave-module => ../weave-module
