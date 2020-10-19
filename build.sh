#!/bin/bash

CWD=$(cd $(dirname $0); pwd)
OUTPUT_DIR=$CWD/output

mkdir -p $OUTPUT_DIR/bin

go build -o $OUTPUT_DIR/bin/awesome-bot $CWD/cmd/v1/main.go
cp -rf $CWD/conf $OUTPUT_DIR/conf
cp $CWD/scripts/bootstrap_need_driver.sh $OUTPUT_DIR/bootstrap_need_driver.sh
cp $CWD/scripts/bootstrap_no_driver.sh $OUTPUT_DIR/bootstrap_no_driver.sh
cp $CWD/scripts/entry_point.sh $OUTPUT_DIR/entry_point.sh
