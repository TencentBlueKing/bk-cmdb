### Functional description

batch delete kube pod (version: v3.10.23+, auth: Delete Kube Pod)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field | Type         | Required | Description                                                                     |
|-------|--------------|----------|---------------------------------------------------------------------------------|
| data  | object array | yes      | The array of pod info to be deleted, the sum of all pods in data is at most 200 |

#### data

| Field     | Type         | Required | Description                                                                       |
|-----------|--------------|----------|-----------------------------------------------------------------------------------|
| bk_biz_id | int          | yes      | biz id                                                                            |
| ids       | int array    | yes      | The array of pod cc IDs to be deleted, the sum of all pods in data is at most 200 |

### Request Parameters Example

```json
{
    "bk_app_code": "code",
    "bk_app_secret": "secret",
    "bk_username": "xxx",
    "bk_token": "xxxx",
    "data": [
        {
            "bk_biz_id": 123,
            "ids": [
                5,
                6
            ]
        }
    ]
}
```

### Return Result Example

```json
{
    "result": true,
    "code": 0,
    "message": "",
    "permission": null,
    "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Return Result Parameters Description

#### response

| Name       | Type   | Description                                                                             |
|------------|--------|-----------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                              |
| message    | string | Error message returned by request failure                                               |
| permission | object | Permission information                                                                  |
| request_id | string | Request chain id                                                                        |