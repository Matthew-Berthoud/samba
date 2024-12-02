Simple lightweight faas orchestrator that runs on a single machine

[repo](https://github.com/openfaas/faasd?tab=readme-ov-file)


# faasd with multipass
[video](https://www.youtube.com/watch?v=WX1tZoSXy8E)

Follow [the tutorial](https://github.com/openfaas/faasd/blob/master/docs/MULTIPASS.md).

Generate ssh key if not already done, and put it in the cloud config file in the proper spot
```shell
cd ~
ssh-keygen
# hit enter through the prompts
cat ~/.ssh/id_rsa.pub # or similar, if not rsa
```

# faasd with just qemu
__HAVENT GOTTEN THIS WORKING__

- cloudinit with ubuntu [tutorial](https://cloudinit.readthedocs.io/en/latest/tutorial/qemu.html)
- launch cloudinit with qemu [tutorial](https://cloudinit.readthedocs.io/en/latest/howto/launch_qemu.html#launch-qemu)

Make these files to fill out the genisoimage signature
```shell
touch network-config
touch meta-data
cat >user-data <<EOF
#cloud-config
password: password
chpasswd:
  expire: False
ssh_pwauth: True
EOF
```

Generate ISO disk. Use the first to do the simple cloud-init example, and use the following to do the openfaas-enabled example.
```shell
genisoimage \
    -output seed.img \
    -volid cidata -rational-rock -joliet \
    user-data meta-data network-config
```
```shell
genisoimage \
    -output seed.img \
    -volid cidata -rational-rock -joliet \
    cloud-config.txt
```

Download ubuntu image
```shell
wget https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img
```

Boot with qemu
```shell
qemu-system-x86_64 -m 1024 -net nic -net user \
    -drive file=jammy-server-cloudimg-amd64.img,index=0,format=qcow2,media=disk \
    -drive file=seed.img,index=1,media=cdrom \
    -machine accel=kvm:tcg \
    -nographic
```

