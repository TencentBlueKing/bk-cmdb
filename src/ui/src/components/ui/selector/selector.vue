<template>
    <bk-select v-if="hasChildren"
        v-model="selected"
        :placeholder="placeholder"
        :searchable="searchable"
        :clearable="allowClear"
        :disabled="disabled"
        :loading="loading"
        :font-size="fontSize"
        :popover-options="popoverOptions"
        :readonly="readonly">
        <bk-option-group v-for="(group, index) in list"
            :key="index"
            :name="group[displayKey]">
            <bk-option v-for="option in group.children || []"
                :key="option[settingKey]"
                :id="option[settingKey]"
                :name="option[displayKey]">
                <slot v-bind="option" />
            </bk-option>
        </bk-option-group>
    </bk-select>
    <bk-select v-else
        v-model="selected"
        :placeholder="placeholder"
        :searchable="searchable"
        :clearable="allowClear"
        :disabled="disabled"
        :loading="loading"
        :font-size="fontSize"
        :popover-options="popoverOptions"
        :readonly="readonly">
        <bk-option
            v-for="option in list"
            :key="option[settingKey]"
            :id="option[settingKey]"
            :name="option[displayKey]">
            <slot v-bind="option" />
        </bk-option>
    </bk-select>
</template>

<script>
    export default {
        name: 'cmdb-selector',
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
            },
            emptyText: {
                type: String,
                default: ''
            },
            fontSize: {
                type: String,
                default: 'medium'
            },
            searchable: {
                type: Boolean,
                default: false
            },
            loading: Boolean,
            popoverOptions: {
                type: Object,
                default: () => ({})
            },
            readonly: Boolean
        },
        data () {
            return {
                selected: ''
            }
        },
        computed: {
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
