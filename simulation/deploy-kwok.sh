#! /bin/bash
#
# Deploy KWOK (Kubernetes without Kubelet) in a cluster

# KWOK repo
KWOK_REPO=kubernetes-sigs/kwok

# Get latest
KWOK_LATEST_RELEASE=$(curl "https://api.github.com/repos/${KWOK_REPO}/releases/latest" | jq -r '.tag_name')

# deplo kwok and set up CRDs
kubectl apply -f "https://github.com/${KWOK_REPO}/releases/download/${KWOK_LATEST_RELEASE}/kwok.yaml"

# set up default CRs of stages
kubectl apply -f "https://github.com/${KWOK_REPO}/releases/download/${KWOK_LATEST_RELEASE}/stage-fast.yaml"

