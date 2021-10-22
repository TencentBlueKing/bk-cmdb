### 功能描述
动态调整api限流规则

假设apiserver IP为127.0.0.1 端口为8080

### 请求参数
| 字段                |  类型       | 必选   |  描述            |
|---------------------|-------------|--------|------------------------|
|action| string | 是| 更改配置的相关动作|

#### 接口参数
(1) 查找所有限流规则，无请求参数

(2) 查找限流规则

| 字段                  |  类型        | 必选	 |  描述    |
|---------------------|-------------|--------|-----------------------------|
| rule_names           | array | 是     | 要查找的规则的名称 |

(3) 添加限流规则

| 字段     | 类型   | 必选 | 描述                                                         |
|----------|--------|------|--------------------------------------------------------------|
| rule_name | string | 是   | 策略名                                                       |
| app_code  | string | 否   | 应用ID                                                       |
| user     | string | 否   | 请求发起的用户名                                             |
| ip       | string | 否   | api的来源ip                                                  |
| method   | string | 否   | 请求的类型，配置的情况下只能为POST、GET、PUT、DELETE中的一种 |
| url      | string | 否   | api的url正则表达式                                           |
| limit    | int64  | 否   | api请求限制总次数                                            |
| ttl      | int64  | 否   | 策略存活时间，单位为秒                                       |
| deny_all  | bool   | 否   | 是否直接禁掉请求，默认为false，为true时忽略limit和ttl参数    |

(4) 删除限流规则

| 字段                  |  类型        | 必选	 |  描述    |
|---------------------|-------------|--------|----------------------------------|
| rule_names           | array | 是     | 要查找的规则的名称 |

### 请求参数示例
均为POST请求

(1) 查找所有限流规则，http://127.0.0.1:8080/set_limit_rule?action=getAll

(2) 根据名称查找限流规则，http://127.0.0.1:8080/set_limit_rule?action=get
```json
{
   "rule_names": [
      "test1",
      "test2"
  ] 
}
```

(3) 添加限流规则，http://127.0.0.1:8080/set_limit_rule?action=add
```json
{
  "rule_name": "rule1",
  "app_code": "gse",
  "user": "admin",  
  "ip": "",       
  "method": "POST",  
  "url": "^/api/v3/module/search/[^\\s/]+/[0-9]+/[0-9]+/?$",       
  "limit": 1000,    
  "ttl": 60,      
  "deny_all": false
}
```

(4) 删除限流规则，http://127.0.0.1:8080/set_limit_rule?action=delete
```json
{
   "rule_names": [
      "test1",
      "test2"
  ] 
}
```