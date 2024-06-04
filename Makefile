run: build
	@./bin/dreampicai

install:
	@go install github.com/a-h/templ/cmd/templ@latest
	@go get ./...
	@go mod vendor
	@go mod tidy
	@go mod download
	@npm install -D tailwindcss
	@npm install -D daisyui@latest
	@npm install -D drizzle-kit

css:
	@npx tailwindcss -i view/css/app.css -o public/styles.css --minify --watch 

templ:
	@templ generate --watch --proxy=http://localhost:3131 --open-browser=false --proxybind="localhost" --proxyport="3000"

air:
	@air

live: 
	make -j2 css templ

build:
	@npx tailwindcss -i view/css/app.css -o public/styles.css --minify
	@templ generate view
	@go build -tags dev -o bin/dreampicai main.go 

# db migration stuff
migrate: ## Database migration up
	# @go run cmd/migrate/main.go up
	@npx drizzle-kit migrate

reset:
	# @go run cmd/reset/main.go up

down: ## Database migration down
	# @go run cmd/migrate/main.go down
	@npx drizzle-kit drop

generate: ## Migrations against the database
	#@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))
	@npx drizzle-kit generate

miup:
	@npx drizzle-kit generate
	@npx drizzle-kit migrate

seed:
	# @go run cmd/seed/main.go

studio:
	@npx drizzle-kit studio
