<template>
    <bk-input type="text"
        v-model="localValue"
        :placeholder="localPlaceholder"
        :maxlength="maxlength"
        :disabled="disabled"
        v-bind="$attrs"
        @change="handleChange"
        @enter="handleEnter">
    </bk-input>
</template>

<script>
    export default {
        name: 'cmdb-form-singlechar',
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
                default: 256
            },
            placeholder: {
                type: String,
                default: ''
            }
        },
        computed: {
            localPlaceholder () {
                return this.placeholder || this.$t('请输入短字符')
            },
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
                this.$el.querySelector('input').focus()
            }
        }
    }
</script>
