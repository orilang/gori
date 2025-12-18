#!/bin/bash

usage(){
  echo "$1 argument is mantatory"
  echo "Example: $0 1.21.6 1.22.0"
}

if [ -n "$1" ]
then
  export OLD_VERSION="$1"
else
  usage OLD_VERSION
  exit 1
fi
shift

if [ -n "$1" ]
then
  export NEW_VERSION="$1"
else
  usage NEW_VERSION
  exit 1
fi
shift

OS=$(uname -s | tr -s "A-Z" "a-z")
if [[ "${OS}" == "darwin" ]]
then
  find . -type f -name go.mod -exec sed -i '' "s#^go ${OLD_VERSION}#go ${NEW_VERSION}#g" {} \;
  sed -i '' "s#${OLD_VERSION}#${NEW_VERSION}#g" .github/workflows/*.yml
else
  find . -type f -name go.mod -exec sed -i "s#^go ${OLD_VERSION}#go ${NEW_VERSION}#g" {} \;
  sed -i "s#${OLD_VERSION}#${NEW_VERSION}#g" .github/workflows/*.yml
fi

go mod tidy