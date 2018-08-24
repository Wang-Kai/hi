rm -rf pb/*.pb.go

protoc -I=. pb/*.proto --go_out=plugins=grpc:.