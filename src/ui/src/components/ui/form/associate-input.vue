<template>
    <div class="form-associate-input">
        <input class="cmdb-form-input" type="text" v-model="localValue">
    </div>
</template>

<script>
    export default {
        name: 'cmdb-form-associate-input',
        props: ['value'],
        data () {
            return {
                localValue: ''
            }
        },
        watch: {
            value (value) {
                this.localValue = value
            },
            localValue (localValue) {
                this.$emit('input', localValue)
                this.$emit('on-change', localValue)
            }
        },
        created () {
            this.localValue = this.value
        }
    }
</script>