<template>
    <div class="cmdb-form form-singlechar">
        <input class="cmdb-form-input form-singlechar-input" type="text"
            :placeholder="$t('Form[\'请输入短字符\']')"
            :maxlength="maxlength"
            :value="value"
            :disabled="disabled"
            @input="handleInput($event)"
            @change="handleChange">
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-singlechar',
        props: {
            value: {
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            maxlength: {
                type: Number,
                default: 256
            }
        },
        methods: {
            handleInput (event) {
                let value = event.target.value.trim()
                this.$emit('input', value)
            },
            handleChange () {
                this.$emit('on-change', this.value)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-singlechar-input {
        height: 36px;
        width: 100%;
        padding: 0 10px;
        background-color: #fff;
        border: 1px solid $cmdbBorderColor;
        font-size: 14px;
        outline: none;
        &:focus{
            border-color: $cmdbBorderFocusColor;
        }
    }
</style>