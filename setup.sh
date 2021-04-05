#!/usr/bin/env bash

# kubectl delete  -f deployment/deployment.yaml

docker build -f Dockerfile.build . -t local/build-kdaudit
docker build -f Dockerfile.kdaudit . -t local/kdaudit

# kubectl apply -f deployment/
