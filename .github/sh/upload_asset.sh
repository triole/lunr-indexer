#!/bin/bash

repo="triole/lunr-indexer"
file="${1}"

curl \
    -H "Authorization: token $GITHUB_TOKEN" \
    -H "Content-Type: $(file -b --mime-type ${file})" \
    --data-binary @${file} \
    "https://uploads.github.com/repos/hubot/singularity/releases/123/assets?name=$(basename ${file})"

# "https://uploads.github.com/repos/${repo}/releases/123/assets?name=$(basename ${file})"
