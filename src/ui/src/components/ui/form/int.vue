<template>
    <bk-input type="text" ref="input"
        :placeholder="placeholder || $t('请输入数字')"
        :maxlength="maxlength"
        :disabled="disabled"
        v-model="localValue"
        @blur="handleInput"
        @change="handleChange">
        <template slot="append" v-if="unit">
            <div class="unit" :title="unit">{{unit}}</div>
        </template>
    </bk-input>
</template>

<script>
    export default {
        name: 'cmdb-form-int',
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
            maxlength: {
                type: Number,
                default: 11
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
        computed: {
            localValue: {
                get () {
                    return this.value === null ? '' : this.value
                },
                set (value) {
                    const emitValue = value === '' ? null : value
                    this.$emit('input', emitValue)
                    this.$emit('change', emitValue)
                    this.$emit('on-change', emitValue)
                }
            }
        },
        methods: {
            handleInput (value, event) {
                const originalValue = String(event.target.value).trim()
                const intValue = originalValue.length ? Number(event.target.value.trim()) : null
                if (isNaN(intValue)) {
                    value = this.autoCheck ? null : value
                } else {
                    value = intValue
                }
                this.localValue = value
                this.$refs.input.curValue = this.localValue
            },
            handleChange () {
                this.$emit('on-change', this.localValue)
            },
            focus () {
                this.$el.querySelector('input').focus()
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
