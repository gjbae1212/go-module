#!/bin/bash
set -e -o pipefail
trap '[ "$?" -eq 0 ] || echo "Error Line:<$LINENO> Error Function:<${FUNCNAME}>"' EXIT

cd `dirname $0`
CURRENT=`pwd`

function test
{
   set_env
   go test -v $(go list ./... | grep -v vendor) --count 1 -race -coverprofile=$CURRENT/coverage.txt -covermode=atomic
}

function test_part
{
   set_env
   local part=$1
   if [ -z $part ]
   then
      exit 1
   fi

   local modules=(`go list ./... | tr ' ', '\n'`)
   for item in "${modules[@]}"
   do
      if [[ $item =~ .*${part}.* ]]
      then
        go test -v $item --count 1
      fi
   done
}

function codecov
{
   /bin/bash <(curl -s https://codecov.io/bash)
}

function set_env
{
   if [ -e $CURRENT/local_env.sh ]; then
     source $CURRENT/local_env.sh
   fi
}
CMD=$1
shift
$CMD $*
