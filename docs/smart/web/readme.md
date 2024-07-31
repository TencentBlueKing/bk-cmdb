S-Mart包构建
===========

## 目录结构

构建前需要将CMDB的代码目录调整为如下目录结构：

```
.                                       CMDB代码根目录
├── app_desc.yaml                       S-Mart包配置
├── logo.png                            CMDB的logo
├── bin                                 extra-data目录中的bin目录
│     ├── envsubst                      用于将环境变量渲染到SaaS的配置文件中
│     ├── go-pre-compile                S-Mart包编译前的前置处理脚本
│     ├── post-compile                  S-Mart包编译后的后置处理脚本
│     └── start-web.sh                  SaaS启动脚本
├── configure                           extra-data目录中的configure目录
│     ├── readme.md                     配置文件变量说明
│     └── web.yaml.tmpl                 配置文件模板
├── web                                 编译好的前端包
├── conf                                对应cmdb的resources目录
│     ├── errors                        错误码配置
│     └── language                      国际化配置
├── changelog_user                      版本日志目录，对应cmdb的docs/support-file/changelog_user目录
├── src                                 CMDB的代码目录
│     ├── web_server                    web_server的代码目录，S-Mart包构建时会使用该目录进行编译，编译好的二进制会放到bin目录中
...
```

## 构建方式

使用`蓝鲸S-Mart源码包构建工具`蓝盾流水线插件指定源码路径为CMDB代码根目录进行构建
