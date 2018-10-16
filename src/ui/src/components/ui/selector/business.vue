<template>
    <bk-selector
        :list="privilegeBusiness"
        :selected.sync="localSelected"
        :searchable="privilegeBusiness.length > 5"
        :disabled="disabled"
        setting-key="bk_biz_id"
        display-key="bk_biz_name">
    </bk-selector>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        name: 'cmdb-business-selector',
        props: {
            value: {
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                localSelected: ''
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['privilegeBusiness'])
        },
        watch: {
            localSelected (localSelected) {
                window.localStorage.setItem('selectedBusiness', localSelected)
                if (this.$route.meta.requireBusiness) {
                    this.$http.setHeader('bk_biz_id', localSelected)
                } else {
                    this.$http.deleteHeader('bk_biz_id')
                }
                this.$emit('input', localSelected)
                this.$emit('on-select', localSelected)
            },
            value (value) {
                if (value !== this.localSelected) {
                    this.setLocalSelected()
                }
            }
        },
        beforeCreate () {
            this.$http.deleteHeader('bk_biz_id')
        },
        async created () {
            await this.getPrivilegeBusiness()
            if (this.privilegeBusiness.length) {
                this.setLocalSelected()
            } else {
                this.$error(this.$t('Common["您没有业务权限"]'))
            }
        },
        beforeDestroy () {
            this.$http.deleteHeader('bk_biz_id')
        },
        methods: {
            getPrivilegeBusiness () {
                return this.$store.dispatch('objectBiz/searchBusiness', {
                    config: {
                        requestId: 'post_searchBusiness',
                        fromCache: true
                    }
                })
            },
            setLocalSelected () {
                const selected = this.value || parseInt(window.localStorage.getItem('selectedBusiness'))
                const exist = this.privilegeBusiness.some(business => business['bk_biz_id'] === selected)
                if (exist) {
                    this.localSelected = selected
                } else if (this.privilegeBusiness.length) {
                    this.localSelected = this.privilegeBusiness[0]['bk_biz_id']
                }
            }
        }
    }
</script>