<template>
    <bk-select class="form-enum-selector"
        v-model="selected"
        :clearable="allowClear"
        :searchable="searchable"
        :disabled="disabled"
        :multiple="multiple"
        :placeholder="placeholder"
        :font-size="fontSize"
        :popover-options="{
            boundary: 'window'
        }"
        v-bind="$attrs"
        ref="selector">
        <bk-option
            v-for="(option, index) in options"
            :key="index"
            :id="option.id"
            :name="option.name">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        name: 'cmdb-form-enum',
        props: {
            value: {
                type: [Array, String, Number],
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            multiple: {
                type: Boolean,
                default: false
            },
            allowClear: {
                type: Boolean,
                default: false
            },
            autoSelect: {
                type: Boolean,
                default: true
            },
            options: {
                type: Array,
                default () {
                    return []
                }
            },
            placeholder: {
                type: String,
                default: ''
            },
            fontSize: {
                type: [String, Number],
                default: 'medium'
            }
        },
        computed: {
            searchable () {
                return this.options.length > 7
            },
            selected: {
                get () {
                    if (this.isEmpty(this.value)) {
                        return this.getDefaultValue()
                    }
                    return this.value
                },
                set (value) {
                    let emitValue = value
                    if (value === '') {
                        emitValue = this.multiple ? [] : null
                    }
                    this.$emit('input', emitValue)
                    this.$emit('on-selected', emitValue)
                }
            }
        },
        watch: {
            value: {
                immediate: true,
                handler (value) {
                    this.checkSelected()
                }
            }
        },
        methods: {
            isEmpty (value) {
                return ['', undefined, null].includes(value)
            },
            getDefaultValue () {
                if (this.autoSelect) {
                    const defaultOption = this.options.find(option => option['is_default'])
                    return defaultOption
                        ? this.multiple ? [defaultOption.id] : defaultOption.id
                        : ''
                }
                return this.multiple ? [] : ''
            },
            checkSelected () {
                const selected = this.selected
                if (this.value !== selected) {
                    this.selected = selected
                }
            },
            focus () {
                this.$refs.selector.show()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-enum-selector{
        width: 100%;
    }
</style>
