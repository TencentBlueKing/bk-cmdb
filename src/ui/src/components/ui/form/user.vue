<template>
    <blueking-user-selector class="cmdb-form-objuser"
        ref="userSelector"
        display-list-tips
        :api="api"
        v-model="localValue"
        v-bind="$attrs">
    </blueking-user-selector>
</template>

<script>
    import BluekingUserSelector from '@blueking/user-selector'
    export default {
        name: 'cmdb-form-objuser',
        components: {
            BluekingUserSelector
        },
        props: {
            value: {
                type: String,
                default: ''
            }
        },
        computed: {
            api () {
                return window.ESB.userManage
            },
            localValue: {
                get () {
                    return this.value.length ? this.value.split(',') : []
                },
                set (val) {
                    this.$emit('input', val.toString())
                    this.$emit('change', val.toString, this.value)
                }
            }
        },
        methods: {
            focus () {
                setTimeout(() => {
                    debugger
                    this.$refs.userSelector.focus()
                    this.$refs.userSelector.calcOverflow()
                }, 0)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-form-objuser {
        width: 100%;
    }
</style>
