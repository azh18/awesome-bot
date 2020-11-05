#!/bin/bash
CWD=$(cd $(dirname $0); pwd)

bash /opt/bin/entry_point.sh > $CWD/driver.log 2>&1 &
sleep 5
bash $CWD/bootstrap_no_driver.sh
