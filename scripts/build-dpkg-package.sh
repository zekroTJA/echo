#!/usr/bin/env bash

set -e

if [[ -z $1 ]]; then
    echo "usage: $(basename "$0") {arch} [output-file-name]"
    exit 1
fi

arch=$1
tag=$(git describe --tag --abbrev=0)

sed 's/<<VERSION>>/'"${tag:1}"'/' -i dist/dpkg/DEBIAN/control
sed 's/<<ARCH>>/'"${arch}"'/' -i dist/dpkg/DEBIAN/control

dpkg-deb -Zxz --build --root-owner-group \
    dist/dpkg "${2:-"echo-${arch}.deb"}"
