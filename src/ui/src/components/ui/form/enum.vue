<template>
    <bk-select class="form-enum-selector"
        v-model="selected"
        :clearable="allowClear"
        :searchable="searchable"
        :disabled="disabled"
        :popover-options="{
            boundary: 'window'
        }">
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
                type: [String, Number],
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            allowClear: {
                type: Boolean,
                default: false
            },
            options: {
                type: Array,
                default () {
                    return []
                }
            }
        },
        data () {
            return {
                selected: ''
            }
        },
        computed: {
            searchable () {
                return this.options.length > 7
            }
        },
        watch: {
            value (value) {
                if (value !== null) {
                    this.selected = value
                }
            },
            selected (selected) {
                this.$emit('input', selected)
                this.$emit('on-selected', selected)
            },
            disabled (disabled) {
                this.setInitData()
            }
        },
        created () {
            this.setInitData()
        },
        methods: {
            setInitData () {
                if (this.value === '') {
                    const defaultOption = this.options.find(option => option['is_default'])
                    if (defaultOption) {
                        this.selected = defaultOption.id
                    } else {
                        this.$emit('input', null)
                    }
                } else {
                    this.selected = this.value
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-enum-selector{
        width: 100%;
    }
</style>
