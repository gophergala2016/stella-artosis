Stella Artosis
--------------

![Stella](artosis.jpg)

Stella Artosis is a simple program used to quickly hash files on your
filesystem for use in forensic analysis.

Usage:

```
go get github.com/gophergala2016/stella-artosis
go install github.com/gophergala2016/stella-artosis/cmd/artosis
export GOMAXPROCS=`sysctl -n hw.ncpu`
artosis -include /usr/lib
```
