#!/bin/bash

set -e

export REGISTRY_SERVER="https://index.docker.io/v1/"
export REGISTRY_USER="matthewberthoud"
export REGISTRY_PASSWORD='ag.6S?5wLLY9dSJ'

kubectl create secret docker-registry push-secret \
    --docker-server=$REGISTRY_SERVER \
    --docker-username=$REGISTRY_USER \
    --docker-password=$REGISTRY_PASSWORD
