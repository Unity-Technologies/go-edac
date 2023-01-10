# go-edac [![Go Report Card](https://goreportcard.com/badge/github.com/multiplay/go-edac)](https://goreportcard.com/report/github.com/multiplay/go-edac) [![License](https://img.shields.io/badge/license-BSD-blue.svg)](https://github.com/multiplay/go-edac/blob/master/LICENSE) [![GoDoc](https://pkg.go.dev/badge/github.com/multiplay/go-edac?utm_source=godoc)](https://pkg.go.dev/github.com/multiplay/go-edac)

go-edac is a [golang](http://golang.org/) library and utility for Linux kernel [Error Detection and Correction (EDAC)](https://www.kernel.org/doc/html/v4.14/admin-guide/ras.html#edac-error-detection-and-correction).

Features
--------
* Reading memory status
* Resetting memory counters
* Setting memory scrub rate

Installation
------------
```shell
go get -u github.com/multiplay/go-edac
```

Requirements
------------

In order for EDAC information to be available your system needs the required kernel modules
installed and loaded.

The following command will return as list of modules if they are:
```shell
lsmod |grep edac
```

If you don't see any modules and you have a CPU which supports ECC memory checking, you may
need a newer kernel which adds support for your specific hardware.

For example Ubuntu 16.04 LTS GA Kernel (4.4) only supports up to Intel E3-1200 v5 CPU's.

If you want to stick with 16.04 LTS and have a newer Intel CPU you can opt into the
[Rolling HWE Stacks](https://wiki.ubuntu.com/Kernel/RollingLTSEnablementStack)

```shell
apt-get install --install-recommends linux-generic-hwe-16.04
```

This will update your kernel to 4.15 which supports up to the 7th Generation Intel E3's.
