### 变量说明
| 变量名称                                   | 变量含义                             | 可直接设置成的变量值  |
|----------------------------------------|----------------------------------|-------------|
| BK_CMDB_ES_STATUS                      | 全文检索功能开关(取值：off/on)，默认是off，开启是on | "off"       |
| BK_COMPONENT_API_URL                   | esb地址                            |             |
| BK_CMDB_APP_CODE                       | cmdb app code                    |             |
| BK_CMDB_APP_SECRET                     | cmdb app secret                  |             |
| BKPAAS_CC_URL                          | 该值表示部署完成后,输入到浏览器中访问的cmdb 网址      |
| BKPAAS_LOGIN_URL                       | 蓝鲸统一登录访问地址                       |             |
| BK_CMDB_AUTH_SCHENE                    | 权限模式，web页面使用，可选值: internal, iam  | iam         |
| BKPAAS_BK_DOMAIN                       | 用于配置前端需要的cookie domain地址         |             |
| BK_HTTP_SCHEMA                         | 访问协议                             |             |
| BKPAAS_SHARED_RES_URL                  | 蓝鲸共享资源URL                        |             |
| BK_IAM_V3_APP_CODE                     | 权限中心app code                     |             |
| BKPAAS_NODEMAN_URL                     | 节点管理地址                           |             |
| BK_CMDB_ENABLE_BK_NOTICE               | 是否启用消息通知, true或false             |             |
| BK_CMDB_APIGW_JWT_ENABLED              | 是否通过jwt调用apigw, true或false       | true        |
| BK_CMDB_APIGW_JWT_PUBLICKEY            | cmdb API GATEWAY网关公钥             |             |
| BK_APIGW_BK_NOTICE_URL                 | 消息通知中心API GATEWAY网关地址            |             |
| BK_APIGW_BK_CMDB_URL                   | cmdb API GATEWAY网关地址             |             |
| BK_APIGW_BK_LOGIN_URL                  | 蓝鲸登录API GATEWAY网关地址              |             |
| BK_CMDB_MONGODB_HOST                   | cmdb mongodb地址                   |             |
| BK_CMDB_MONGODB_PORT                   | cmdb mongodb端口                   |             |
| BK_CMDB_MONGODB_USERNAME               | cmdb mongodb用户                   |             |
| BK_CMDB_MONGODB_PASSWORD               | cmdb mongodb密码                   |             |
| BK_CMDB_MONGODB_DATABASE               | cmdb mongodb数据库名称                | cmdb        |
| BK_CMDB_MONGODB_MAX_OPEN_CONNS         | cmdb mongodb最大连接数                | 3000        |
| BK_CMDB_MONGODB_MAX_IDLE_CONNS         | cmdb mongodb最大空闲连接数              | 100         |
| BK_CMDB_MONGODB_MECHANISM              | cmdb mongodb mechanism           | SCRAM-SHA-1 |
| BK_CMDB_MONGODB_RS_NAME                | cmdb mongodb  rsName             | rs0         |
| BK_CMDB_MONGODB_SOCKET_TIMEOUT_SECONDS | cmdb mongodb socket连接的超时时间       | 10          |
| BK_CMDB_REDIS_SENTINEL_HOST            | cmdb redis sentinel地址            |             |
| BK_CMDB_REDIS_SENTINEL_PORT            | cmdb redis sentinel端口            |             |
| BK_CMDB_REDIS_PASSWORD                 | cmdb redis密码                     |             |
| BK_CMDB_REDIS_SENTINEL_PASSWORD        | cmdb redis sentinel密码            |             |
| BK_CMDB_REDIS_DATABASE                 | cmdb redis数据库名称                  | "0"         |
| BK_CMDB_REDIS_MAX_OPEN_CONNS           | cmdb redis最大连接数                  | 3000        |
| BK_CMDB_REDIS_MAX_IDLE_CONNS           | cmdb redis最大空闲连接数                | 1000        |
| BK_CMDB_REDIS_MASTER_NAME              | cmdb redis master 名称             |             |
