### Function description

batch create project (version: v3.10.23+, permission: creation permission of the project)

### Request parameters

{{ common_args_desc }}


#### Interface parameters

| field | type | required | description |
| ----------------------------|------------|----------|--------------------------------------------|
| data | array| yes      | array, limited to 200 at a time|

#### data

| field | type | required | description                                                                                                                           |
|--------------------|------------|----------|---------------------------------------------------------------------------------------------------------------------------------------|
| bk_project_id | string | no       | project_id, if pass this parameter, it needs to be a 32-bit uuid without underscore; if not, it will be automatically generated       |
| bk_project_name | string | yes      | project_name                                                                                                                          |
| bk_project_code | string | yes      | project english name                                                                                                                  |
| bk_project_desc | string | no       | project_description                                                                                                                   |
| bk_project_type | enum | no       | project type, optional values: "mobile_game", "pc_game", "web_game", "platform_prod", "support_prod", "other", default value: "other" |
| bk_project_sec_lvl | enum | no       | confidentiality level, optional values: "public", "private", "classified", default: "public"                                          |
| bk_project_owner | string | yes      | project owner                                                                                                                         |
| bk_project_team | array | no       | project team                                                                                                                          |
| bk_project_icon | string | no       | project icon                                                                                                                          |

### Request parameter examples

```json
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
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

### Return Result Example

```json
{
    "result":true,
    "code":0,
    "message": "success",
    "permission":null,
    "data": {
        "ids": [1]
    },
    "request_id": "dsda1122adasadadada2222"
}
```
**Note:**
- The order of the array of ids in the returned data remains the same as the order of the array data in the parameters.

### Return result parameter description
#### response

| name | type | description |
| ------- | ------ | ------------------------------------- |
| result | bool | The success or failure of the request. true: the request was successful; false: the request failed.|
| code | int | The error code. 0 means success, >0 means failure error.|
| message | string | The error message returned by the failed request.|
| permission | object | Permission information |
| request_id | string | request_chain_id |
| data | object | data returned by the request|

#### data

| field | type | description                       |
| ----------- |----------|-----------------------------------|
| ids | array  | array of unique identifiers in cc |
