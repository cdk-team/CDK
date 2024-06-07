#!/bin/bash

current_file=$0
pdir=`dirname "${current_file}"`
ppdir=`dirname "${pdir}"`
project_dir=`dirname "${ppdir}"`

# CGO_CFLAGS="-Wno-undef-prefix=TARGET_OS_ -Wno-deprecated-declarations"
#CGO_CFLAGS="-Wno-deprecated-declarations"
#export CGO_CFLAGS="-Wno-undef -Wno-deprecated-declarations"

go run "$project_dir/cmd/cdk/cdk.go" "$@" 2>&1 | sed '/# github.com/,/has been explicitly marked deprecated here$/d'
