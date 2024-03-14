pb:
	./scripts/update_proto.sh

docker:
	./scripts/build_all.sh

test:
	go test -v ./...
