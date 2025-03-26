# Knative
Knative has FaaS stuff. Woohoo!

## Where do we code?
[description of scaling architecture](https://knative.dev/docs/serving/request-flow/#scale-from-zero)

## [Development](https://github.com/etclab/serving/DEVELOPMENT.md)
These steps are required to get the local registry working.
The rest of [these instructions](https://github.com/etclab/serving/DEVELOPMENT.md) worked and got me set up locally on Mac.
```
minikube start
# https://minikube.sigs.k8s.io/docs/handbook/registry/
minikube addons enable registry
# I ran the following but it didn't work by itself, and then the other commands made things work so I'm not sure if it's necessary
# docker run --rm -it --network=host alpine ash -c "apk add socat && socat TCP-LISTEN:5000,reuseaddr,fork TCP:$(minikube ip):5000"
kubectl port-forward --namespace kube-system service/registry 5000:80
eval $(minikube docker-env)
```

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
Refer to these bulletpoints if you get stuck on the [end-to-end tutorial](https://knative.dev/docs/bookstore/page-0/welcome-knative-bookstore-tutorial/)

- coming soon

