<template>
    <div class="form-label">
        <span class="label-text">{{$t('ModelManagement["正则校验"]')}}</span>
        <textarea
            v-model="localValue"
            :disabled="isReadOnly"
            @input="handleInput"
        ></textarea>
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
