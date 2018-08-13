# cc rpc

## 相对 go rpc 优势

- 使用function接口而非反射, 提高调用效率
- client 支持服务发现
- client 支持连接池, 可以同时连接多个服务端
- client 支持断链重连, 而 go rpc 的client一旦连接断掉后不在重连, 调用Call会直接报错
