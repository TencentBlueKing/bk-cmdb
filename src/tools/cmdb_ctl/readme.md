# cmdb 小工具

### 动态调整日志级别
- 使用方式

  ```
  ./tool_ctl log [flags]
  ```

- 命令行参数
  ```
  --set-default[=false]: set log level to default value
  --set-v="": set log level for V logs
  --addrport="": the ip address and port for the hosts to apply command, separated by comma
  --zk-addr="": the ip address and port for the zookeeper hosts, separated by comma, corresponding environment variable is ZK_ADDR
  ```
- 示例

  - ```
    ./tool_ctl log --set-default --addrport=127.0.0.1:8080 --zk-addr=127.0.0.1:2181
    ```

  - ```
    ./tool_ctl log --set-v=3 --addrport="127.0.0.1:8080" --zk-addr="127.0.0.1:2181"
    ```

### 优雅退出
- 使用方式

  ```
  ./tool_ctl shutdown [flags]
  ```

- 命令行参数
  ```
  --pids="": the processes to be shutdown, separated by comma
  --show-pids[=false]: show cmdb processes pid
  ```
- 示例

  - ```
    ./tool_ctl shutdown --pids=1234
    ```

  - ```
    ./tool_ctl shutdown --show-pids
    ```

### zookeeper操作
- 使用方式

  ```
  ./tool_ctl zk [command]
  ```

- 子命令
  ```
  ls          list children of specified zookeeper node
  get         get value of specified zookeeper node
  del         delete specified zookeeper node
  set         set value of specified zookeeper node
  ```
- 命令行参数
  ```
  --zk-path="": the zookeeper path
  --zk-addr="": the ip address and port for the zookeeper hosts, separated by comma, corresponding environment variable is ZK_ADDR
  --value="": the value to be set（仅用于set命令）
  ```
- 示例

  - ```
    ./tool_ctl zk ls --zk-path=/test --zk-addr=127.0.0.1:2181
    ```

  - ```
    ./tool_ctl zk get --zk-path=/test --zk-addr=127.0.0.1:2181
    ```

  - ```
    ./tool_ctl zk del --zk-path=/test --zk-addr=127.0.0.1:2181
    ```

  - ```
    ./tool_ctl zk set --zk-path=/test --zk-addr=127.0.0.1:2181 --value=test
    ```
    
### 回显服务器
- 使用方式

  ```
  ./tool_ctl echo [flags]
  ```

- 命令行参数
  ```
  --url="": the url for echo server to listen
  ```
- 示例

  - ```
    ./tool_ctl echo --url=127.0.0.1:8080/echo
    ```
### 检查主机快照
- 使用方式

  ```
  ./tool_ctl snapshot [flags]
  ```

- 命令行参数
  ```
  --bizId=2: blueking business id. e.g: 2
  --zk-addr="": the ip address and port for the zookeeper hosts, separated by comma, corresponding environment variable is ZK_ADDR
  ```
- 示例

  - ```
    ./tool_ctl snapshot --bizId=2 --zk-addr=127.0.0.1:2181
    ```

 ### 权限中心资源操作
 - 使用方式
 
   ```
   ./tool_ctl auth [command]
   ```

 - 子命令
   ```
   check       check if user has the authority to operate resources
   ```

 - 命令行参数
  ```
   -h, --help[=false]: help for auth
   -v, --logV=4: the log level of request, default, request body log level is 4
   -r, --resource="": the resource for authorize
   -f, --rsc-file="": the resource file path for authorize
  --zk-addr="": the ip address and port for the zookeeper hosts, separated by comma, corresponding environment variable is ZK_ADDR
   --supplier-account="0": the supplier id that this user belongs to（仅用于check命令）
   --user="": the name of the user（仅用于check命令）
   ```

 - 示例
 
   - ```
     ./tool_ctl auth check --app-code=test --app-secret=test --auth-address=http://127.0.0.1 --resource=[{\"type\":\"modelInstance\",\"action\":\"update\",\"Name\":\"\",\"InstanceID\":1,\"InstanceIDEx\":\"\",\"SupplierAccount\":\"\",\"business_id\":2,\"Layers\":[{\"type\":\"model\",\"action\":\"\",\"Name\":\"\",\"InstanceID\":1,\"InstanceIDEx\":\"\"}]}] --supplier-account=0 --user=test
     ```

   - ```
     ./tool_ctl auth check --app-code=test --app-secret=test --auth-address=http://127.0.0.1 --rsc-file=resource.json
     
     resource.json file content example:
     [
         {
             "type": "modelInstance",
             "action": "update",
             "Name": "",
             "InstanceID": 1,
             "InstanceIDEx": "",
             "SupplierAccount": "",
             "business_id": 2,
             "Layers": [
                 {
                     "type": "model",
                     "action": "",
                     "Name": "",
                     "InstanceID": 1,
                     "InstanceIDEx": ""
                 }
             ]
         }
     ]
     ```

### 检查拓扑结构
- 使用方式

  ```
  ./tool_ctl topo [flags]
  ```

- 命令行参数
  ```
  --bizId=2: blueking business id. e.g: 2
  --mongo-uri="": the mongodb URI, eg. mongodb://127.0.0.1:27017/cmdb, corresponding environment variable is MONGO_URI
  ```
- 示例

  - ```
    ./tool_ctl topo --bizId=2 --mongo-uri=mongodb://127.0.0.1:27017/cmdb
    ```

### 操作api请求限流策略
- 使用方式
    ```
      ./tool_ctl limiter [flags]
      ./tool_ctl limiter [command]
    ```

- 子命令
    ```
      set         set api limiter rule, use with flag --rule
      get         get api limiter rules according rule names,use with flag --rulenames
      del         del api limiter rules, use with flag --rulenames
      ls          list all api limiter rules
    ```
- 命令行参数
    ```
     --rule="": the api limiter rule to set, a json like '{"rulename":"rule1","appcode":"gse","user":"","ip":"","method":"POST","url":"^/api/v3/module/search/[^\\s/]+/[0-9]+/[0-9]+/?$","limit":1000,"ttl":60,"denyall":false}'
     --rulenames="": the api limiter rule names to get or del, multiple names is separated with ',',like 'name1,name2'
    ```

- 示例
    ```
      ./tool_ctl limiter set --rule='{"rulename":"rule1","appcode":"gse","user":"","ip":"","method":"POST","url":"^/api/v3/module/search/[^\\s/]+/[0-9]+/[0-9]+/?$","limit":1000,"ttl":60,"denyall":false}'
      ./tool_ctl limiter get --rulenames=test1,test2
      ./tool_ctl limiter del --rulenames=test1,test2
      ./tool_ctl limiter ls
    ```
### 检查yaml格式文件是否正确
- 使用方式
    ```
      ./tool_ctl checkconf [flags]
    ```
- 命令行参数
    ```
      --path="": the directory to store configuration
    ```
- 示例
    ```
      ./tool_ctl checkconf
      ./tool_ctl checkconf --path="/data/cmdb/cmdb_adminserver/configures"
    ```