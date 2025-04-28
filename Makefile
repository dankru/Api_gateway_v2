migrate-add: ## Create new migration file, usage: make migrate-add [name=<migration_name>]
	goose -dir database/migrations create $(name) sql