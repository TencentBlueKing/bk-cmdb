<template>
    <div class="form-bool" @click="handleClick">
        <input class="form-bool-input" type="checkbox"
            ref="input"
            :checked="checked"
            :disabled="disabled"
            :true-value="trueValue"
            :false-value="falseValue"
            @click.stop
            @change="handleChange($event)">
        <slot></slot>
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
            }
        },
        mounted () {
            for (let attr in this.$attrs) {
                this.$el.setAttribute(attr, '')
                this.$refs.input.setAttribute(attr, this.$attrs[attr])
            }
        },
        methods: {
            handleClick () {
                const value = this.$refs.input.checked ? this.falseValue : this.trueValue
                this.$emit('change', value)
                this.$emit('on-change', value)
            },
            handleChange (event) {
                const value = event.target.checked ? this.trueValue : this.falseValue
                this.$emit('change', value)
                this.$emit('on-change', value)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .form-bool{
        display: inline-block;
        vertical-align: middle;
        line-height: 1;
        cursor: pointer;
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
    }
</style>