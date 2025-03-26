# Knative
Knative has FaaS stuff. Woohoo!

## [Autoscale Sample App](https://knative.dev/docs/serving/autoscaling/autoscale-go/)
Instead of doing the first step and installing serving from yaml, you can start by doing the Development tutorial
instead including installing Kourier as a network layer.
When you get the running this command: `kubectl get ksvc autoscale-go`, if it shows a URL that ends in `.local` instead of `.sslip.io`, run `minikube tunnel`.
This is another command, like the docker port forward one in Development, that you'll need to keep a terminal open for so it can run in the foreground.

## Where do we code?
- [description of scaling architecture](https://knative.dev/docs/serving/request-flow/#scale-from-zero)
- [video, where second half explains custom autoscaling](https://www.youtube.com/watch?v=OPSIPr-Cybs)

## [Development](https://github.com/etclab/serving/DEVELOPMENT.md)
These steps are required to get the local registry working.
The rest of [these instructions](https://github.com/etclab/serving/DEVELOPMENT.md) worked and got me set up locally on Mac.
```
minikube start
# https://minikube.sigs.k8s.io/docs/handbook/registry/
minikube addons enable registry
# The following you'll have to run in another terminal pane, unless you want it in the background
kubectl port-forward --namespace kube-system service/registry 5000:80
eval $(minikube docker-env)
```
Don't forget to install Kourier as a network layer

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


