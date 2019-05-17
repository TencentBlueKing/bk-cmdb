同步第三方系统框架

config file example:

```
[core]
supplierAccount=0
user=build_user
ccaddress=http://test.apiserver:8080
```

usage:

```
go run main.go --regdiscv=<zk-ip>:<zk-port> --config conf/demo.conf --addrport 127.0.0.1:8086
```
