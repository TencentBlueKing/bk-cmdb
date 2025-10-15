# cmdb 小工具

### 国际化：新增语言/新增内容自动翻译
- 使用方式

  ```
  ./translate_ctl [flags]
  ```

  - 命令行参数
    ```
    --path                                                  :  设置待翻译文件路径，文件组织形式需符合规范
    --lang              [可选值：""/任意语言类型]              :  新增语言时配置，代表新加入语言
    ```
  - 示例

    ```
      ./translate_ctl --path=../../resource/translations --lang=ko
      ```