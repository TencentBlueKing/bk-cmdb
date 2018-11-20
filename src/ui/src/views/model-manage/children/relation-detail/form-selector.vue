<template>
    <bk-selector
        :disabled="disabled"
        :has-children="hasChildren"
        :list="list"
        :selected.sync="selected"
    ></bk-selector>
</template>

<script>
    export default {
        props: {
            value: {
                default: ''
            },
            hasChildren: {
                type: Boolean,
                default: false
            },
            disabled: {
                type: Boolean,
                default: false
            },
            list: {
                type: Array
            }
        },
        data () {
            return {
                selected: ''
            }
        },
        watch: {
            value (value) {
                this.selected = value
            },
            selected (selected) {
                this.$emit('input', selected)
                this.$emit('on-selected', selected)
            }
        },
        created () {
            this.selected = this.value
        }
    }
</script>

<style lang="scss" scoped>

</style>
