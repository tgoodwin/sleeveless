## ZooKeeper bug 312
https://github.com/pravega/zookeeper-operator/issues/312

This controller is compatible with Kubernetes version 1.19.11.

## Exercising the bug
To test, manually run
1.
```
kubectl apply -f zkc-1.yaml
```
2.
```
kubectl delete -f zkc-1.yaml
```
3.
```
kubectl apply -f zkc-yaml
```
This does not guarantee bug reproduction, only creates the environment in which the bug may occur.

The bug occurs when after step (3), a stale message created in step (1), which is the message updating the ZKC pod with a non-nil deletionTimestamp, is delivered to the controller. The controller is unable to tell that the message pertains to a no longer extant ZKC pod and interprets it as an update to the existing pod created in step (3). As such, it mistakenly marks the PVC for the ZKC pod as marked for deletion, even though there has been no command to delete the ZKC pod from step (3).
