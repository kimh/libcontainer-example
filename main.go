package main

import (
	"fmt"
	"os"
	"runtime"
	"syscall"
	"github.com/opencontainers/runc/libcontainer"
	"github.com/opencontainers/runc/libcontainer/configs"
	"github.com/Sirupsen/logrus"
)

	func init() {
		if len(os.Args) > 1 && os.Args[1] == "init" {
			runtime.GOMAXPROCS(1)
			runtime.LockOSThread()
			factory, _ := libcontainer.New("")
			if err := factory.StartInitialization(); err != nil {
				logrus.Fatal(err)
			}
			panic("--this line should have never been executed, congratulations--")
		}
	}

func main() {
	rootfs := os.Getenv("ROOTFS")
	
	factory, err := libcontainer.New(rootfs,libcontainer.Cgroupfs)

	if err != nil {
		fmt.Println(err)
	}

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	config := &configs.Config{
		Rootfs: rootfs,
		Capabilities: []string{
			"CAP_CHOWN",
			"CAP_DAC_OVERRIDE",
			"CAP_FSETID",
			"CAP_FOWNER",
			"CAP_MKNOD",
			"CAP_NET_RAW",
			"CAP_SETGID",
			"CAP_SETUID",
			"CAP_SETFCAP",
			"CAP_SETPCAP",
			"CAP_NET_BIND_SERVICE",
			"CAP_SYS_CHROOT",
			"CAP_KILL",
			"CAP_AUDIT_WRITE",
		},
		Namespaces: configs.Namespaces([]configs.Namespace{
			{Type: configs.NEWNS},
			{Type: configs.NEWUTS},
			{Type: configs.NEWIPC},
			{Type: configs.NEWPID},
			{Type: configs.NEWNET},
		}),
		Cgroups: &configs.Cgroup{
			Name:   "test-container",
			Parent: "system",
			Resources: &configs.Resources{
				MemorySwappiness: -1,
				AllowAllDevices:  false,
				AllowedDevices:   configs.DefaultAllowedDevices,
			},
		},
		MaskPaths: []string{
			"/proc/kcore",
		},
		ReadonlyPaths: []string{
			"/proc/sys", "/proc/sysrq-trigger", "/proc/irq", "/proc/bus",
		},
		Devices:  configs.DefaultAutoCreatedDevices,
		Hostname: "testing",
		Mounts: []*configs.Mount{
			{
				Source:      "proc",
				Destination: "/proc",
				Device:      "proc",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "tmpfs",
				Destination: "/dev",
				Device:      "tmpfs",
				Flags:       syscall.MS_NOSUID | syscall.MS_STRICTATIME,
				Data:        "mode=755",
			},
			{
				Source:      "devpts",
				Destination: "/dev/pts",
				Device:      "devpts",
				Flags:       syscall.MS_NOSUID | syscall.MS_NOEXEC,
				Data:        "newinstance,ptmxmode=0666,mode=0620,gid=5",
			},
			{
				Device:      "tmpfs",
				Source:      "shm",
				Destination: "/dev/shm",
				Data:        "mode=1777,size=65536k",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "mqueue",
				Destination: "/dev/mqueue",
				Device:      "mqueue",
				Flags:       defaultMountFlags,
			},
			{
				Source:      "sysfs",
				Destination: "/sys",
				Device:      "sysfs",
				Flags:       defaultMountFlags | syscall.MS_RDONLY,
			},
		},
		Networks: []*configs.Network{
			{
				Type:    "loopback",
				Address: "127.0.0.1/0",
				Gateway: "localhost",
			},
		},
		Rlimits: []configs.Rlimit{
			{
				Type: syscall.RLIMIT_NOFILE,
				Hard: uint64(1025),
				Soft: uint64(1025),
			},
		},
	}

	container, err := factory.Create("abcde", config)

	if err != nil {
		fmt.Println(err)
	}

	process := &libcontainer.Process{
		Args:   []string{"/bin/sh"},
		Env:    []string{"PATH=/bin"},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = container.Start(process)
	if err != nil {
		fmt.Println(err)
	}

	// wait for the process to finish.
	status, err := process.Wait()
	
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(status)

	container.Destroy()
	
	fmt.Println("done")
}

