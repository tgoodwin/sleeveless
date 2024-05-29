#!/bin/bash
set -x

NAMESPACE=tracey

kubectl create --namespace $NAMESPACE -f crds
kubectl create --namespace $NAMESPACE -f default_ns/rbac.yaml
kubectl create --namespace $NAMESPACE -f default_ns/operator.yaml
