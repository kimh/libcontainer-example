Use Go.1.5 and install all dependencies of `main.go` in `GOPATH`.

```
# You need to become root first
sudo -s

# Extract rootfs of busybox
mkdir rootfs && tar -C rootfs -xf busybox.tar

# Start container with /bin/sh
ROOTFS=$(pwd)/rootfs go run main.go
```

You should be able to login `/bin/sh` if successful.
