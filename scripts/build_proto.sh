BASE_PATH=pkg

protoc \
    --proto_path=$BASE_PATH \
    --go_out=$BASE_PATH \
    --go_opt=paths=source_relative \
    --go-grpc_out=$BASE_PATH \
    --go-grpc_opt=paths=source_relative \
    $BASE_PATH/tunnels/protobuf/*.proto

protoc \
    --proto_path=$BASE_PATH \
    --go_out=$BASE_PATH \
    --go_opt=paths=source_relative \
    --go-grpc_out=$BASE_PATH \
    --go-grpc_opt=paths=source_relative \
    $BASE_PATH/pipelines/protobuf/*.proto