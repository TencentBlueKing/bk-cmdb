### 功能描述

根据主机id列表和云区域id,更新主机的云区域字段

### 请求参数

{{ common_args_desc }}

#### 接口参数

| 字段                 |  类型      | 必选	   |  描述                 |
|----------------------|------------|--------|-----------------------|
| bk_biz_id            | int  | 否   | 业务ID |
| bk_cloud_id         | int  | 是   | 云区域ID |
| bk_host_ids         | array  | 是   | 主机IDs, 最多2000个 |


### 请求参数示例

```python
{
	"bk_host_ids": [43, 44], 
	"bk_cloud_id": 27,
	"bk_biz_id": 1
}
```

### 返回结果示例

```python
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": ""
}
```

### 返回结果实例 - 云区域 + 内网IP 重复

```python
{
  "result": false,
  "code": 1199014,
  "message": "数据唯一性校验失败， bk_host_innerip 重复",
  "permission": null,
  "data": null
}
```

### 返回结果实例 - 一次操作主机数太多
```python
{
  "result": false,
  "code": 1199077,
  "message": "一次操作记录数超过最大限制：2000",
  "permission": null,
  "data": null
}
```

```

### 返回结果参数说明

#### response

| 名称  | 类型  | 描述 |
|---|---|---|
| result | bool | 请求成功与否。true:请求成功；false请求失败 |
| code | int | 错误编码。 0表示success，>0表示失败错误 |
| message | string | 请求失败返回的错误信息 |
| data | object | 无数据返回 |
