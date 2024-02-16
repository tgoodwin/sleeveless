#!/usr/bin/env sh

# tempo distributor service needs static ip
# and clusterIPs are immutable

# ./longhorn-values.yaml && \
./grafana-values.yaml && \
./tempo-values.yaml && \
kubectl delete svc tempo-distributor-service.yaml && \
./tempo-distributor-service.yaml
