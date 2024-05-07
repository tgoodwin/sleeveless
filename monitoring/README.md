# tracey
tracing util with tempo and grafana

### install
run `bootstrap.sh`

test that the plumbing has been configured correctly: `kubectl apply -f trace-generator.yaml`

### k8s version compatibility
- works on v1.19.11
- should work on higher versions too
