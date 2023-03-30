### 功能描述

新建项目(版本：v3.10.23+，权限：项目的创建权限)

### 请求参数

{{ common_args_desc }}


#### 接口参数

| 字段                       |  类型      | 必选   |  描述                                      |
|----------------------------|------------|--------|--------------------------------------------|
| data | array| 是 |数组, 一次限制创建200|

#### data

| 字段                 |  类型      | 必选   |  描述    |
|--------------------|------------|--------|----------|
| bk_project_id      |  string  | 否     | 项目id, 若传此参数，需要是32位的uuid，不带中划线的id；若不传，系统会自动生成|
| bk_project_name    |  string  | 是     | 项目名称|
| bk_project_code    |  string  | 是     | 项目英文名|
| bk_project_desc    |  string  | 否     | 项目描述|
| bk_project_type    |  enum  | 否     | 项目类型，可选值："mobile_game"(手游)、"pc_game"(端游)、"web_game"(页游)、"platform_prod"(平台产品)、"support_prod"(支撑产品)、"other"(其他)，默认值："other"|
| bk_project_sec_lvl | enum   | 否     | 保密级别，可选值："public"(公开)、"private"(私有)、"classified"(机密)，默认值："public"|
| bk_project_owner   |  string  | 是     | 项目负责人|
| bk_project_team    | array   |  否    | 所属团队|
| bk_project_icon    |  string  | 否     | 项目图标 |

### 请求参数示例

```json
{
    "bk_app_code":"esb_test",
    "bk_app_secret":"xxx",
    "bk_username":"xxx",
    "bk_token":"xxx",
    "data": [
        {
            "bk_project_id": "21bf9ef9be7c4d38a1d1f2uc0b44a8f2",
            "bk_project_name": "test",
            "bk_project_code": "test",
            "bk_project_desc": "test project",
            "bk_project_type": "mobile_game",
            "bk_project_sec_lvl": "public",
            "bk_project_owner": "admin",
            "bk_project_team": [1, 2],
            "bk_project_icon": "https://127.0.0.1/file/png/11111"
        }
    ]  
}
```

### 返回结果示例

```json

{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data": {
        "ids": [1]
    },
    "request_id": "dsda1122adasadadada2222"
}
```
**注意：**
- 返回的data中的ID数组顺序与参数中的数组数据顺序保持一致。

### 返回结果参数说明
#### response

| 名称    | 类型   | 描述                                    |
| ------- | ------ | ------------------------------------- |
| result  | bool   | 请求成功与否。true:请求成功；false请求失败 |
| code    | int    | 错误编码。 0表示success，>0表示失败错误    |
| message | string | 请求失败返回的错误信息                    |
| permission    | object | 权限信息    |
| request_id    | string | 请求链id    |
| data    | object | 请求返回的数据                           |

#### data

| 字段       | 类型      | 描述     |
|----------- |-----------|----------|
| ids |    array    |  在cc中的唯一标识数组  |
