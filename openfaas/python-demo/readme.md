# Python HTTP app demo

These instructions came from this [video tutorial](https://www.youtube.com/watch?v=igv9LRPzZbE).

Install `faas-cli`.
```shell
curl -sSL https://get.arkade.dev | sudo -E sh
```
```shell
arkade get faas-cli
```

Follow the instructions it provides after this command runs to add the binary to your path.

Show availible "templates" for FaaS applications. Templates are essentially Dockerfiles with some metadata.
```shell
faas-cli template store list
```

Pull in the python http app template.
```shell
faas-cli template store pull python3-http
```

Make a new Function for this demo.
```shell
faas-cli new --lang python3-http --prefix etclab openfaas-python-demo
```
