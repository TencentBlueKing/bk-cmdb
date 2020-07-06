<template>
    <blueking-user-selector class="cmdb-form-objuser"
        ref="userSelector"
        display-list-tips
        v-bind="props"
        v-model="localValue">
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
                    return (this.value && this.value.length) ? this.value.split(',') : []
                },
                set (val) {
                    this.$emit('input', val.toString())
                    this.$emit('change', val.toString, this.value)
                }
            },
            props () {
                const props = { ...this.$attrs }
                if (this.api) {
                    props.api = this.api
                } else {
                    props.fuzzySearchMethod = this.fuzzySearchMethod
                    props.exactSearchMethod = this.exactSearchMethod
                    props.pasteValidator = this.pasteValidator
                }
                return props
            }
        },
        methods: {
            focus () {
                this.$refs.userSelector.focus()
            },
            async fuzzySearchMethod (keyword, page = 1) {
                const users = await this.$http.get(`${window.API_HOST}user/list`, {
                    params: {
                        fuzzy_lookups: keyword
                    },
                    config: {
                        cancelPrevious: true
                    }
                })
                return {
                    next: false,
                    results: users.map(user => ({
                        username: user.english_name,
                        display_name: user.chinese_name
                    }))
                }
            },
            exactSearchMethod (usernames) {
                const isBatch = Array.isArray(usernames)
                return Promise.resolve(isBatch ? usernames.map(username => ({ username })) : { username: usernames })
            },
            pasteValidator (usernames) {
                return Promise.resolve(usernames)
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-form-objuser {
        width: 100%;
    }
</style>
