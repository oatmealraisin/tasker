all: proto release

proto:
	protoc -I=proto/ --go_out=pkg/models proto/task.proto

release:
	go build tasker.go


.PHONY: proto release all
