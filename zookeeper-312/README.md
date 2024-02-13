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


## Hacking on the source code
A copy of the zookeeper source code is in `./src/zookeeper-operator`.
Inside this folder is another directory `./custom` which contains copies of `client-go` and `controller-runtime`. Both of these local repos have been checked out to the git tag that corresponds with the version used in `zookeeper-operator/go.mod`.

NOTE: `zookeeper-operator/go.mod` calls for `client-go v0.27.5` and `controller-runtime v0.15.2` whereas `src/zookeeper-operator/custom/controller-runtime/go.mod` (on version v0.15.2) calls for `client-go v0.27.2`. _Seems like our hacked controller binary might need 2 versions of client-go, but in practice, the Go build step will only pull in client-go v0.27.5. So there's some silent version promotion happening behind the scenes. Just something to be aware of._

The `go.mod` file in `./src/zookeeper-operator` has been updated to use these copies of the client libraries.

Further, the `zookeeper-operator/Makefile` has been modified to build the image and tag it with "sleeveless" and push it to `docker.io/tlg2132` and `default_ns/operator.yaml` has been modified to pull the image down from `docker.io/tlg2132/zookeeper-operator:sleeveless`. Modify as needed.