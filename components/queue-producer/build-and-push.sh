#!/bin/bash
set -eu

tag=${1:?Specify tag as first argument}

image=europe-west1-docker.pkg.dev/fluent-buckeye-343615/workshop/producer:"$tag"

docker build -t "$image" .
docker push "$image"
