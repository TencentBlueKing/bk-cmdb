<template>
    <bk-input type="text" ref="input"
        :placeholder="placeholder || $t('请输入数字')"
        :value="value"
        :maxlength="maxlength"
        :disabled="disabled"
        @blur="handleInput"
        @change="handleChange">
    </bk-input>
</template>

<script>
    export default {
        name: 'cmdb-form-int',
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
            maxlength: {
                type: Number,
                default: 11
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
                value = parseInt(event.target.value.trim())
                if (isNaN(value)) {
                    value = null
                }
                this.$refs.input.curValue = value
                this.localValue = value
            },
            handleChange () {
                this.$emit('on-change', this.localValue)
            }
        }
    }
</script>
