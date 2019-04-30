<template>
    <div class="cmdb-form form-float">
        <input class="cmdb-form-input form-float-input" type="text"
            :placeholder="$t('Form[\'请输入浮点数\']')"
            :value="value"
            :disabled="disabled"
            @blur="handleInput"
            @change="handleChange">
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-float',
        props: {
            value: {
                default: null,
                validator (val) {
                    return typeof val === 'number' || val === '' || val === null
                }
            },
            disabled: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                localValue: null
            }
        },
        watch: {
            value (value) {
                this.localValue = this.value === '' ? null : this.value
            },
            localValue (localValue) {
                if (localValue !== this.value) {
                    this.$emit('input', localValue)
                }
            }
        },
        created () {
            this.localValue = this.value === '' ? null : this.value
        },
        methods: {
            handleInput (event) {
                if (this.validateFloat(event.target.value)) {
                    this.localValue = parseFloat(event.target.value)
                } else {
                    event.target.value = null
                    this.localValue = null
                }
            },
            handleChange () {
                this.$emit('on-change', this.localValue)
            },
            validateFloat (val) {
                return /^[+-]?([0-9]*[.]?[0-9]+|[0-9]+[.]?[0-9]*)([eE][+-]?[0-9]+)?$/.test(val)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-float-input {
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
