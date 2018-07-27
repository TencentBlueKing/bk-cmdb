<template id="">
    <li class="bk-select-list-item"
        :title="localLabel"
        :class="{
            selected: selected,
            disabled: disabled}"
        @click.stop="optionClick($event)">
        <span class="bk-form-checkbox bk-checkbox-small"
            v-if="$parent.multiple && !isEmptyMark">
            <input type="checkbox"
                :value="value"
                :disabled="disabled"
                v-model="selected">
            {{ localLabel }}
        </span>
        <span v-else>{{ localLabel }}</span>
    </li>
</template>

<script>
    import { isObject, isInArray } from '../../../assets/js/utils'

    export default {
        name: 'bk-select-option',
        props: {
            value: {
                required: true
            },
            label: {
                type: [String, Number]
            },
            disabled: {
                type: Boolean,
                default: false
            },
            isEmptyMark: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                localData: {
                    value: this.value,
                    label: this.label
                },
                // 当前选项在列表中的index值
                optIndex: -1,
                isOption: true
            }
        },
        computed: {
            localLabel () {
                return this.label || this.value
            },
            // 当前项是否被选中
            selected: {
                get () {
                    if (this.$parent.multiple) {
                        return isInArray(this.$parent.curValue, this.value).result
                    } else {
                        return this.value === this.$parent.curValue
                    }
                },
                set () {}
            }
        },
        watch: {
            value (value) {
                this.$parent.updateOption(this.optIndex, this)
            },
            label (label) {
                this.localData.label = label
                this.$parent.updateOption(this.optIndex, this)
            }
        },
        methods: {
            optionClick (e) {
                if (this.disabled || this.isEmptyMark) return

                this.$parent.optionClickHandlder(isObject(this.value) ? this.value : this.localData, this.optIndex)
            }
        },
        created () {
            if (this.isEmptyMark) return
            // 将当前节点传入bk-select组件
            this.$parent.create ? this.optIndex = this.$parent.addOption(this, false) : ''
        },
        beforeDestroy () {
            this.$parent.removeOption(this)
        }
    }
</script>
