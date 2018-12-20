#!/bin/bash

cd `dirname $0`
CURRENT=`pwd`
echo $CURRENT

function test
{
   if [ -e $CURRENT/local_env.sh ]; then
     source $CURRENT/local_env.sh
   fi
   go test -v $(go list ./... | grep -v vendor) --count 1
}

CMD=$1
shift
$CMD $*