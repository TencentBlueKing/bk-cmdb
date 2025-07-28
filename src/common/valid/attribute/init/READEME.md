### 说明

后续可以将插件的注入放到改目录下

##### 注意事项

- package 名字，必须为 init
- 不允许写业务逻辑

```
package init
import (
 _ "configcenter/src/common/valid/attribute/plugins"
)

```
