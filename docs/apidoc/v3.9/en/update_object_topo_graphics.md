### Functional description

update object topology graphics

### Request Parameters

{{ common_args_desc }}

#### General Parameters

| Field          |  Type      | Required	   |  Description                                           |
|---------------|------------|--------|-------------------------------------------------|
| action        | string     | Yes     | Update methods, optional: update,override                    |
| scope_type    | string     | Yes     | Graphics range type,global,biz,cls(global) |
| scope_id      | string     | Yes     | ID of graphics range type, if it'sglobal, fill in '0'           |
| node_type     | string     | Yes     | Node type, obj, inst                           |
| bk_obj_id     | string     | Yes     | Object ID                                    |
| bk_inst_id    | int        | Yes     | Instance ID                                          |
| position      | string     | No     | The position of node in graphic                                |
| ext           | object     | No     | Front extension field                                    |
| bk_obj_icon   | string     | No     | Object icon                                  |

> scope_type,scope_id  Determine a graphic only
> node_type,bk_obj_id,bk_inst_id Three of them determines a node of each graphic, so it must be filled


### Request Parameters Example

```python

{
    "action": "update",
    "scope_tpye": "global",
    "scope_id": "0",
    "node_type": "obj",
    "bk_obj_id": "switch",
    "bk_inst_id": 0,
    "position": {
        "x": 100,
        "y": 100
    },
    "ext": {
        "a":"test",
        "b":"test"
    },
    "bk_obj_icon": "icon-cc-switch2",
}

```


### Return Result Example

```python

{
    "result": true,
    "code": 0,
    "message": "",
    "data": "success"
}
```
