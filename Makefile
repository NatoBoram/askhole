install:
	go install .

build:
	go build .

run:
	go run .

test:
	go test .

clean:
	rm askhole

# Docker

build-docker:
	docker build --tag askhole .

run-docker:
	docker run --publish 127.0.0.1:9123:9123 --name askhole askhole

build-run-docker: build-docker run-docker

kill-docker:
	docker ps --format '{{.Image}} {{.ID}}' | grep askhole | awk '{print $2}' | xargs docker kill
