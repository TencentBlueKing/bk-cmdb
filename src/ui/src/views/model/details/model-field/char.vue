<template>
    <div class="form-item">
        <label class="form-label">{{$t('Common["正则验证"]')}}</label>
        <input type="text" class="cmdb-form-input" name="" placeholder="" 
            v-model="localValue"
            :disabled="isReadOnly"
            @input="handleInput">
    </div>
</template>

<script>
    export default {
        props: {
            value: {
                default: ''
            },
            isReadOnly: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                localValue: ''
            }
        },
        watch: {
            value () {
                this.localValue = this.value === '' ? '' : this.value
            }
        },
        created () {
            this.localValue = this.value === '' ? '' : this.value
        },
        methods: {
            handleInput () {
                this.$emit('input', this.localValue)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-form-input {
        width: 625px !important;
    }
</style>
