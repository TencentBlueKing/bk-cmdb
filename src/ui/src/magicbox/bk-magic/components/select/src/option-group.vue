<template>
    <li class="bk-select-group">
        <ul>
            <li class="bk-select-group-name">{{ label }}</li>
            <slot></slot>
        </ul>
    </li>
</template>

<script>
    export default {
        name: 'bk-option-group',
        props: {
            label: {
                type: [String, Number],
                default: ''
            }
        },
        data () {
            return {
                isOption: false,
                localOptions: [],
                localOptionsLoaded: false,
                multiple: this.$parent.multiple,
                render: this.$parent.render,
                preLength: -1 // 当前group组件在父组件中的index
            }
        },
        computed: {
            curValue () {
                return this.$parent.curValue
            },
            curLabel () {
                return this.$parent.curLabel
            },
            create () {
                return this.$parent.create
            }
        },
        methods: {
            addOption (item) {
                let loaded = this.localOptionsLoaded

                if (!loaded) {
                    this.localOptions.push(item)

                    // 若当前localOptions的长度等于挂载的default子组件的个数，则代表子组件已挂载完成
                    if (this.localOptions.length === this.$slots.default.length) {
                        this.localOptionsLoaded = true
                    }
                } else {
                    // 清除当前localOptions供新子组件挂载
                    this.localOptions.splice(0, this.localOptions.length)
                    this.localOptions.push(item)
                    this.localOptionsLoaded = false
                }

                // 返回localOptions的长度 - 1给当前子组件作为index值
                return this.preLength + this.localOptions.length - 1
            },
            updateOption (optIndex, opt) {
                this.$parent.updateOption(optIndex, opt)
            },
            optionClickHandlder (v, k) {
                this.$parent.optionClickHandlder(v, k)
            },
            removeOption (el) {
                this.$parent.removeOption(el)
            }
        },
        created () {
            this.preLength = this.$parent.addOption(this, true)
        }
    }
</script>
