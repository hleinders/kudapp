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


docker build . -t kudapp:latest
docker build . -t kudapp:red --build-arg DEAFULT_COLOR=red
docker build . -t kudapp:blue --build-arg DEAFULT_COLOR=blue
docker build . -t kudapp:green --build-arg DEAFULT_COLOR=green

if [[ -n "${CONTAINER_REGISTRY}" ]]
then
    PREFIX="${CONTAINER_REGISTRY}/"

    for t in latest red blue green
    do
        docker tag kudapp:$t ${PREFIX}kudapp:$t
    done
fi

cat <<EOM


Images build:

    * ${PREFIX}kudapp:latest
    * ${PREFIX}kudapp:red
    * ${PREFIX}kudapp:blue
    * ${PREFIX}kudapp:green

Please don't forget to push them if neccessary
EOM
