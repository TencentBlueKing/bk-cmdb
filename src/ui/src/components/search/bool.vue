<template>
    <bk-select
        v-model="localValue"
        v-bind="$attrs">
        <bk-option id="true" name="true"></bk-option>
        <bk-option id="false" name="false"></bk-option>
    </bk-select>
</template>

<script>
    export default {
        name: 'cmdb-search-bool',
        props: {
            value: [String, Boolean]
        },
        computed: {
            localValue: {
                get () {
                    return this.value.toString()
                },
                set (value) {
                    let newValue = value
                    if (value.length) {
                        newValue = value === 'true'
                    }
                    this.$emit('input', newValue)
                    this.$emit('change', newValue)
                }
            }
        }
    }
</script>
