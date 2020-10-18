#!/bin/bash

CWD=$(cd $(dirname $0); pwd)
BINARY=$CWD/bin/awesome-bot
CONF_DIR=$CWD/conf
LARK_CONF=$CONF_DIR/lark_config.yaml

export TZ=Asia/Shanghai
echo "$BINARY -lark-config=$LARK_CONF"
exec $BINARY -lark-config=$LARK_CONF
