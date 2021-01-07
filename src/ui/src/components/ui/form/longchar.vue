<template>
    <bk-input
        v-model="localValue"
        :placeholder="placeholder || $t('请输入长字符')"
        :disabled="disabled"
        :type="'textarea'"
        :rows="row"
        :maxlength="maxlength"
        :clearable="!disabled"
        @enter="handleEnter"
        @on-change="handleChange">
    </bk-input>
</template>

<script>
    export default {
        name: 'cmdb-form-longchar',
        props: {
            value: {
                type: [String, Number],
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            maxlength: {
                type: Number,
                default: 2000
            },
            minlength: {
                type: Number,
                default: 2000
            },
            row: {
                type: Number,
                default: 3
            },
            placeholder: {
                type: String,
                default: ''
            }
        },
        computed: {
            localValue: {
                get () {
                    return (this.value === null || this.value === undefined) ? '' : this.value
                },
                set (value) {
                    this.$emit('input', value)
                }
            }
        },
        methods: {
            handleChange (value) {
                this.$emit('on-change', value)
            },
            handleEnter (value) {
                this.$emit('enter', value)
            },
            focus () {
                this.$el.querySelector('textarea').focus()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .bk-form-control {
        /deep/ .bk-textarea-wrapper {
            .bk-form-textarea {
                min-height: auto !important;
                padding: 5px 10px 8px;
                @include scrollbar-y;
                &.textarea-maxlength {
                    margin-bottom: 0 !important;
                }
            }
        }
        /deep/ .bk-limit-box {
            display: none !important
        }
    }
</style>
