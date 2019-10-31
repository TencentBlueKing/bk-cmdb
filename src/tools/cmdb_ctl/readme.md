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
  --addrport="": the ip address and port for the hosts to apply command, separated by ',', corresponding environment variable is ADDR_PORT
  --zkaddr="": the ip address and port for the zookeeper hosts, separated by ',', corresponding environment variable is ZK_ADDR
  ```
- 示例

  - ```
    ./tool_ctl log --set-default --addrport=127.0.0.1:8080 --zkaddr=127.0.0.1:2181
    ```

  - ```
    ./tool_ctl log --set-v=3 --addrport="127.0.0.1:8080" --zkaddr="127.0.0.1:2181"
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
  --zkaddr="": the ip address and port for the zookeeper hosts, separated by ',', corresponding environment variable is ZK_ADDR
  --value="": the value to be set（仅用于set命令）
  ```
- 示例

  - ```
    ./tool_ctl ls --zk-path=/test --zkaddr=127.0.0.1:2181
    ```

  - ```
    ./tool_ctl get --zk-path=/test --zkaddr=127.0.0.1:2181
    ```

  - ```
    ./tool_ctl del --zk-path=/test --zkaddr=127.0.0.1:2181
    ```

  - ```
    ./tool_ctl set --zk-path=/test --zkaddr=127.0.0.1:2181 --value=test
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
  ./tool_ctl snapshot-check [flags]
  ```

- 命令行参数
  ```
  --bizId=2: blueking business id. e.g: 2
  --zkaddr="": the ip address and port for the zookeeper hosts, separated by ',', corresponding environment variable is ZK_ADDR
  ```
- 示例

  - ```
    ./tool_ctl snapshot-check --bizId=2 --zkaddr=127.0.0.1:2181
    ```
