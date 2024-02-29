build:
	docker build -t forumv2 .
image-run:
	docker run --name forum -d -p 8080:8080 --rm forumv2
run:
	go run ./cmd/forum/
stop:
	docker stop forum