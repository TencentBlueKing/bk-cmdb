# snapshot check server



### run parameters 

* regdiscv cmdb 服务配置regdiscv, 服务发现组件的地址

* appID  蓝鲸业务的id， 非必填， 默认值为2

* interval 检查服务的间隔，单位是分钟， 非必填，默认值10，最小值是10分钟

* log-dir 日志目录


eg :


./cmdb_tool_snapshotcheck --regidscv=127.0.0.1:2181 --log-dir=./log 
