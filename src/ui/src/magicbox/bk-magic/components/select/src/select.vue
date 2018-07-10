<template>
    <div class="bk-select"
        :class="{
            open: open,
            multiple: multiple
        }"
        @click="listToggle"
        v-clickoutside="close">
        <div class="bk-select-wrapper">
            <input type="text" class="bk-select-input" readonly="readonly"
                :placeholder="placeholder?placeholder:t('select.placeholder')"
                :disabled="disabled"
                v-model="model"
            >
            <i class="bk-icon icon-close bk-select-clear" v-if="showClear && isSelected" @click.stop="clear"></i>
            <i class="bk-icon icon-angle-down bk-select-icon" v-show="!showClear || (showClear && !isSelected)"></i>
        </div>
        <transition name="toggle-slide">
            <div class="bk-select-list"
                v-show="open">
                <slot name="pre-ext"></slot>
                <div class="bk-select-list-filter"
                    v-if="filterable"
                    @click.stop="open = true">
                    <input type="text" class="bk-select-filter-input" autofocus
                        v-model="filter">
                    <i class="bk-icon icon-search"></i>
                </div>
                <ul>
                    <template v-if="$slots.default">
                        <slot></slot>
                    </template>
                    <template v-if="!$slots.default">
                        <select-option
                            :value="'bk-no-value'"
                            :label="t('select.noData')"
                            :is-empty-mark="isEmptyMark"></select-option>
                    </template>
                </ul>
                <slot name="post-ext"></slot>
            </div>
        </transition>
    </div>
</template>

