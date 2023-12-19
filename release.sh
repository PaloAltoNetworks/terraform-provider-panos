#!/bin/bash

set -x 
version="1.11.6"

for os in darwin linux; 
do
    for arch in amd64;
    do
        builddir=terraform-provider-panos_${version}_${os}_${arch}
        mkdir -p $builddir
        cp LICENSE $builddir/
        cp CHANGELOG.md $builddir/
        cp README.md $builddir/
        GOOS=$os GOARCH=amd64 go build -o $builddir/terraform-provider-panos_$version
        cd $builddir
        zip -vr -X ../$builddir.zip . -x "*.DS_Store"
        cd ..
    done
done

shasum -a 256 *.zip > terraform-provider-panos_${version}_SHA256SUMS
gpg --local-user fpluchorg --detach-sign terraform-provider-panos_${version}_SHA256SUMS
