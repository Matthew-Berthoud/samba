# OpenFunction
[openfunction.dev](https://www.openfunction.dev)

## Build hello world function in Go
Follow [these](https://openfunction.dev/docs/concepts/function_build/#build-functions-with-the-pack-cli) and [these](https://github.com/OpenFunction/samples/tree/main/functions/knative/hello-world-go#sample-function-go) instructions, mainly the first link.

You'll have to follow the links in the first link to install the `pack` cli.
Then once the function is running you'll have to visit it with `curl http://localhost:8080/world` for it to say `Hello, world!`, or whatever other thing you want to say hello to!

## Async Function Example
- Click the Go Kafka input and HTTP output binding example on the [Create Async Functions](https://openfunction.dev/docs/getting-started/quickstarts/async-functions/) page
- [This](https://kubesphere.io/blogs/serverless-way-for-kubernetes-log-alert/) is a more detailed explanation of the example
- Uses [Kafka](#Kafka)


## Kafka
- Started looking into Kafka to maybe get a better grasp on why I can't get the above example to work. The [quickstart](https://kafka.apache.org/quickstart) has a helpful high level [video](https://www.youtube.com/watch?v=vHbvbwSEYGo)