<script>
    import selectOption from './select-option'
    import clickoutside from './../../../utils/clickoutside'
    import {
        findValueInObj,
        isObject,
        isInArray,
        findValueInArrByRecord } from '../../../assets/js/utils'

    export default {
        name: 'bk-select',
        props: {
            placeholder: {
                type: String,
                default: ''
            },
            selected: {
                require: true,
                default: ''
            },
            valueKey: {
                type: String,
                default: 'value'
            },
            disabled: {
                type: Boolean,
                default: false
            },
            multiple: {
                type: Boolean,
                default: false
            },
            filterable: {
                type: Boolean,
                default: false
            },
            filterFn: {
                type: Function
            },
            list: {
                type: Array,
                default () {
                    return []
                }
            },
            showClear: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                open: false,
                curValue: this.multiple ? [] : '',
                curLabel: this.multiple ? [] : '',
                model: '',
                localOptions: [],         // 缓存所有待选项的value和label
                localOptionsLoaded: false,      // 记录当前子组件是否已全部挂载
                filter: '',
                create: true,
                isEmptyMark: true
            }
        },
        components: {
            'select-option': selectOption
        },
        directives: {
            clickoutside
        },
        watch: {
            filter (val) {
                let arr = []
                val = val.toLowerCase()
                this.localOptions.map((v) => {
                    if (v.isOption) {
                        let display = v.$el.style.display

                        if (v.label.toLowerCase().indexOf(val) === -1) {
                            v.$el.style.display = 'none'
                        } else {
                            v.$el.style.display = 'block'
                        }
                    } else {
                        let count = 0

                        v.$children.map((_v) => {
                            let _display = _v.$el.style.display

                            if (_v.label.toLowerCase().indexOf(val) === -1) {
                                _v.$el.style.display = 'none'
                                count++
                            } else {
                                _v.$el.style.display = 'block'
                                arr.push(_v)
                            }
                        })

                        if (count === v.$children.length) {
                            v.$el.style.display = 'none'
                        } else {
                            v.$el.style.display = 'block'
                        }
                    }
                })

                this.$emit('on-filter', this.filter, arr)
            },
            selected (val, oldVal) {
                this.$nextTick(() => {
                    let value

                    if (this.multiple) {
                        value = val.length ? val : ''
                    } else {
                        value = val
                    }

                    this.setSelected(value)
                })
            },
            localOptions (localOptions) {
                this.$nextTick(() => {
                    this.setSelected(this.selected)
                })
            }
        },
        computed: {
            isSelected () {
                return ![null, undefined, ''].includes(this.selected)
            }
        },
        methods: {
            listToggle () {
                if (!this.disabled) {
                    this.open = !this.open
                    this.$emit('on-toggle', this.open)
                }
            },
            close () {
                this.open = false
                this.$emit('on-toggle', false)
                if (this.filterable) {
                    setTimeout(() => {
                        this.filter = ''
                    }, 200)
                }
            },
            /**
             *  挂载子组件后缓存子组件
             *  @param item 当前挂载的子组件的对象
             */
            addOption (item, isGroup) {
                let loaded = this.localForceUpdate
                let preLength = 0

                this.localOptions.push(item)
                this.$nextTick(() => {
                    if (this.list.length === this.localOptions.length) {
                        this.selected && this.setSelected(this.selected)
                    }
                })

                if (isGroup) {
                    this.localOptions.map((v) => {
                        preLength += v.$children.length
                    })
                }

                // 返回localOptions的长度 - 1给当前子组件作为index值
                return isGroup ? preLength : this.localOptions.length - 1
            },
            removeOption (el) {
                let index = this.localOptions.indexOf(el)

                index > -1 && this.localOptions.splice(index, 1)
            },
            updateOption (optIndex, opt) {
                this.model = opt.localData.label
                this.curLabel = opt.localData.label
                this.localOptions.splice(optIndex, 1, opt)
            },
            // 点击选项后的handler
            optionClickHandlder (child, index) {
                let multipleArr = []
                let curValueStr

                this.setSelected(child)

                if (!this.multiple) {
                    this.close()
                } else {
                    let curLabel = this.curLabel
                    let curValue = this.curValue

                    if (curLabel.length === curValue.length) {
                        for (let i = 0, len = curLabel.length; i < len; i++) {
                            multipleArr.push({
                                label: curLabel[i],
                                value: curValue[i]
                            })
                        }
                    }

                    curValueStr = this.curValue.join(',')
                }

                this.$emit('update:selected', this.multiple ? curValueStr : this.curValue)
                this.$emit('on-selected', this.multiple ? curValueStr : child, index, this.multiple ? multipleArr : undefined)
            },
            // 根据传入的值/对象获取指定的option数据
            getOption (args) {
                let _isObject = isObject(args)
                let {
                    valueKey
                } = this
                let $opt
                let isBreak
                
                function testEqual (args, opt) {
                    let isEqual

                    if (args instanceof Array) {        // 传入的数据是数组
                        for (let _args of args) {
                            isEqual = _isObject ? findValueInObj((isObject(opt.value) ? opt.value : opt), valueKey) === findValueInObj(_args, valueKey) : opt.value.toString() === _args.toString()
                        }
                    } else {
                        isEqual = _isObject ? findValueInObj((isObject(opt.value) ? opt.value : opt), valueKey) === findValueInObj(args, valueKey) : opt.value.toString() === args.toString()
                    }

                    return isEqual
                }

                // 遍历当前缓存的选项内容，找到option
                for (let opt of this.localOptions) {
                    let realLocalOptions

                    if (opt.value === undefined) { // 当前子组件是option-group
                        realLocalOptions = opt.localOptions

                        for (let _opt of realLocalOptions) {
                            if (testEqual(args, _opt)) {
                                $opt = _opt
                                isBreak = true
                                break
                            }
                        }
                    } else {
                        realLocalOptions = opt

                        if (testEqual(args, opt)) {
                            $opt = opt
                            isBreak = true
                            break
                        }
                    }

                    if (isBreak) break
                }

                if ($opt) return $opt

                // 如果之前未匹配到，则返回一个新的对象
                return {
                    value: args,
                    label: !isObject ? args : ''
                }
            },
            _isEmpty (val) {
                return typeof val === 'number' ? val : (val || '')
            },
            // 设置选中的项
            setSelected (val, index) {
                let target = this.getOption(Number(val) || val)
                let $value = target.value
                let {
                    valueKey
                } = this

                if (!this.multiple) {
                    if (isObject($value)) {
                        this.curValue = $value[valueKey]
                        this.curLabel = $value.label
                    } else {
                        this.curValue = this._isEmpty($value)
                        this.curLabel = this._isEmpty(target.label)
                    }
                } else {
                    if (typeof val === 'string') {  // 传入的是字符串，用于初始化赋值
                        this.curValue.splice(0, this.curValue.length)
                        this.curLabel.splice(0, this.curLabel.length)

                        if (val.length) {
                            let value = val.split(',')

                            for (let item of value) {
                                let $item

                                if (isObject(item)) {
                                    $item = item[valueKey]
                                } else {
                                    $item = item
                                }

                                this.curValue.push($item)
                                this.curLabel.push(findValueInArrByRecord(this.localOptions, { value: $item }, 'label'))
                            }
                        }
                    } else if (!(val instanceof Array)) {   // 传入的是对象，用于点击选择
                        let isInArr = isInArray(this.curValue, $value)
                        if (isInArr.result) {      // 已经在数组中
                            this.curValue.splice(isInArr.index, 1)
                            this.curLabel.splice(isInArr.index, 1)
                        } else {                   // 不在数组中
                            this.curValue.push($value)
                            this.curLabel.push(target.label)
                        }
                    }
                }

                this.model = this.curLabel
            },
            clear () {
                this.$emit('update:selected', '')
                this.$emit('on-selected', '', undefined)
            }
        },
        mounted () {
            this.selected.length ? this.setSelected(this.selected) : ''
        },
        beforeDestroy () {
            this.open = false,
            this.curValue = this.multiple ? [] : '',
            this.curLabel = this.multiple ? [] : '',
            this.model = '',
            this.localOptions = [],
            this.localOptionsLoaded = false,
            this.filter = '',
            this.create = true,
            this.isEmptyMark = true
        }
    }
</script>
