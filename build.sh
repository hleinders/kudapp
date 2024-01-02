#!/bin/bash

APP="kudapp"
HTML="html"
TMPL="templates"

errf() {
    echo $1
    exit $2
}


[[ -x ${APP} ]]  || errf "Executable ${APP} not found. Please compile first" 1
[[ -x ${HTML} ]] || errf "Document root ${HTML} not found. Please change location" 1
[[ -x ${TMPL} ]] || errf "Template dir ${TMPL} not found. Please change location" 1


[[ -n "${CONTAINER_REGISTRY}" ]] && PREFIX="${CONTAINER_REGISTRY}/"

docker build . -t ${PREFIX}kudapp:red --build-arg DEAFULT_COLOR=red
docker build . -t ${PREFIX}kudapp:blue --build-arg DEAFULT_COLOR=blue
docker build . -t ${PREFIX}kudapp:green --build-arg DEAFULT_COLOR=green

docker tag ${PREFIX}kudapp:red ${PREFIX}kudapp:latest

cat <<EOM


Images build:

    * ${PREFIX}kudapp:latest
    * ${PREFIX}kudapp:red
    * ${PREFIX}kudapp:blue
    * ${PREFIX}kudapp:green

Please don't forget to push them if neccessary
EOM
