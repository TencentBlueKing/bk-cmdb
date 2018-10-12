<template>
    <div class="bk-member-selector">
		<tag-input 
            v-model="tags"
            :placeholder="placeholder"
            :disabled="isDisabled"
            :save-key="saveKey"
            :display-key="displayKey"
            :search-key="searchKey"
            :has-delete-icon="hasDeleteIcon"
            :maxData="maxData"
            :maxResult="maxResult"
            :list="renderList"
            :tpl="tpl"
            @change="change"
            @select="select"
            @remove="remove"></tag-input>
    </div>
</template>

<script>
    import tagInput from '../tag-input/tag-input.vue'
    import { uuid } from '../../util'
    
    export default {
        name: 'bk-member-selector',
        components: {
            tagInput
        },
        props: {
            placeholder: {
                type: String,
                default: '请输入'
            },
            disabled: {
                type: Boolean,
                default: false
            },
            hasDeleteIcon: {
                type: Boolean,
                default: false
            },
            type: {
                type: String,
                default: 'rtx'
            },
            value: {
                type: Array,
                default () {
                    return []
                }
            },
            maxData: {
                type: Number,
                default: -1
            },
            maxResult: {
                type: Number,
                default: 5
            }
        },
        data () {
            return {
                isDisabled: true,
                saveKey: this.type === 'email' ? 'Name' : 'english_name',
                displayKey: this.type === 'email' ? 'Name' : this.type === 'all' ? 'english_name' : 'english_name',
                searchKey: this.type === 'email' ? 'FullName' : 'english_name',
                renderList: [],
                tags: []
            }
        },
        watch: {
            tags (newVal) {
                this.$emit('input', newVal)
            }
        },
        mounted () {
            this.requestList()
            this.tags = [...this.value]
        },
        methods: {
            requestList () {
                let self = this
                let cookie = document.cookie.match(/bk_uid=\S*/)
                if (!cookie) {
                    let msg = '蓝鲸版人员选择器需要从Cookie 中获取 bk_uid 才能拉取人员信息\n本地开发时你需要：\n1）登录蓝鲸平台： open.oa.com \n2）本地配置host：例如：127.0.0.1 demo.open.oa.com\n3）通过域名（demo.open.oa.com）来访问本地服务'
                    alert(msg)
                    return false
                }

                let host = location.host
                let typeList = ['rtx', 'email', 'all']
                let prefix = host.indexOf('o.ied.com') > -1 ? 'http://o.ied.com/component/compapi/tof3/' : 'http://open.oa.com/component/compapi/tof3/'
                let config = {
                    url: '',
                    data: {}
                }

                if (!typeList.includes(this.type)) {
                    console.error('请配置正确的选择器类型')
                    return false
                }

                switch (this.type) {
                    case 'rtx':
                        config.url = `${prefix}get_all_staff_info/`
                        config.data = {
                            'query_type': 'simple_data',
                            'app_code': 'workbench'
                        }
                        break
                    case 'email':
                        config.url = `${prefix}get_all_ad_groups/`
                        config.data['query_type'] = undefined
                        config.data = {
                            'app_code': 'workbench'
                        }
                        break
                    case 'all':
                        config.url = `${prefix}get_all_rtx_and_mail_group/`
                        config.data = {
                            'app_code': 'workbench'
                        }
                        break
                    default:
                        break
                }

                this.ajaxRequest({
                    url: config.url,
                    jsonp: 'callback' + uuid(),
                    data: config.data,
                    success: function (res) {
                        if (res.result) {
                            res.data.map(val => {
                                self.renderList.push(val)
                            })
                            self.isDisabled = self.disabled
                        } else {
                            console.error(res.message)
                        }
                    },
                    error: function (error) {
                        console.error(error)
                    }
                })
            },
            ajaxRequest (params) {
                params = params || {}
                params.data = params.data || {}
                
                let callbackName = params.jsonp
                let head = document.getElementsByTagName('head')[0]
                params.data['callback'] = callbackName

                // 设置传递给后台的回调参数名
                let data = this.formatParams(params.data)
                let script = document.createElement('script')
                head.appendChild(script)

                // 创建jsonp回调函数
                window[callbackName] = function (res) {
                    head.removeChild(script)
                    clearTimeout(script.timer)
                    window[callbackName] = null
                    params.success && params.success(res)
                }

                // 发送请求
                script.src = params.url + '?' + data
            },
            // 格式化参数
            formatParams (data) {
                let arr = []
                for (let name in data) {
                    arr.push(encodeURIComponent(name) + '=' + encodeURIComponent(data[name]))
                }
                return arr.join('&')
            },
            tpl (node, ctx) {
                let parentClass = 'bk-selector-node bk-selector-member'
                let textClass = 'text'
                let imgClass = 'avatar'
                let template

                switch (this.type) {
                    case 'rtx':
                        template = (
                            <div class={parentClass}>
                                <img class={imgClass} src={`http://dayu.oa.com/avatars/${node.english_name}/avatar.jpg`} />
                                <span class={textClass}>{node.english_name}({node.chinese_name})</span>
                            </div>
                        )
                        break
                    case 'email':
                        template = (
                            <div class={parentClass}>
                                <span domPropsInnerHTML={node.FullName} class={textClass}></span>
                            </div>
                        )
                        break
                    case 'all':
                        template = (
                            <div class={parentClass}>
                                <span domPropsInnerHTML={node.english_name} class={textClass}></span><span>-</span>
                                <span domPropsInnerHTML={node.chinese_name} class={textClass}></span>
                            </div>
                        )
                        break
                    default:
                        break
                }
                
                return template
            },
            change (data) {
                this.$emit('change', data)
            },
            select () {
                this.$emit('select')
            },
            remove (data) {
                this.$emit('remove')
            }
        }
    }
</script>

<style lang="scss">
    @import '../../bk-magic-ui/src/member-selector.scss'
</style>
