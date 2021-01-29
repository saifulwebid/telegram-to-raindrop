.PHONY: run
run:
	go run cmd/main.go

.PHONY: deploy
deploy:
	./bin/deploy