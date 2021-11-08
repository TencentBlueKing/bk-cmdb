### 功能描述

克隆主机属性

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段        |  类型   | 必选   |  描述                       |
|-------------|---------|--------|-----------------------------|
| bk_org_ip   | string  | 是     | 源主机内网ip   |
| bk_dst_ip   | string  | 是     | 目标主机内网ip |
| bk_org_id   | int  | 是     | 源主机身份ID    |
| bk_dst_id   | int  | 是     | 目标主机身份ID |
| bk_biz_id   | int     | 是     | 业务ID                      |
| bk_cloud_id | int     | 否     | 云区域ID                    |


注： 使用主机内网IP进行克隆与使用主机身份ID进行克隆，这两种方式只能使用期中的一种，不能混用。

### 请求参数示例

```json
{
    "bk_biz_id":2,
    "bk_org_ip":"127.0.0.1",
    "bk_dst_ip":"127.0.0.2",
    "bk_cloud_id":0
}
```
或

```json
{
    "bk_biz_id":2,
    "bk_org_id": 10,
    "bk_dst_id": 11,
    "bk_cloud_id":0
}
```

### 返回结果示例

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": null
}
```
