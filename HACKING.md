# Getting Started

### Hacking on Kubernetes Clients
Lets say we have some controller (e.g. Zookeeper Operator) and we want to hack on its behavior and how it interacts with the Kubernetes API.
1. Download the controller source code
```
git clone https://github.com/your-controller
```

2. Figure out how the controller interacts with Kubernetes (probably client-go, controller-runtime, or both).

3. Clone additional dependencies you need to hack on

4. In the controller's source repo, there will be a `go.mod` file which tells Go which dependencies the controller needs, and which versions of those dependencies to use. Get the controller to use our modified versions of these dependencies by adding following "replace" statements to the `go.mod` file. For example, assume we want to hack on `controller-runtime` and have the zookeeper-operator controller use this hacked version.

In the zookeeper-operator's `go.mod` file, find the line for controller runtime, and make note of the version being used. At the bottom of the file, add the following:
```
replace sigs.k8s.io/controller-runtime => /path/to/your/controller-runtime
```
and then also be sure to check out your `/path/to/your/controller-runtime` to the correct version by doing, for example, `git checkout v0.15.2`. The value for the branch name can be the version you find in the `go.mod` file.

5. Build your controller to pull in your custom changes.
NOTE: if producing a docker image, you will likely need to modify the `Dockerfile` to `COPY` in your local copy of controller-runtime (or any other dependencies you modify) so that they are available in the container's filesystem during the build process. This may also mean that you need to move your local copy to a location that is a child of the controller's root level directory.

For example, when hacking on `zookeeper-operator` I created a `/custom` folder in the zookeeper-operator repo and moved my copies of `controller-runtime` and `client-go` in there. Then, I added the following line to the `Dockerfile`:
```
# custom tim stuff
COPY custom/ custom/
```
6. Build and push your image.

7. Modify the K8s manifest for this controller
- update the image location and name (e.g. `docker.io/your-docker-username/your-controller:your-tag)
- probably want to set `imagePullPolicy: Always`

8. Repeat! That's the hacking dev loop

### Hacking on Kubernetes Source
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


#### example workflow
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

