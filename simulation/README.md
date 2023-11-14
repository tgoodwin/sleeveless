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

