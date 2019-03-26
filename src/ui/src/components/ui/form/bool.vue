<template>
    <div :class="['form-bool', {disabled}]" @click="handleChange">
        <input class="form-bool-input" type="checkbox"
            ref="input"
            :class="{
                indeterminate: indeterminate
            }"
            :style="style"
            :checked="checked"
            :disabled="disabled"
            :true-value="trueValue"
            :false-value="falseValue"
            v-model="localChecked"
            @click.prevent>
        <span class="form-bool-label" v-if="$slots.default">
            <slot></slot>
        </span>
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-bool',
        model: {
            prop: 'checked',
            event: 'change'
        },
        props: {
            checked: {
                default: false
            },
            disabled: {
                default: false
            },
            trueValue: {
                default: true
            },
            falseValue: {
                default: false
            },
            size: {
                type: Number,
                default: 0
            },
            indeterminate: Boolean
        },
        data () {
            return {
                localChecked: this.checked
            }
        },
        computed: {
            style () {
                let size = this.size ? this.size : 18
                return {
                    transform: `scale(${size / 18})`
                }
            }
        },
        watch: {
            checked (checked) {
                this.localChecked = checked
            },
            localChecked (localChecked) {
                this.$emit('change', localChecked, this)
                this.$emit('on-change', localChecked, this)
            }
        },
        created () {
            this.localChecked = this.checked
        },
        mounted () {
            for (let attr in this.$attrs) {
                this.$el.setAttribute(attr, '')
                this.$refs.input.setAttribute(attr, this.$attrs[attr])
            }
        },
        methods: {
            handleChange () {
                if (this.disabled) {
                    return false
                }
                this.localChecked = this.localChecked ? this.falseValue : this.trueValue
                this.$emit('click')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-bool{
        display: inline-block;
        vertical-align: middle;
        line-height: 1;
        white-space: nowrap;
        cursor: pointer;
        font-size: 0;
        &.disabled {
            cursor: not-allowed;
        }
    }
    .form-bool-input{
        display: inline-block;
        vertical-align: middle;
        width: 18px;
        height: 18px;
        cursor: pointer;
        outline: none;
        -webkit-appearance: none;
        background: #fff url("../../../assets/images/checkbox-sprite.png") no-repeat;
        background-position: 0 -62px;
        transform-origin: left center;
        &:checked{
            background-position: -33px -62px;
            &:disabled{
                background-position: -99px -62px;
            }
        }
        &:disabled{
            background-position: -66px -62px;
            cursor: not-allowed;
        }
        &.indeterminate {
            padding: 7px 3px;
            border: 1px solid $cmdbBorderFocusColor;
            border-radius: 2px;
            background-image: none;
            background-color: $cmdbBorderFocusColor;
            background-clip: content-box;
        }
    }
    .form-bool-label {
        display: inline-block;
        vertical-align: middle;
        font-size: 14px;
        @include ellipsis;
    }
</style>