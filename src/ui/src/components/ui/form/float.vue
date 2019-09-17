<template>
    <bk-input type="text" ref="input"
        :placeholder="placeholder || $t('请输入浮点数')"
        :value="value"
        :disabled="disabled"
        @blur="handleInput"
        @change="handleChange">
        <template slot="append" v-if="unit">
            <div class="unit" :title="unit">{{unit}}</div>
        </template>
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
            },
            unit: {
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
                if (this.validateFloat(value)) {
                    this.localValue = parseFloat(value)
                } else {
                    this.$refs.input.curValue = null
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
    .unit {
        max-width: 120px;
        font-size: 12px;
        @include ellipsis;
        padding: 0 10px;
        height: 30px;
        line-height: 30px;
        background: #f2f4f8;
        color: #63656e;
    }
</style>
