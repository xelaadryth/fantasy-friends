VERSION = `cat VERSION`

.PHONY: build run down

build:
	@export VERSION=$(VERSION) && docker-compose build

run:
	@export VERSION=$(VERSION) && docker-compose up

down:
	@export VERSION=$(VERSION) && docker-compose down

deploy:
	@heroku container:push web
