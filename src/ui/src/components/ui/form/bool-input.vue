<template>
    <bk-input class="cmdb-form-input" type="text" v-model.trim="localValue"></bk-input>
</template>

<script>
    export default {
        name: 'cmdb-form-bool-input',
        // eslint-disable-next-line
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
                if (['true', 'false'].includes(localValue)) {
                    localValue = [true, false].find(value => value.toString() === localValue)
                }
                this.$emit('input', localValue)
                this.$emit('on-change', localValue)
            }
        },
        created () {
            this.localValue = this.value
        }
    }
</script>
