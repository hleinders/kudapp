#!/bin/bash

rm -f kudapp
for t in latest red blue green
do
    docker image rm kudapp:$t
done

if [[ -n "${CONTAINER_REGISTRY}" ]]
then
    PREFIX="${CONTAINER_REGISTRY}/"

    for t in latest red blue green
    do
        docker image rm ${PREFIX}kudapp:$t
    done
fi
