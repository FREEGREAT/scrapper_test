migrate:
	goose -dir ./sub_worker_cron/migration postgres "postgres://user:user@localhost:5430/subscribe_db" up

migrate-down:
	goose -dir ./sub_worker_cron/migration postgres "postgres://user:user@localhost:5430/subscribe_db" down