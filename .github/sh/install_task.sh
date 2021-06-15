#!/bin/bash

bin="task"
url="https://github.com/go-task/task/releases"
grep_scheme='(?<=href\=\").*linux_amd64.*\.tar.gz(?=")'
# additional_args="--strip-components 1"

tfol="/home/runner/.local/bin"

function fail() {
    x=$?
    echo "last command failed, exitcode ${x}"
    if ((${x} == 3)); then
        echo "malformed url ${bin_url}"
    fi
    echo "abort"
    exit
}

function rcmd() {
    echo -e "\n\033[0;93m${1}\033[0m"
    ${1} || fail
}

function install() {
    tmpfil="/tmp/install_tmp.tar.gz"

    echo "Install ${bin}"
    bin_url="https://github.com$(
        curl -Ls "${url}" |
            grep -Po "${grep_scheme}" | head -n 1
    )"

    mkdir -p "${tfol}"
    rcmd "curl -L ${bin_url} -o ${tmpfil}"
    rcmd "tar -xzf ${tmpfil} -C ${tfol} ${additional_args}"

    echo -e "${bin} installed\n"
    ${tfol}/${bin} --version
}

${bin} version >/dev/null 2>&1 || install
