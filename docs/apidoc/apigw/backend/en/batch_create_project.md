### Description

Create a new project (Version: v3.10.23+, Permission: Project creation permission)

### Parameters

| Name | Type  | Required | Description                          |
|------|-------|----------|--------------------------------------|
| data | array | Yes      | Array, limit to create 200 at a time |

#### data

| Name               | Type   | Required | Description                                                                                                                                                                                                           |
|--------------------|--------|----------|-----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| bk_project_id      | string | No       | Project ID, if this parameter is passed, it needs to be a 32-character uuid without hyphens; if not passed, the system will generate it automatically                                                                 |
| bk_project_name    | string | Yes      | Project name                                                                                                                                                                                                          |
| bk_project_code    | string | Yes      | Project code                                                                                                                                                                                                          |
| bk_project_desc    | string | No       | Project description                                                                                                                                                                                                   |
| bk_project_type    | enum   | No       | Project type, optional values: "mobile_game" (mobile game), "pc_game" (PC game), "web_game" (web game), "platform_prod" (platform product), "support_prod" (support product), "other" (other), default value: "other" |
| bk_project_sec_lvl | enum   | No       | Confidentiality level, optional values: "public" (public), "private" (private), "classified" (confidential), default value: "public"                                                                                  |
| bk_project_owner   | string | Yes      | Project owner                                                                                                                                                                                                         |
| bk_project_team    | array  | No       | Team it belongs to                                                                                                                                                                                                    |
| bk_project_icon    | string | No       | Project icon                                                                                                                                                                                                          |

### Request Example

```json
{
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

### Response Example

```json
{
    "result":true,
    "code":0,
    "message":"success",
    "permission":null,
    "data": {
        "ids": [1]
    },
}
```

**Note:**

- The order of the IDs array in the returned data is consistent with the order of the array data in the parameters.

### Response Parameters

| Name       | Type   | Description                                                       |
|------------|--------|-------------------------------------------------------------------|
| result     | bool   | Whether the request was successful. true: success; false: failure |
| code       | int    | Error code. 0 indicates success, >0 indicates a failure error     |
| message    | string | Error message returned for a failed request                       |
| permission | object | Permission information                                            |
| data       | object | Data returned for the request                                     |

#### data

| Name | Type  | Description                   |
|------|-------|-------------------------------|
| ids  | array | Unique identifier array in cc |
