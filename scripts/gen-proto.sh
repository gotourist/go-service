#!/bin/bash
CURRENT_DIR=$(pwd)

# shellcheck disable=SC2086
# shellcheck disable=SC2044
for module in $(find $CURRENT_DIR/protos/* -type d); do
    echo $module
    protoc -I /usr/local/include \
           -I $GOPATH/src/github.com/gogo/protobuf/gogoproto \
           -I $CURRENT_DIR/protos/ \
            --gofast_out=plugins=grpc:$CURRENT_DIR/genproto/ \
            $module/*.proto;
done;

# shellcheck disable=SC2044
for module in $(find $CURRENT_DIR/genproto/* -type d); do
  if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i "" -e "s/,omitempty//g" $module/*.go
  else
    sed -i -e "s/,omitempty//g" $module/*.go
  fi
done;