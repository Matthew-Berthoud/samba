Simple lightweight faas orchestrator that runs on a single machine

[repo](https://github.com/openfaas/faasd?tab=readme-ov-file)



# faasd with multipass
[video](https://www.youtube.com/watch?v=WX1tZoSXy8E)

Follow [the tutorial](https://github.com/openfaas/faasd/blob/master/docs/MULTIPASS.md) to use multipass. Here's what I did to just use raw qemu, without multipass:

# faasd with qemu
- cloudinit with ubuntu [tutorial](https://cloudinit.readthedocs.io/en/latest/tutorial/qemu.html)
- launch cloudinit with qemu [tutorial](https://cloudinit.readthedocs.io/en/latest/howto/launch_qemu.html#launch-qemu)


1. Create 20G qemu image
```shell
qemu-img create -f qcow2 faasd.img 20G
```

2. Get a latest cloud-specific Ubuntu LTS release
```shell
wget 'https://cloud-images.ubuntu.com/minimal/releases/noble/release/ubuntu-24.04-minimal-cloudimg-arm64.img'
```

3. Get the faasd config file
```shell
curl -sSLO https://raw.githubusercontent.com/openfaas/faasd/master/cloud-config.txt
```

4. Generate ssh key if not already done, and put it in the cloud config file in the proper spot
```shell
cd ~
ssh-keygen
# hit enter through the prompts
cat ~/.ssh/id_rsa.pub # or similar, if not ssh
# now replace the key in cloud-config.txt
```

4. Create ISO from the cloud config
```shell
genisoimage -output cloud-init.iso -volid cidata -joliet -rock user-data cloud-config.txt
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
