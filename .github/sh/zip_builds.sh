#!/bin/bash

scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir=$(echo "${scriptdir}" | grep -Po "^.*(?=/.*/)")
builddir="${basedir}/build"

appname=$(echo ${basedir} | grep -Po "[^/]+$")
tmpdir="/tmp/assets"

function getver() {
    f=$(
        find "${builddir}" -type f -executable | grep "$(arch)" | head -n 1
    )
    eval "${f}" -V | grep -Po "(?<=Version:\s).*"
}

function rcmd() {
    echo -e "\n\033[0;93m${1}\033[0m"
    eval ${1}
}

ver=$(getver)
mkdir -p "${tmpdir}"

for fol in $(find "${builddir}" -maxdepth 1 -mindepth 1 -type d); do
    farch=$(echo "${fol}" | grep -Po "[^/]+$")
    cd "${fol}"
    rcmd "tar -zcvf ${tmpdir}/${appname}_v${ver}_${farch}.tar.gz *"
done
