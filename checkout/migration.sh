goose -dir ./migrations postgres "postgres://user:password@localhost:5432/checkout?sslmode=disable" status

goose -dir ./migrations postgres "postgres://user:password@localhost:5432/checkout?sslmode=disable" up