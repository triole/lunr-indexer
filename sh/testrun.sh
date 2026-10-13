#!/bin/bash
scriptdir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
basedir="${scriptdir%/*}"

cd "${basedir}" &&
  r ../testdata/ -f -o /dev/null $*
