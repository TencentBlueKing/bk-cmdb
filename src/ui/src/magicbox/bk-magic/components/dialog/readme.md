<script>
    export default {
        data () {
            return {
                preventCloseByConfirm: true,
                defaultSetting: {
                    isShow: false
                },
                basicCustomSetting: {
                    isShow: false,
                    width: 800,
                    title: '自定义配置对话框',
                    content: '简单配置的对话框',
                    closeIcon: false
                },
                basicCustomSetting2: {
                    hasHeader: false,
                    isShow: false
                },
                basicCustomSetting3: {
                    isShow: false
                },
                textareaVal: ''
            }
        },
        methods: {
            defaultDialog () {
                this.defaultSetting.isShow = true
            },
            basicCustomDialog () {
                this.basicCustomSetting.isShow = true
            },
            basicCustomDialog2 () {
                this.basicCustomSetting2.isShow = true
            },
            basicCustomDialog3 () {
                this.basicCustomSetting3.isShow = true
            },
            confirmFn (done) {
                if (!this.textareaVal) {
                    this.$bkInfo({
                        title: '请填写内容',
                        theme: 'danger'
                    })
                } else {
                    this.$bkInfo('你点击了确定')
                    done()
                }
            },
            cancelFn () {
                this.$bkInfo('你点了取消')
            }
        }
    }
</script>

<style>
    .demo-box .bk-switcher{
        margin-right: 20px;
    }
</style>

## Dialog 对话框

可完全定制内容的弹窗

### 基本用法

:::demo 默认配置的对话框
```html
<template>
    <bk-button
        :type="'primary'"
        @click="defaultDialog">
        默认配置
    </bk-button>

    <bk-dialog
        :is-show.sync="defaultSetting.isShow"
        @confirm="done => done()">
    </bk-dialog>
</template>

<script>
    export default {
        date () {
            return {
                defaultSetting: {
                    isShow: false
                }
            }
        },
        methods: {
            defaultDialog () {
                this.defaultSetting.isShow = true
            }
        }
    }
</script>
```
:::

### 自定义内容

:::demo 可自定义标题，内容，组件大小，按钮文字和回调函数等
```html
<template>
    <bk-button
        :type="'primary'"
        @click="basicCustomDialog">自定义配置1</bk-button>

    <bk-dialog
        :is-show.sync="basicCustomSetting.isShow"
        :width="basicCustomSetting.width"
        :title="basicCustomSetting.title"
        :content="basicCustomSetting.content"
        :close-icon="basicCustomSetting.closeIcon"
        @confirm="done => done()">
    </bk-dialog>

    <bk-button
        :type="'primary'"
        @click="basicCustomDialog2">
        自定义配置2
    </bk-button>

    <bk-dialog
        :has-header="basicCustomSetting2.hasHeader"
        :is-show.sync="basicCustomSetting2.isShow"
        width="600"
        @confirm="confirmFn"
        @cancel="cancelFn">
        <div slot="content">
            <textarea class="bk-form-textarea" style="height: 180px;" placeholder="请输入内容..." v-model="textareaVal"></textarea>
        </div>
    </bk-dialog>
</template>

<script>
    export default {
        data () {
            return {
                basicCustomSetting: {
                    isShow: false,
                    width: 800,
                    title: '自定义配置对话框',
                    content: '简单配置的对话框',
                    closeIcon: false
                },
                basicCustomSetting2: {
                    hasHeader: false,
                    isShow: false
                }
            }
        },
        methods: {
            basicCustomDialog () {
                this.basicCustomSetting.isShow = true
            },
            basicCustomDialog2 () {
                this.basicCustomSetting2.isShow = true
            },
            confirmFn (done) {
                this.$bkInfo('你点击了确定')
                done()
            },
            cancelFn () {
                this.$bkInfo('你点了取消')
            }
        }
    }
</script>
```
:::

### 属性
| 参数 | 说明    | 类型      | 可选值       | 默认值   |
| ---- | ------ | --------- | ----------- | -------- |
| is-show | 是否显示弹窗，支持.sync修饰符 | Boolean | —— | false |
| width | 弹窗的宽度，支持数字和百分比 | Number/String | —— | —— |
| title | 弹窗的标题 | String | —— | —— |
| content | 弹窗的内容，当内容很简单，仅为字符串时可以使用 | String | —— | —— |
| has-header | 是否显示头部 | Boolean | —— | true |
| has-footer | 是否显示底部按钮 | Boolean | —— | true |
| ext-cls | 自定义的样式，传入的CSS类会被加在组件最外层的DOM上 | String | —— | —— |
| padding | 弹窗内容区的内边距 | Number/String | —— | —— |
| close-icon | 是否显示关闭按钮 | Boolean | —— | true |
| theme | 组件的主题色 | String | primary/info/warning/success/danger | primary |
| confirm | 确定按钮的文字 | String | —— | —— |
| cancel | 取消按钮的文字 | String | —— | —— |
| quick-cLose | 是否允许点击遮罩关闭弹窗 | Boolean | —— | true |

### 事件
| 事件名称 | 说明 | 回调参数 |
|---------|------|---------|
| confirm | 点击确定按钮的回调函数 | 关闭`Dialog`的函数，可手动调用执行关闭动作 |
| cancle | 点击确定按钮的回调函数 | —— |
