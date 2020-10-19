#!/bin/bash

CWD=$(cd $(dirname $0); pwd)
BINARY=$CWD/bin/awesome-bot
CONF_DIR=$CWD/conf
LARK_CONF=$CONF_DIR/lark_config.yaml
CHROME_PATH=/Users/bytedance/code/awesome-bot/chromedriver

export TZ=Asia/Shanghai
echo "$BINARY -lark-config=$LARK_CONF -chrome-driver-path=$CHROME_PATH -chrome-driver-port=9090"
exec $BINARY -lark-config=$LARK_CONF -chrome-driver-path=$CHROME_PATH -chrome-driver-port=9090
