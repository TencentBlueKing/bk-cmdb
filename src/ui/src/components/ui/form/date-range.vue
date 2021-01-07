<template>
    <div class="form-date-range">
        <bk-date-range class="bk-date-range"
            :placeholder="placeholder"
            :disabled="disabled"
            :position="position"
            :rangeSeparator="rangeSeparator"
            :quickSelect="quickSelect"
            :ranges="ranges"
            :timer="timer"
            :startDate="startDate"
            :endDate="endDate"
            @change="handleChange">
        </bk-date-range>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-date-range',
        props: {
            value: {
                type: Array,
                default () {
                    return []
                }
            },
            placeholder: {
                type: String,
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            position: {
                type: String,
                default: 'bottom-right',
                validator (val) {
                    return ['top', 'bottom', 'left', 'right', 'top-left', 'top-right', 'bottom-left', 'bottom-right'].includes(val)
                }
            },
            rangeSeparator: {
                type: String,
                default: ' - '
            },
            quickSelect: {
                type: Boolean,
                default: false
            },
            ranges: {
                type: Object
            },
            timer: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                localSelected: []
            }
        },
        computed: {
            startDate () {
                return this.localSelected.length ? this.localSelected[0] : ''
            },
            endDate () {
                return this.localSelected.length ? this.localSelected[1] : ''
            },
            separator () {
                return ` ${this.rangeSeparator} `
            }
        },
        watch: {
            value (value) {
                this.setLocalSelected()
            },
            localSelected (localSelected, oldSelected) {
                if (localSelected.join(this.separator) !== this.value.join(this.separator)) {
                    this.$emit('input', [...localSelected])
                    this.$emit('on-change', [...localSelected], [...oldSelected])
                }
            }
        },
        created () {
            this.setLocalSelected()
        },
        methods: {
            setLocalSelected () {
                if (this.localSelected.join(this.separator) !== this.value.join(this.separator)) {
                    this.localSelected = [...this.value]
                }
            },
            handleChange (oldValue, newValue) {
                this.localSelected = newValue.split(this.separator)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-date-range{
        position: relative;
        display: inline-block;
        vertical-align: middle;
    }
    .bk-date-range{
        width: 100%;
        white-space: normal;
    }
</style>