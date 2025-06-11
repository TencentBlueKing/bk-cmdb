### 描述

根据主机id列表和管控区域id,更新主机的管控区域字段(权限：业务主机编辑权限)

### 输入参数

| 参数名称        | 参数类型  | 必选 | 描述             |
|-------------|-------|----|----------------|
| bk_biz_id   | int   | 否  | 业务ID           |
| bk_cloud_id | int   | 是  | 管控区域ID         |
| bk_host_ids | array | 是  | 主机IDs, 最多2000个 |

### 调用示例

```json
{
    "bk_host_ids": [43, 44], 
    "bk_cloud_id": 27,
    "bk_biz_id": 1
}
```

### 响应示例

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": null
}
```

### 响应参数说明
