<template>
    <bk-select class="form-list-selector"
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
        ref="selector">
        <bk-option v-for="(option, index) in options"
            :key="index"
            :id="option"
            :name="option">
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        name: 'cmdb-form-list',
        props: {
            value: {
                type: [Array, String],
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
                default: () => []
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
        data () {
            return {
                selected: this.multiple ? [] : ''
            }
        },
        computed: {
            searchable () {
                return this.options.length > 7
            }
        },
        watch: {
            value (value) {
                this.selected = value
            },
            selected (selected) {
                this.$emit('input', selected)
                this.$emit('on-selected', selected)
            }
        },
        created () {
            this.initValue()
        },
        methods: {
            initValue () {
                try {
                    if (this.autoSelect && (!this.value || (this.multiple && !this.value.length))) {
                        this.selected = this.multiple ? [this.options[0]] : (this.options[0] || '')
                    } else {
                        this.selected = this.value
                    }
                } catch (error) {
                    this.selected = this.multiple ? [] : ''
                }
            },
            focus () {
                this.$refs.selector.show()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-list-selector {
        width: 100%;
    }
</style>
