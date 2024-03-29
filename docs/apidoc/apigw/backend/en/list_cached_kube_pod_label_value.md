### Description

Query the cached pod label value list by biz id and label key (version: v3.13.5+, permission: biz access)

**Note:**
- This interface will return the deduplicated pod label value list of the specified label key under the specified business.
- This interface is a cache interface, and the default full cache refreshing time is once a day.
- If the pod data changes, the cached data corresponding to the label values in the pod will be refreshed in real time through the event mechanism.
- This interface is only used for selecting label values from the drop-down list on the front-end page, and is not recommended for use in other scenarios. If it is used in other scenarios and causes abnormal situations, you should take responsibility for the consequences.

### Parameters

| Name      | Type   | Required | Description |
|-----------|--------|----------|-------------|
| bk_biz_id | int    | yes      | Business ID |
| key       | string | yes      | label key   |

### Request Example

```json
{
  "bk_biz_id": 3,
  "key": "key1"
}
```

### Response Example

```json
{
  "result": true,
  "code": 0,
  "message": "success",
  "permission": null,
  "data": {
    "values": [
      "value1",
      "value2"
    ]
  }
}
```

### Response Parameters

| Name       | Type   | Description                                                                               |
|------------|--------|-------------------------------------------------------------------------------------------|
| result     | bool   | Whether the request was successful or not. true:request successful; false request failed. |
| code       | int    | The error code. 0 means success, >0 means failure error.                                  |
| message    | string | The error message returned by the failed request.                                         |
| data       | object | The data returned by the request.                                                         |
| permission | object | Permission information                                                                    |

#### data

| Name   | Type         | Description                                                             |
|--------|--------------|-------------------------------------------------------------------------|
| values | string array | Deduplicated pod label value list of specified label key under business |
