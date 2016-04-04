#!/bin/sh
set -e
go-bindata -pkg goqmllib -nometadata -o goqmllib/bindata.go templates resources packageconfig
pushd goqmllib >> /dev/null
go fmt
popd >> /dev/null
go fmt
go build -ldflags="-w -s" -o goqml

#rm -rf ./workbench*

#echo "application test"
#mkdir workbench
#cp test/main.go workbench/
#pushd workbench >> /dev/null
#../goqml build
#../goqml pack
#popd >> /dev/null

