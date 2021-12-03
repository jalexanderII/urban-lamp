.PHONY: protos

protos:
	protoc -I=./protos --go_opt=paths=source_relative --go_out=plugins=grpc:./server/sensorpb ./protos/sensor.proto