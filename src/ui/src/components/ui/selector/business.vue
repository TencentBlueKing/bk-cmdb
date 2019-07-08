<template>
    <bk-select style="text-align: left;"
        v-model="localSelected"
        :searchable="authorizedBusiness.length > 5"
        :clearable="false"
        :placeholder="$t('请选择业务')"
        :disabled="disabled">
        <bk-option
            v-for="(option, index) in authorizedBusiness"
            :key="index"
            :id="option.bk_biz_id"
            :name="option.bk_biz_name">
        </bk-option>
    </bk-select>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        name: 'cmdb-business-selector',
        props: {
            value: {
                type: [String, Number],
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            }
        },
        data () {
            return {
                authorizedBusiness: [],
                localSelected: ''
            }
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            requireBusiness () {
                return this.$route.meta.requireBusiness
            }
        },
        watch: {
            localSelected (localSelected, prevSelected) {
                window.localStorage.setItem('selectedBusiness', localSelected)
                if (prevSelected !== '') {
                    window.location.reload()
                    return
                }
                this.setHeader()
                this.$emit('input', localSelected)
                this.$emit('on-select', localSelected)
                this.setLocalSelected()
            },
            value (value) {
                if (value !== this.localSelected) {
                    this.setLocalSelected()
                }
            },
            bizId (value) {
                this.localSelected = value
            },
            requireBusiness () {
                this.setHeader()
            }
        },
        async created () {
            this.authorizedBusiness = await this.$store.dispatch('objectBiz/getAuthorizedBusiness')
            if (this.authorizedBusiness.length) {
                this.setLocalSelected()
            }
        },
        methods: {
            setHeader () {
                if (this.requireBusiness) {
                    this.$http.setHeader('bk_biz_id', this.localSelected)
                } else {
                    this.$http.deleteHeader('bk_biz_id')
                }
            },
            setLocalSelected () {
                const selected = this.value || parseInt(window.localStorage.getItem('selectedBusiness'))
                const exist = this.authorizedBusiness.some(business => business['bk_biz_id'] === selected)
                if (exist) {
                    this.localSelected = selected
                } else if (this.authorizedBusiness.length) {
                    this.localSelected = this.authorizedBusiness[0]['bk_biz_id']
                }
                this.$store.commit('objectBiz/setBizId', this.localSelected)
            }
        }
    }
</script>
