# OpenFunction
[openfunction.dev](https://www.openfunction.dev)

_Wed Mar 19_
## Async pubsub example
- https://github.com/OpenFunction/samples/tree/main/functions/async/pubsub
- The three [prereq](https://github.com/OpenFunction/samples/tree/main/functions/async/pubsub#prerequisites)
links are very similar / the same to the [async function example](#Async Function Example) section I tried earlier.
See that for the files I made to shortcut the process.
In short,
    1. do the helm stuff for openfunction
    2. Do the helm stuff for kafka, then do  `kubectl apply -f kafka-server-config.yaml` 
    3. run `bash registry-credential.sh` to get the local registry configured
- Run the rest of the commands from the instructions from within the `samples/functions/async/pubsub` directory

```
cd producer
```
```
docker build -t matthewberthoud/v1beta1-autoscaling-producer:latest -f Dockerfile.producer .
```
```
docker push matthewberthoud/v1beta1-autoscaling-producer:latest
```



_Tue Mar 18_

## Building Openfunction 

Sources:
- https://kubernetes.io/docs/tutorials/hello-minikube/
- https://www.docker.com/blog/how-to-use-your-own-registry-2/

I made some changes to the Openfunction Dockerfile.
I got rid of the chinese GOPROXY stuff since I don't think it's necessary and it was hanging for a while.
I also got rid of the `-a` flag so that it doesn't disable caching of dependencies.
I don't think these changes are the source of my difficulties.

Run the following from within the Openfunction directory (forked from their repo)
```
docker run -d -p 5000:5000 --name registry registry:latest
```
```
minikube start --insecure-registry localhost:5000
```
```
docker build -t openfunction-local .
```
```
docker tag openfunction-local localhost:5000/openfunction-local
```
```
docker push localhost:5000/openfunction-local
```
```
docker rmi localhost:5000/openfunction-local
```
```
docker pull localhost:5000/openfunction-local
```
```
kubectl create deployment of-node --image=localhost:5000/openfunction-local
```
```
kubectl get deployments
```
```
kubectl get pods
```

Still have `ImagePullBackOff` status. More debugging to be done.



_Earlier_

## Kubernetes
`minikube start`
- I run into issues whenever I stop and start minikube, since I don't always know how to gracefully stop and start the kubernetes processes within. To remove all pods and namespaces and start fresh, I just run `minikube delete` followed by `minikube start`. Not the greatest solution but it works for now, and keeps the state consistent-ish.

## Helm
Helm is the "package manager for kubernetes." You use it to add "charts" (config files, essentially) and repos to a kubernetes project. I installed it on `shs1` with the commands in the "from script" section of [this link](https://helm.sh/docs/intro/install/#from-script).

## Build hello world function in Go
Followed and modified [these](https://openfunction.dev/docs/concepts/function_build/#build-functions-with-the-pack-cli) and [these](https://github.com/OpenFunction/samples/tree/main/functions/knative/hello-world-go#sample-function-go) instructions.
You'll have to follow the links in the first link to install the `pack` cli.

```bash
cd samples/functions/knative/hello-world-go
pack build func-helloworld-go --builder openfunction/builder-go:v2.4.0-1.17 --env FUNC_NAME="HelloWorld"  --env FUNC_CLEAR_SOURCE=true
docker run --rm --env="FUNC_CONTEXT={\"name\":\"HelloWorld\",\"version\":\"v1.0.0\",\"port\":\"8080\",\"runtime\":\"Knative\"}" --env="CONTEXT_MODE=self-host" --name func-helloworld-go -p 8080:8080 func-helloworld-go
curl http://localhost:8080/world # to say "Hello, world!"
```

## Async Function Example
### Setup
- Click the Go Kafka input and HTTP output binding example on the [Create Async Functions](https://openfunction.dev/docs/getting-started/quickstarts/async-functions/) page
- [This](https://kubesphere.io/blogs/serverless-way-for-kubernetes-log-alert/) is a more detailed explanation of the example
- Uses [Kafka](#Kafka)
- When you're doing these [prerequisites](https://github.com/OpenFunction/samples/blob/main/Prerequisites.md#openfunction), the kafka version they have in step 2 is out of date and will cause errors, and fail to launch kafka and zookeeper.
So, for step 2 just run the `kafka-server-config.yaml` file in this repo with the following command:
```bash
kubectl apply -f kafka-server-config.yaml
```
- (I discovered this error by realizing they weren't running in step 3, and then running `kubectl describe kafka kafka-server -n default`, which will show the version error in its output)
- Step 3 should work as it says. Just note that for the third command, within the pod, you'll have to run the following, since `kafka-server` is the very creative name I gave the `<kafka-server>`
```bash
kafkacat -L -b kafka-server-kafka-brokers:9092
```
- Note that later when following the article, they name their server `kafka-logs-receiver` so you can't directly copy all commands
- For the [Registry Credential](https://github.com/OpenFunction/samples/blob/main/Prerequisites.md#registry-credential) section you can just use mine by running the file in this repo
```bash
bash regitry-credential.sh
```
- To do the rest of this tutorial you'll have to be in the correct folder of the samples submodule I put in this repo
```bash
cd samples/functions/async/logs-handler-function
```

### Kubesphere Article
- Follow the instructions in [this article](https://kubesphere.io/blogs/serverless-way-for-kubernetes-log-alert/) to see the example in action.
- STILL WORKING THRU


## Kafka
- Started looking into Kafka to maybe get a better grasp on why I can't get the above example to work. The [quickstart](https://kafka.apache.org/quickstart) has a helpful high level [video](https://www.youtube.com/watch?v=vHbvbwSEYGo)

