#!/bin/bash
echo the listening port is:
read varport
#psname="chatserver"

#echo "PS name is $psname"

#ps=`pgrep $psname`

defaultLogFile="chat_log_${varport}.log"

#保存 log 文件
if test -e $defaultLogFile
then
    date1=$(date "+%Y_%m_%d_%H_%M_%S")
    logFile=${defaultLogFile}"."${date1}
    mv ${defaultLogFile} $logFile
else
    echo "log file not found"
fi

 nohup ./chatserver_linux64 --config=conf/config.json --port=${varport} 2>&1 | bash ./log_split.sh &
