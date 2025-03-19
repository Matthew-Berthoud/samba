# Knative
Knative has FaaS stuff. Woohoo!

## [Tutorials](https://knative.dev/docs/getting-started/tutorial/)

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
- more here...


### End to End
Refer to these bulletpoints if you get stuck on the [end-to-end tutorial](https://knative.dev/docs/bookstore/page-0/welcome-knative-bookstore-tutorial/)

- coming soon

