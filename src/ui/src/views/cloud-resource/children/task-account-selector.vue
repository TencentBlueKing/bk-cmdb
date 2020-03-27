<template>
    <bk-select v-if="display === 'selector'"
        searchable
        :readonly="readonly"
        :disabled="disabled"
        :placeholder="$t('请选择xx', { name: $t('账户名称') })"
        :loading="$loading(requestId)"
        v-model="selected">
        <bk-option v-for="account in accounts"
            :key="account.bk_account_id"
            :name="account.bk_account_name"
            :id="account.bk_account_id">
        </bk-option>
    </bk-select>
    <span v-else>{{getAccountInfo()}}</span>
</template>

<script>
    export default {
        name: 'task-account-selector',
        props: {
            display: {
                type: String,
                default: 'selector'
            },
            readonly: Boolean,
            disabled: Boolean,
            value: {
                type: [String, Number]
            }
        },
        data () {
            return {
                accounts: [],
                requestId: 'taskAccountSelectorRequest'
            }
        },
        computed: {
            selected: {
                get () {
                    return this.value
                },
                set (value, oldValue) {
                    this.$emit('input', value)
                    this.$emit('change', value, oldValue)
                }
            }
        },
        created () {
            this.getAccounts()
        },
        methods: {
            async getAccounts () {
                try {
                    const { info: accounts } = await this.$store.dispatch('cloud/account/findMany', {
                        params: {},
                        config: {
                            requestId: this.requestId,
                            fromCache: true,
                            cacheExpire: 'page'
                        }
                    })
                    this.accounts = accounts
                } catch (e) {
                    console.error(e)
                    this.accounts = []
                }
            },
            getAccountInfo () {
                const account = this.accounts.find(account => account.bk_account_id === this.value)
                return account ? account.bk_account_name : '--'
            }
        }
    }
</script>
