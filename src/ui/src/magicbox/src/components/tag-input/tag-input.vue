<template>
	<div :class="['bk-tag-input', {'active': isEdit, 'disabled': disabled}]" @click="foucusInputer">
        <ul class="tag-list">
            <li class="key-node" v-for="(tag, index) in tagList">
                <span class="tag">{{tag}}</span>
                <a href="javascript:void(0)" class="remove-key" @click.stop.prevent="removeTag(tag, index)" v-if="!disabled">
                    <i class="bk-icon icon-close"></i>
                </a>
            </li>
            <li>
                <input
                    type="text"
                    class="input"
                    v-model="curInputValue"
                    v-if="!disabled"
                    :style="inputStyle"
                    @focus="focusHandler"
                    @blur="blurHandler"
                    @keyup="keyupHandler">
            </li>
        </ul>
        <p class="placeholder" v-show="!isEdit && !tagList.length">{{placeholder}}</p>
    </div>
</template>

<script> 
    export default {
        name: 'bk-tag-input',
        props: {
            placeholder: {
                type: String,
                default: '请输入并按Eeter结束'
            },
            tags: {
                type: Array,
                default () {
                    return []
                }
            },
            disabled: {
                type: Boolean,
                default: false
            },
            separator: {
                type: String,
                default: ''
            },
            isBlurTrigger: {
                type: Boolean,
                default: true
            }
        },
        data () {
            return {
                curInputValue: '',
                isCanRemoveTag: false,
                tagList: this.tags,
                isInputFocus: false,
                timer: 0,
                isEdit: false
            }
        },
        computed: {
            // 动态计算输入长度
            inputStyle () {
                let tag = this.curInputValue
                let charLen = this.getCharLength(tag) + 1
                return {width: charLen * 8 + 'px'}
            }
        },
        watch: {
            tags () {
                this.tagList = this.tags
            },
            curInputValue (newVal, oldVal) {
                if (newVal === '') {
                    this.isInputEmpty = true
                } else {
                    this.isInputEmpty = false
                }
            }
        },
        methods: {
            // 获取字符长度，汉字两个字节
            getCharLength (str) {
                let len = str.length
                let bitLen = 0
                for (let i = 0; i < len; i++) {
                    if ((str.charCodeAt(i) & 0xff00) !== 0) {
                        bitLen++
                    }
                    bitLen++
                }
                return bitLen
            },
            keyupHandler (event) {
                switch (event.code) {
                    case 'Enter':
                        this.addTag()
                        break
                    case 'Backspace':
                        this.backspaceHandler()
                        break
                    default:
                        this.isCanRemoveTag = false
                        break
                }
            },
            // 添加标签
            addTag () {
                if (this.separator) {
                    let tags = this.curInputValue.split(this.separator)
                    tags.forEach(tag => {
                        if (tag && !this.tagList.includes(tag)) {
                            this.tagList.push(tag)
                        }
                    })
                } else {
                    let tag = this.curInputValue
                    if (tag && !this.tagList.includes(tag)) {
                        this.tagList.push(tag)
                    }
                }
                
                this.$emit('change', this.tagList)
                this.$emit('update:tags', this.tagList)
                this.clearInput()
            },
            // 输入清除
            backspaceHandler () {
                // 如果清空输入
                console.log(this.isInputEmpty)
                if (!this.curInputValue) {
                    if (this.isCanRemoveTag) {
                        this.tagList.pop()
                        this.$emit('change', this.tagList)
                        this.$emit('update:tags', this.tagList)
                    }
                    this.isCanRemoveTag = true
                }
            },
            // 删除标签
            removeTag (data, index) {
                this.tagList.splice(index, 1)
                let tags = []
                this.$emit('change', this.tagList)
                this.$emit('update:tags', this.tagList)
            },
            // 清空输入框
            clearInput () {
                this.curInputValue = ''
                this.isCanRemoveTag = true
            },
            focusHandler () {
                this.isInputFocus = true
            },
            blurHandler () {
                if (this.isBlurTrigger) {
                    this.addTag()
                }
                this.timer = setTimeout(() => {
                    this.clearInput()
                    this.isEdit = false
                }, 300)
            },
            // 输入框获取焦点时触发
            foucusInputer () {
                if (this.disabled) return
                clearTimeout(this.timer)
                this.isEdit = true
                this.clearInput()
                this.$nextTick(() => {
                    this.$el.querySelector('.input').focus()
                })
            }
        }
    }
</script>

<style lang="scss">
    @import '../../bk-magic-ui/src/tag-input.scss'
</style>

