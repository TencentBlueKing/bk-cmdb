<template>
    <blueking-user-selector class="cmdb-form-objuser"
        ref="userSelector"
        display-list-tips
        v-bind="props"
        v-model="localValue"
        :class="{
            'has-fast-select': fastSelect
        }"
        @focus="$emit('focus')"
        @blur="$emit('blur')">
    </blueking-user-selector>
</template>

<script>
    import BluekingUserSelector from '@blueking/user-selector'
    import { mapGetters } from 'vuex'
    import Vue from 'vue'
    export default {
        name: 'cmdb-form-objuser',
        components: {
            BluekingUserSelector
        },
        props: {
            value: {
                type: String,
                default: ''
            },
            fastSelect: Boolean
        },
        computed: {
            ...mapGetters(['userName']),
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
        mounted () {
            this.setupFastSelect()
        },
        methods: {
            setupFastSelect () {
                if (!this.fastSelect) return
                const FastSelect = new Vue({
                    render: (h) => {
                        return (
                            <span class="fast-select"
                                on-click={ this.handleFastSelect }>
                                { this.$i18n.locale === 'en' ? 'me' : '我' }
                            </span>
                        )
                    }
                })
                FastSelect.$mount()
                FastSelect.$el.setAttribute([this.$options._scopeId], true)
                const container = this.$refs.userSelector.$refs.container
                container.parentElement.append(FastSelect.$el)
            },
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
            },
            handleFastSelect (event) {
                event.stopPropagation()
                const value = [...this.localValue]
                const exist = value.includes(this.userName)
                if (exist) return
                if (this.$refs.userSelector.multiple) {
                    value.push(this.userName)
                } else {
                    value.splice(0, value.length, this.userName)
                }
                this.localValue = value
            }
        }
    }
</script>

<style lang="scss" scoped>
    .cmdb-form-objuser {
        width: 100%;
        &.has-fast-select {
            /deep/ .user-selector-container {
                padding-right: 20px;
            }
        }
    }
    .fast-select {
        position: absolute;
        top: 50%;
        right: 4px;
        font-size: 12px;
        line-height: 16px;
        margin-top: -8px;
        color: $textColor;
        z-index: 2;
        cursor: pointer;
    }
</style>
