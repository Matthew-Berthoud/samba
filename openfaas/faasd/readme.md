Simple lightweight faas orchestrator that runs on a single machine

[repo](https://github.com/openfaas/faasd?tab=readme-ov-file)



# faasd with multipass
[video](https://www.youtube.com/watch?v=WX1tZoSXy8E)

Follow [the tutorial](https://github.com/openfaas/faasd/blob/master/docs/MULTIPASS.md). Here's what I did.

```shell
curl -sSLO https://raw.githubusercontent.com/openfaas/faasd/master/cloud-config.txt
brew install multipass
```

Note: if you manually change the ssh key in the `cloud-config.txt`, launch with
```shell
multipass launch --name faasd --cloud-init cloud-config.txt
```

```shell
multipass info faasd
```

```shell
```

```shell
```

```shell
```

```shell
```
