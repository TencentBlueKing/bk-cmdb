<template>
    <bk-select class="form-timezone-selector"
        searchable
        v-model="selected"
        :clearable="false"
        :disabled="disabled"
        :multiple="multiple"
        :placeholder="placeholder"
        ref="selector">
        <bk-option
            v-for="(option, index) in timezoneList"
            :key="index"
            :id="option.id"
            :name="option.name">
        </bk-option>
    </bk-select>
</template>

<script>
    import TIMEZONE from './timezone.json'
    export default {
        name: 'cmdb-form-timezone',
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
            placeholder: {
                type: String,
                default: ''
            }
        },
        data () {
            const timezoneList = TIMEZONE.map(timezone => {
                return {
                    id: timezone,
                    name: timezone
                }
            })
            return {
                timezoneList,
                selected: this.multiple ? [] : ''
            }
        },
        watch: {
            value (value) {
                if (value !== this.selected) {
                    this.selected = value
                }
            },
            selected (selected) {
                this.$emit('input', selected)
                this.$emit('on-selected', selected)
            },
            disabled (disabled) {
                if (!disabled) {
                    this.selected = this.getDefaultValue()
                }
            }
        },
        created () {
            this.selected = this.getDefaultValue()
        },
        methods: {
            getDefaultValue () {
                let value = this.value || ''
                if (this.multiple && !value.length) {
                    value = ['Asia/Shanghai']
                } else {
                    value = value || 'Asia/Shanghai'
                }
                return value
            },
            focus () {
                this.$refs.selector.show()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-timezone-selector{
        width: 100%;
    }
</style>
