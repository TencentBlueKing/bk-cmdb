<template>
    <bk-input type="text"
        :placeholder="localPlaceholder"
        :maxlength="maxlength"
        :value="value"
        :disabled="disabled"
        @input="handleInput"
        @change="handleChange">
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
