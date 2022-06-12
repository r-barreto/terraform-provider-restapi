#!/bin/bash


cd fakeservercli || exit 1

case $1 in
  start)

    go build
    nohup ./fakeservercli &
    echo $! > fakeserver_pid.txt
    ;;
  stop)
    kill -9 "$(cat fakeserver_pid.txt)"
    rm fakeserver_pid.txt fakeservercli nohup.out
    ;;
esac