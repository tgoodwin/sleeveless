# Simulation

## KWOK
todo

```
brew install kwok
```

## cluster setup
1. create a kind cluster with kwok configured
```
kwokctl create cluster --runtime kind
```

2. create a virtual node so we have somewhere to schedule virtual pods
```
kubectl apply -f nodes.yaml
```

3. optionally, create some virtual pods
```
kubectl apply -f fakepods.yaml
```

Or, when starting with a pre-existing kind cluster:
```
sh deploy-kwok.sh
```

At this point, the `kwok-controller` is running in a cluster pod. It creates _virtual nodes_, but since we're running all this in a kind cluster, Kind will put a kindnet pod on this node as it does for all other nodes. The problem is that `kwok` by default assigns the K8s IP (different from the host IP range) of the kwok-controller pod to all virtual nodes that the controller creates. We want to configure the virual nodes to have IPs in the same CIDR range as other IPs in the Kind cluster so that Kind's networking won't crash in the presence of virutal nodes. To do so, we want to configure the kwok-controller's pod to run with `hostNetwork: true`.
```
kubectl patch deployment kwok-controller -n kube-system --patch '{"spec":{"template":{"spec":{"hostNetwork":true}}}}'
```
After applying this patch, we can go ahead and create some virtual nodes.

```
kubectl apply -f nodes.yaml
```
