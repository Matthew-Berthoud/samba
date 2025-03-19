# Knative
Knative has FaaS stuff. Woohoo!

Tutorials: [https://knative.dev/docs/getting-started/tutorial/](https://knative.dev/docs/getting-started/tutorial/)

## Quickstart
- I downloaded the various binaries from source, using the `amd64` version since I'm on the shs machines.
You can use `wget LINK_TO_ASSET` for command-line downloading.
- On [this page](https://knative.dev/docs/getting-started/install-func/) I chose the first option (`func` CLI).
I think Collin mentioned doing the second option: installing it as a `kn` plugin.
- On [this page](https://knative.dev/docs/getting-started/build-run-deploy-func/#procedure) I ran:
```
cd hello
export FUNC_REGISTRY=matthewberthoud
func run
```

## End to End

