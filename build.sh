#! /bin/sh

PROTOC=${PROTOC:-/opt/protoc}

# protoc-gen-go for go 1.22
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.0
# protoc-gen-go-grpc for go 1.22
#go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
export PATH=$PATH:$(go env GOPATH)/bin

PROTO_FILE="bfe_access_pb/bfe_access.proto"

echo "start building ......"

${PROTOC} --go_out=./ --go_opt=paths=source_relative ${PROTO_FILE}
if [ $? != 0 ] 
then 
    exit 1
fi

echo "build succeed!"
