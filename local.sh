#!/bin/bash

cd `dirname $0`
CURRENT=`pwd`
echo $CURRENT

function test
{
   go test -v $(go list ./... | grep -v vendor) --count 1
}

CMD=$1
shift
$CMD $*