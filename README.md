### The official repository for this project has been moved to [LXC](https://github.com/lxc/) organization.

This repository is kept here for historical purposes only.

# Go Bindings for LXC 0.9

This package implements [Go](http://golang.org) bindings for the [LXC](http://linuxcontainers.org/) C API.

## Requirements

This package requires [LXC 0.9](https://github.com/lxc/lxc/releases) and [Go 1.x](https://code.google.com/p/go/downloads/list). 

It has been tested on 

+ Ubuntu 12.10 (quantal) by manually installing LXC 0.9 
+ Ubuntu 13.04 (raring) by using distribution [provided packages](https://launchpad.net/ubuntu/raring/+package/lxc)
+ Ubuntu 13.10 (saucy) by using distribution [provided packages](https://launchpad.net/ubuntu/saucy/+package/lxc)

## Installing

The typical `go get github.com/caglar10ur/lxc` will install LXC Go Bindings.

## Documentation

Documentation can be found at [GoDoc](http://godoc.org/github.com/caglar10ur/lxc)

## Examples

See the [examples](https://github.com/caglar10ur/lxc/tree/master/examples) directory for some.

## Notes

Note that since LXC 0.9 does not have full user namespaces support, any code using the LXC API needs to run as root.
