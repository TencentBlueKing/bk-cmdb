### 描述

获取主机基础信息详情(权限：主机池主机查看权限)

### 输入参数

| 参数名称                | 参数类型   | 必选 | 描述                    |
|---------------------|--------|----|-----------------------|
| bk_supplier_account | string | 否  | 开发商账号                 |
| bk_host_id          | int    | 是  | 主机身份ID，即bk_host_id字段值 |

### 调用示例

```json
{
    "bk_host_id": 10000
}
```

### 响应示例

```json
{
    "result": true,
    "code": 0,
    "message": "success",
    "data": [
        {
            "bk_property_id": "bk_host_name",
            "bk_property_name": "host name",
            "bk_property_value": "centos7"
        },
        ......
        {
            "bk_property_id": "bk_host_id",
            "bk_property_name": "host ID",
            "bk_property_value": "10000"
        }
    ],
    "permission": null,
}
```

### 响应参数说明

| 参数名称       | 参数类型   | 描述                         |
|------------|--------|----------------------------|
| result     | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code       | int    | 错误编码。 0表示success，>0表示失败错误  |
| message    | string | 请求失败返回的错误信息                |
| permission | object | 权限信息                       |
| data       | object | 请求返回的数据                    |

#### data

| 参数名称              | 参数类型   | 描述   |
|-------------------|--------|------|
| bk_property_id    | string | 属性id |
| bk_property_name  | string | 属性名称 |
| bk_property_value | string | 属性值  |

**注意**

-
如果主机的属性字段为表格类型，返回的bk_property_value为null，要查询表格类型字段的值，请使用list_quoted_inst接口，文档链接：https://github.com/TencentBlueKing/bk-cmdb/blob/v3.12.x/docs/apidoc/cc/zh_hans/list_quoted_inst.md
