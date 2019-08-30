<template>
    <bk-input type="text"
        :placeholder="placeholder || $t('请输入长字符')"
        :maxlength="maxlength"
        :value="value"
        :disabled="disabled"
        @input="handleInput"
        @change="handleChange">
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
            placeholder: {
                type: String,
                default: ''
            }
        },
        methods: {
            handleInput (value) {
                this.$emit('input', value.trim())
            },
            handleChange () {
                this.$emit('on-change', this.value)
            }
        }
    }
</script>
