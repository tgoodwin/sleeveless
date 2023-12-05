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

# patch the kwok-controller to use the host network so that pods scheduled onto the virtual node(s) are
# given IP ranges that match the IP ranges of pods that kind runs on cluster nodes. This step allows KWOK to play
# nicely with Kind and the way Kind implements the kubernetes network model.
kubectl patch deployment kwok-controller -n kube-system --patch '{"spec":{"template":{"spec":{"hostNetwork":true}}}}'

# now, let's deploy a KWOK virutal node!
kubectl apply -f nodes.yaml
