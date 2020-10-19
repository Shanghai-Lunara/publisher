.PHONY: mod

mod:
	go mod download
	go mod tidy

proto:
	cd scripts && bash ./gen.sh api