<template>
    <bk-selector class="form-selector"
        :placeholder="placeholder"
        :searchable="searchable"
        :list="list"
        :disabled="disabled"
        :allow-clear="allowClear"
        :selected.sync="selected"
        :has-children="hasChildren"
        :setting-key="settingKey"
        :search-key="displayKey"
        :display-key="displayKey">
    </bk-selector>
</template>

<script>
    export default {
        name: 'cmdb-selector',
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
            list: {
                type: Array,
                default () {
                    return []
                }
            },
            settingKey: {
                type: String,
                default: 'id'
            },
            displayKey: {
                type: String,
                default: 'name'
            },
            autoSelect: {
                type: Boolean,
                default: true
            },
            placeholder: {
                type: String,
                default: ''
            },
            hasChildren: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                selected: ''
            }
        },
        computed: {
            searchable () {
                if (this.hasChildren) {
                    let list = []
                    this.list.forEach(group => {
                        list = [...list, ...group.children]
                    })
                    return list.length > 7
                } else {
                    return this.list.length > 7
                }
            },
            selectedOption () {
                return this.list.find(option => option[this.settingKey] === this.selected)
            }
        },
        watch: {
            value (value) {
                this.selected = value
            },
            selected (selected) {
                this.$emit('input', selected)
                this.$emit('on-selected', selected, this.selectedOption)
            },
            list () {
                this.setInitData()
            }
        },
        created () {
            this.setInitData()
        },
        methods: {
            setInitData () {
                let value = this.value
                if (this.autoSelect) {
                    const currentOption = this.list.find(option => option[this.settingKey] === this.value)
                    if (!currentOption) {
                        value = this.list.length ? this.list[0][this.settingKey] : this.value
                    }
                }
                this.selected = value
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-selector{
        display: inline-block;
        vertical-align: middle;
        width: 100%;
    }
</style>