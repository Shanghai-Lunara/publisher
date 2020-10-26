.PHONY: mod proto dc

mod:
	go mod download
	go mod tidy

proto:
	cd scripts && bash ./gen.sh api

dc:
	go mod vendor
	cd scripts && bash ./gen.sh deepcopy
	rm -rf vendor