### Functional description

Query instance Association topology

### Request Parameters

{{ common_args_desc }}

#### Interface Parameters

| Field                | Type   | Required| Description|
| ------------------- | ------ | ---- | ---- |
| bk_obj_id           |  string |yes   | Model id   |
| bk_inst_id          |  int    | yes | Instance id   |


### Request Parameters Example

``` python
{
    "bk_app_code": "esb_test",
    "bk_app_secret": "xxx",
    "bk_username": "xxx",
    "bk_token": "xxx",
    "bk_obj_id":"test",
    "bk_inst_id":1
}
```


### Return Result Example

```python
{
    "result": true,
    "code": 0,
    "data": [
        {
            "id": "",
            "bk_obj_id": "biz",
            "bk_obj_icon": "icon-cc-business",
            "bk_inst_id": 0,
            "bk_obj_name": "business",
            "bk_inst_name": "",
            "asso_id": 0,
            "count": 1,
            "children": [
                {
                    "id": "6",
                    "bk_obj_id": "biz",
                    "bk_obj_icon": "icon-cc-business",
                    "bk_inst_id": 6,
                    "bk_obj_name": "business",
                    "bk_inst_name": "",
                    "asso_id": 558
                }
            ]
        }
    ],
    "message": "success",
    "permission": null,
    "request_id": "94c85fdf6a9341e18750a44d6e18c127"
}
```

### Return Result Parameters Description
#### response

| Name    | Type   | Description                                       |
| ------- | ------ | ------------------------------------------ |
| result  | bool   | Whether the request was successful or not. True: request succeeded;false request failed|
| code    |  int    | Wrong code. 0 indicates success,>0 indicates failure error    |
| message | string |Error message returned by request failure                     |
| permission    |  object |Permission information    |
| request_id    |  string |Request chain id    |
| data    |  object |Data returned by request                             |

#### data

| Field         | Type         | Description                           |
| ------------ | ------------ | ------------------------------ |
| bk_inst_id   |  int          | Instance ID                         |
| bk_inst_name | string       | The name the instance is used to present             |
| bk_obj_icon  | string       | The name of the model icon                 |
| bk_obj_id    |  string       | Model ID                         |
| bk_obj_name  | string       | The name the model is used to present             |
| children     |  object array |The set of all associated instances in this model|
| count        |  int          | Children contains the number of nodes        |

#### children

| Field        | Type   | Description               |
|-------------|--------|--------------------|
|bk_inst_id   |  int    | Instance ID            |
|bk_inst_name | string |The name the instance is used to present|
|bk_obj_icon  | string |The name of the model icon     |
|bk_obj_id    |  string |Model ID             |
|bk_obj_name  | string |The name the model is used to present|
|asso_id  | string |Association id|
