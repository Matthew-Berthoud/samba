# Knative
Knative has FaaS stuff. Woohoo!

## [Development](https://github.com/etclab/serving/DEVELOPMENT.md)

Run this in one terminal, leaving it open/running
```bash
# Set this environment variable one way or another
echo "export KO_DOCKER_REPO='ko.local'" >> ~/.bashrc
minikube start
minikube addons enable registry
eval $(minikube -p minikube docker-env)
kubectl port-forward --namespace kube-system service/registry 5000:80
```
In another terminal,
```bash
eval $(minikube -p minikube docker-env)
minikube tunnel
```
In another terminal, from within the `serving` repo directory:
```bash
eval $(minikube -p minikube docker-env)
kubectl apply -f ./third_party/cert-manager-latest/cert-manager.yaml
kubectl wait --for=condition=Established --all crd
kubectl wait --for=condition=Available -n cert-manager --all deployments
ko apply --selector knative.dev/crd-install=true -Rf config/core/
kubectl wait --for=condition=Established --all crd
ko apply -Rf config/core/
ko delete -f config/post-install/default-domain.yaml --ignore-not-found
ko apply -f config/post-install/default-domain.yaml
kubectl apply -f ./third_party/kourier-latest/kourier.yaml
kubectl patch configmap/config-network \
  -n knative-serving \
  --type merge \
  -p '{"data":{"ingress.class":"kourier.ingress.networking.knative.dev"}}'
```
You may be prompted to enter your computer password in the `minikube tunnel` terminal.
Look at [these instructions](https://github.com/etclab/serving/DEVELOPMENT.md) if stuck.

## [Autoscale Sample App](https://knative.dev/docs/serving/autoscaling/autoscale-go/)
Instead of doing the first step and installing serving from yaml, you can start by doing the Development tutorial above.
Then, clone the knative-docs repo, run this from within it
```bash
kubectl apply -f docs/serving/autoscaling/autoscale-go/service.yaml
kubectl get ksvc autoscale-go
```
This should show a URL that ends in `.sslip.io`, instead of `.local`.
THIS MAY TAKE A BIT TO REGISTER THOUGH, so be patient and run the command a few times.
If it doesn't work, you probably didn't run `minikube tunnel`.
This has to be running before the "default domain" yaml gets run in the Development setup above.

Here are all the commands they run in the tutorial, with IP `127.0.0.1` instead of `1.2.3.4`:
```bash
curl "http://autoscale-go.default.127.0.0.1.sslip.io?sleep=100&prime=10000&bloat=5"
hey -z 30s -c 50 \
  "http://autoscale-go.default.127.0.0.1.sslip.io?sleep=100&prime=10000&bloat=5" \
  && kubectl get pods
hey -z 60s -c 100 \
  "http://autoscale-go.default.127.0.0.1.sslip.io?sleep=100&prime=10000&bloat=5"
hey -z 60s -q 100 \
  "http://autoscale-go.default.127.0.0.1.sslip.io?sleep=10"
hey -z 60s -q 100 \
  "http://autoscale-go.default.127.0.0.1.sslip.io?sleep=1000"
hey -z 60s -q 100 \
  "http://autoscale-go.default.127.0.0.1.sslip.io?prime=40000000"
hey -z 60s -c 5 \
  "http://autoscale-go.default.127.0.0.1.sslip.io?bloat=1000"
```

## Logging
To see logs in one place, I tried to use [Fluent Bit](https://knative.dev/docs/serving/observability/logging/collecting-logs/#procedure). Some of their instructions are out of date.

I ran the initial kubernetes container config with the modified version I have in this repo, that doesn't use a deprecated nginx version.
```bash
kubectl apply -f ./log-collector.yaml
```

Didn't seem to work tho, just gave an nginx Not Found screen.

I'm trying to chang the code and rebuild and it's not working, not sure why.


## Where do we code?
- Current thoughts:
    - Client: queue-proxy sidecars to the functions will be the "client" for proxy-reencryption. They will generate a re-encryption key on startup and send it to the "server".
    - Server: I'm still a little uncertain on what component exactly launches new pods for scaling.
      It's something in the "deployment", but I'm getting worried that it's some underlying kubernetes thing that we'd have to modify.
- image on [this page](https://knative.dev/docs/serving/request-flow/) is helpful conceptually
- [description of scaling architecture](https://knative.dev/docs/serving/request-flow/#scale-from-zero)
- It looks like the autoscaler outputs just a number of pods to scale to, and actual pod launching is accomplished by the Deployment, within the revision. See [this page's](https://github.com/knative/serving/blob/main/docs/scaling/SYSTEM.md) diagrams
- Some of [this](https://github.com/knative/serving/blob/main/docs/encryption/knative-encryption.md) encryption stuff might be useful
- [video, where second half explains custom autoscaling](https://www.youtube.com/watch?v=OPSIPr-Cybs)
    - talks about [this customization](https://knative.dev/docs/serving/autoscaling/autoscale-go/#customization) which won't be enough for us, we'll have to modify whatever actually launches the new containers

## [Knative Overview Tutorials](https://knative.dev/docs/getting-started/tutorial/)

### Quickstart
Refer to these bulletpoints if you get stuck on the [quickstart tutorial](https://knative.dev/docs/getting-started/)
- I downloaded the various binaries from source, using the `amd64` version since I'm on the shs machines.
You can use `wget LINK_TO_ASSET` for command-line downloading.
- On [this page](https://knative.dev/docs/getting-started/install-func/) I chose the first option (`func` CLI).
I think Collin mentioned doing the second option: installing it as a `kn` plugin.
- On [this page](https://knative.dev/docs/getting-started/build-run-deploy-func/#procedure) ran:
```
cd hello
export FUNC_REGISTRY=<your-docker-username>
func run
```
- If you set that environment variable you can omit all the `--registry=<your-docker-username>` arguments
- If you run into issues with `func deploy`, try `docker logout && docker login`, and then try again
- Before starting on [Deploying a Knative Service](https://knative.dev/docs/getting-started/first-service/) you should run `func delete hello` since the next exmample function will also be called `hello`
- __Can't get port forwarding to work__ for [this dashboard](https://knative.dev/docs/getting-started/first-source/#examining-the-cloudevents-player) because of javascript errors I think
- Next little bit doesn't work either for same reason
- Clean up worked fine. DONE

### End to End
I haven't done the [end-to-end tutorial](https://knative.dev/docs/bookstore/page-0/welcome-knative-bookstore-tutorial/)


