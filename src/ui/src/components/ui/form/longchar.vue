<template>
    <search-input :value="value"
        :disabled="disabled"
        :maxlength="maxlength"
        :placeholder="placeholder || $t('请输入长字符')"
        @input="handleInput"
        @enter="handleEnter"
        @on-change="handleChange">
    </search-input>
</template>

<script>
    import searchInput from './search-input.vue'
    export default {
        name: 'cmdb-form-longchar',
        components: {
            searchInput
        },
        props: {
            value: {
                type: [String, Number],
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            },
            maxlength: {
                type: Number,
                default: 2000
            },
            placeholder: {
                type: String,
                default: ''
            }
        },
        methods: {
            handleInput (value) {
                this.$emit('input', value)
            },
            handleChange (value) {
                this.$emit('on-change', value)
            },
            handleEnter (value) {
                this.$emit('enter', value)
            }
        }
    }
</script>
