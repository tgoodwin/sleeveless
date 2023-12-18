# TODO

## Building Kubernetes
1. Download Kubernetes
```bash
mkdir -p fakegopath/src/k8s.io
K8S_VERSION="v1.28.0"
git clone --single-branch --branch $K8S_VERSION https://github.com/kubernetes/kubernetes/git fakegopath/src/k8s.io/kubernetes >> /dev/null
```

2. (Optional) install our custom instrumentation dependencies in K8s source
Skip for now!
TODO: when we have custom golang modules, use this step to modify the k8s `go.mod` file(s)
so our code can be incorporated as a dependency when building kubernetes in step 4.

For now, we'll just hack on the K8s source locally, tracking our modifications on a git branch perhaps.

3. (Future step) instrument kubernetes
Skip for now!
Eventually we may have a script to modify fresh kubernetes source code with our changes.
```bash
# go mod tidy
# go build
# ./instrumentation config.json
```

4. Build kubernetes source into a Kind node image
This step takes a minute or two.
```bash
ORIG_DIR=$(pwd)
cd fakegopath/src/k8s.io/kubernetes

GOPATH=$ORIG_DIR/fakegopath KUBE_GIT_VERSION=siren-${K8S_VERSION} kind build node-image
```

After a successful build, you should see output like
```
Image "kindest/node:latest" build completed.
```
Next, we'll tag the image we just built with something we can reference.

```bash
docker image tag kindest/node:latest <container_registry>/node:<image_tag>

# and then push it to a registry where Kind can find it
# and <container_registry> can be something like docker.io/tlg2132 (this is my docker username)
docker push <container_registry>/node:<image_tag>
```

5. Running modified Kubernetes in a Kind cluster locally
Here we tell Kind to find a Kubernetes container image at the location we just pushed to.
```bash
# clean up any kind cluster(s) that may already exist
# this command allows you to not remember what the cluster is called - it'll just delete whatever's there.
kind get clusters | xargs -n 1 kind delete cluster --name 

kind create cluster --image <container_registry>/node:<image_tag>
```


## example workflow
try hacking on this file (the watch cache) by adding some print statements...
```
tgoodwin@cerulean:~/projects/sleeveless/fakegopath/src/k8s.io/kubernetes (siren) $ fd watch_cache.go
staging/src/k8s.io/apiserver/pkg/storage/cacher/watch_cache.go
tgoodwin@cerulean:~/projects/sleeveless/fakegopath/src/k8s.io/kubernetes (siren) $
```

Then, build your locally modified version of Kubernetes using the steps above.
Once you've pushed the image and booted a Kind cluster using the image, check the logs coming from the pod that runs the code you modified.
In this example, it'd be the API server pod. So, to see its logs...

1. quickly verify that the cluster components are running as expected
```bash
kubectl get pods --namespace kube-system
```

2. get the logs from the API server pod
```bash
kubectl logs kube-apiserver-kind-control-plane --namespace kube-system -f
```

