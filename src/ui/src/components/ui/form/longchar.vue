<template>
    <div class="cmdb-form form-longchar">
        <input class="cmdb-form-input form-longchar-input" type="text"
            :placeholder="$t('Form[\'请输入长字符\']')"
            :maxlength="maxlength"
            :value="value"
            :disabled="disabled"
            @input="handleInput($event)"
            @change="handleChange">
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-longchar',
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
                default: 2000
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
    .form-longchar-input {
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