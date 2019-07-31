<template>
    <bk-input type="text"
        :placeholder="placeholder || $t('Form[\'请输入浮点数\']')"
        :value="value"
        :disabled="disabled"
        @blur="handleInput"
        @change="handleChange">
    </bk-input>
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
            },
            placeholder: {
                type: String,
                default: ''
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
            handleInput (value, event) {
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
