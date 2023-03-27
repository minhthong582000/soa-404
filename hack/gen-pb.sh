#!/bin/bash

cd "$(dirname "${BASH_SOURCE[0]}")"

PB_DIR="../api/v1/pb"
cd $PB_DIR

# Loop through all directories and generate the protobuf files
for dir in $(ls -d */); do
    echo "Generating protobuf files for $dir"
    protoc --go_opt=paths=source_relative -I=$dir --go_out=$dir \
        --go-grpc_opt=paths=source_relative --go-grpc_out=$dir \
        $dir/*.proto
done
