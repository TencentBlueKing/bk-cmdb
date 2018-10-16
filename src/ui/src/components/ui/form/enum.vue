<template>
    <div class="form-enum">
        <bk-selector class="form-enum-selector"
            :searchable="searchable"
            :list="options"
            :disabled="disabled"
            :allow-clear="allowClear"
            :selected.sync="selected">
        </bk-selector>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-enum',
        props: {
            value: {
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
                this.selected = value
            },
            selected (selected) {
                this.$emit('input', selected)
                this.$emit('on-selected', selected)
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
                    }
                } else {
                    this.selected = this.value
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-enum {
        .form-enum-selector{
            width: 100%;
        }
    }
</style>