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
                    return ['string', 'number'].includes(typeof val) || val === null
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
            },
            autoCheck: {
                type: Boolean,
                default: true
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
                const originalValue = String(event.target.value).trim()
                const floatValue = originalValue.length ? Number(event.target.value.trim()) : null
                if (isNaN(floatValue)) {
                    value = this.autoCheck ? null : value
                } else {
                    value = floatValue
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
