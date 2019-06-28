add-submodule:
	git submodule add git@github.com:mf-sakura/bh_proto.git proto

proto-gen-go:
	protoc --proto_path=proto/gateway/v1/ --go_out=plugins=grpc:app/proto/gateway proto/gateway/v1/reserve.proto
	protoc --proto_path=proto/hotel/v1/ --go_out=plugins=grpc:app/proto/hotel proto/hotel/v1/hotel.proto
	protoc --proto_path=proto/user/v1/ --go_out=plugins=grpc:app/proto/user proto/user/v1/user.proto

update-submodule:
	git submodule foreach 'git fetch;git checkout master; git pull'

build:
	go build -o bh_gateway github.com/mf-sakura/bh_gateway/app/cmd