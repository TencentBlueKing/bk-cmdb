### Functional description

delete container clusters in batches (v3.12.1+, permission: kube cluster deletion permission)

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field     | Type  | Required | Description                                        |
|-----------|-------|----------|----------------------------------------------------|
| bk_biz_id | int   | yes      | business id to which the container cluster belongs |
| ids       | array | yes      | ids of the container cluster in cmdb               |

**Note:**

- user needs to ensure that there are no associated resources (such as namespace, pod, node workload, etc.) under the
  cluster to be deleted, otherwise the deletion will fail.- uid and id cannot be empty at the same time, nor can they be
  filled in at the same time.
- the number of clusters to be deleted at one time cannot exceed 10.

### Request Parameters Example

```json
{
  "bk_app_code": "esb_test",
  "bk_app_secret": "xxx",
  "bk_username": "xxx",
  "bk_token": "xxx",
  "bk_biz_id": 2,
  "ids": [
    1,
    2
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
  "data": null,
  "request_id": "87de106ab55549bfbcc46e47ecf5bcc7"
}
```

### Return Result Parameters Description

#### response

| Name       | Type   | Description                                                                        |
|------------|--------|------------------------------------------------------------------------------------|
| result     | bool   | Whether the request succeeded or not. True: request succeeded;false request failed |
| code       | int    | Wrong code. 0 indicates success,>0 indicates failure error                         |
| message    | string | Error message returned by request failure                                          |
| permission | object | Permission information                                                             |
| data       | object | Data returned by request                                                           |
| request_id | string | Request chain id                                                                   |
