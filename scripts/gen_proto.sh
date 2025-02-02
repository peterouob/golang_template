#!/bin/zsh

PROTO_DIR="../api/protobuf"

if [[ -d "$PROTO_DIR" ]]; then
  cd "$PROTO_DIR" || { echo "Failed to navigate to $PROTO_DIR"; exit 1; }
else
  echo "Directory $PROTO_DIR does not exist."
  exit 1
fi

for file in *.proto;do
  if [[ -f "$file" ]]; then
    echo "Process generator proto ..."
    protoc \
        -I . \
        -I /path/to/googleapis \
        --go_out . \
        --go-grpc_out . \
        --grpc-gateway_out . \
        "$file"
    sleep 1
    go mod tidy
  else
    echo "No found any file"
    exit 0
  fi
done

echo "generator proto execute success ..."

# protoc -I . --grpc-gateway_out .   --grpc-gateway_opt paths=source_relative     --grpc-gateway_opt generate_unbound_methods=true  --go_out="." --go-grpc_out="."  user.proto

# protoc \
  #  -I . \
  #  -I /path/to/googleapis \
  #  --go_out . \
  #  --go-grpc_out . \
  #  --grpc-gateway_out . \
  #  yourfile.proto