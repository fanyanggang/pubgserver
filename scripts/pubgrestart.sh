#!/bin/sh
pubgPath="/data/app/pubgserver"
while true; do
        server=`ps aux | grep pubg.server | grep -v grep`
        if [ ! "$server" ]; then
            cd $pubgPath ; nohup $pubgPath/pubg.server > run.log 2>&1 &
            sleep 10
        fi
        sleep 5
done