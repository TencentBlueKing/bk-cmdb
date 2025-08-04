### 说明

后续将插件的import放到目录下

##### 注意事项

- package 名字，必须为 init
- 不允许写业务逻辑

```
package init
import (
 _ "configcenter/src/common/valid/attribute/plugins"
)

```
