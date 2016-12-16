#!/bin/bash

defaultLogFile="log_{{.}}.log"

#保存 log 文件
if test -e $defaultLogFile
then
    date1=$(date "+%Y_%m_%d_%H_%M_%S")
    logFile=${defaultLogFile}"."${date1}
    mv ${defaultLogFile} $logFile
else
    echo "log file not found"
fi

 nohup ./{{.}}_linux64 --config=conf/config.json  2>&1 | bash ./log_split.sh &
