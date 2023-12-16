# TODO

## Building Kubernetes
1. Download Kubernetes
```bash
mkdir -p fakegopath/src/k8s.io

git clone --single-branch --branch $K8S_VERSION https://github.com/kubernetes/kubernetes/git fakegopath/src/k8s.io/kubernetes >> /dev/null
```

2. (Optional) install custom lib for kubernetes
TODO: when we have custom golang modules, use this step to modify the k8s `go.mod` file(s)
so our code can be incorporated as a dependency when building kubernetes in step 4.
```bash
APISERVER_MOD_FILE="fakegopath/src/k8s.io/kubernetes/staging/src/k8s.io/apiserver/go.mod"

```

3. (future step) instrument kubernetes
Eventually we may have a script to modify fresh kubernetes source code with our changes.
```bash
# go mod tidy
# go build
# ./instrumentation config.json
```

4. Build kubernetes source into a Kind node image
```bash
cd fakegopath/src/k8s.io/kubernetes

GOPATH=$/fakegopath KUBE_GIT_VERSION=siren-${K8S_VERSION} kind build node-image
cd -
docker image tag kindest/node:latest <container_registry>/node:<image_tag>

# push to registry
docker push <container_registry>/node:<image_tag>
```

5. Running modified Kubernetes in a Kind cluster locally
```bash
kind create cluster --image <container_registry>/node:<image_tag>
```

