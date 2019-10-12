#!/bin/sh

action="go test -v ./..."

chsum1=""

while [[ true ]]
do
    chsum2=`find . -type f -exec md5sum {} \;`
    if [[ $chsum1 != $chsum2 ]] ; then
        eval "$action"
        chsum1=$chsum2
    fi
    sleep 1
done
