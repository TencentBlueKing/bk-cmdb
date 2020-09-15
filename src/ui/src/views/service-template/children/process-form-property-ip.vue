<template>
    <bk-select style="width: 100%"
        v-model="localValue"
        v-bind="$attrs"
        :clearable="false">
        <bk-option id="1" name="127.0.0.1"></bk-option>
        <bk-option id="2" name="0.0.0.0"></bk-option>
        <bk-option id="3" :name="$t('第一内网IP')"></bk-option>
        <bk-option id="4" :name="$t('第一外网IP')"></bk-option>
    </bk-select>
</template>

<script>
    export default {
        props: {
            value: {
                type: String,
                default: ''
            }
        },
        computed: {
            localValue: {
                get () {
                    return this.value
                },
                set (value) {
                    this.$emit('input', value)
                    this.$emit('change', value)
                }
            }
        }
    }
</script>
