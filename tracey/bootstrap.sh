#!/bin/bash

set -euo

# install cert manager
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.14.3/cert-manager.yaml

# wait for the cert manager pods to be ready before proceeding
sleep 10

# install tempo operator controller
kubectl apply -f https://github.com/grafana/tempo-operator/releases/latest/download/tempo-operator.yaml

# wait for the tempo operator pods to be ready before proceeding
sleep 10

# set up object storage
kubectl apply -f https://raw.githubusercontent.com/grafana/tempo-operator/41d57e9ec1f78bc9789d3cf55241b2fed2faa269/minio.yaml

# configure access to object storage
kubectl apply -f tempo-storage-secret.yaml

# wait for the minio pods to be ready before proceeding
sleep 10

# install tempo stack custom resource
kubectl apply -f tempo.yaml

# install grafana
sudo helm repo add grafana https://grafana.github.io/helm-charts
sudo helm repo update
sudo helm upgrade --install grafana grafana/grafana -f grafana.yaml
