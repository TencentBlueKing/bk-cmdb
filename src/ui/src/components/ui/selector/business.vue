<template>
    <bk-select style="text-align: left;"
        v-model="localSelected"
        :searchable="true"
        :clearable="false"
        :placeholder="$t('请选择业务')"
        :disabled="disabled"
        :popover-options="popoverOptions">
        <bk-option
            v-for="(option, index) in authorizedBusiness"
            :key="index"
            :id="option.bk_biz_id"
            :name="option.bk_biz_name">
        </bk-option>
        <div class="business-extension" slot="extension" v-if="showApplyPermission || showApplyCreate">
            <a href="javascript:void(0)" class="extension-link"
                v-if="showApplyPermission"
                @click="handleApplyPermission">
                <i class="bk-icon icon-plus-circle"></i>
                {{$t('申请业务权限')}}
            </a>
            <a href="javascript:void(0)" class="extension-link"
                v-if="showApplyCreate"
                @click="handleApplyCreate">
                <i class="bk-icon icon-plus-circle"></i>
                {{$t('申请创建业务')}}
            </a>
        </div>
    </bk-select>
</template>

<script>
    import { translateAuth } from '@/setup/permission'
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
            },
            popoverOptions: {
                type: Object,
                default () {
                    return {}
                }
            },
            requestConfig: {
                type: Object,
                default () {
                    return {}
                }
            },
            showApplyPermission: Boolean,
            showApplyCreate: Boolean
        },
        data () {
            return {
                authorizedBusiness: [],
                localSelected: ''
            }
        },
        computed: {
            requireBusiness () {
                return this.$route.meta.requireBusiness
            }
        },
        watch: {
            localSelected (localSelected, old) {
                window.localStorage.setItem('selectedBusiness', localSelected)
                this.setHeader()
                this.$emit('input', localSelected)
                this.$emit('on-select', localSelected, old)
                this.$store.commit('objectBiz/setBizId', localSelected)
            },
            requireBusiness () {
                this.setHeader()
            }
        },
        async created () {
            this.authorizedBusiness = await this.$store.dispatch('objectBiz/getAuthorizedBusiness', this.requestConfig)
            if (this.authorizedBusiness.length) {
                this.init()
            } else {
                this.$emit('business-empty')
            }
        },
        methods: {
            init () {
                const selected = parseInt(window.localStorage.getItem('selectedBusiness'))
                const exist = this.authorizedBusiness.some(business => business.bk_biz_id === selected)
                if (exist) {
                    this.localSelected = selected
                } else if (this.authorizedBusiness.length) {
                    this.localSelected = this.authorizedBusiness[0]['bk_biz_id']
                }
            },
            setHeader () {
                if (this.requireBusiness) {
                    this.$http.setHeader('bk_biz_id', this.localSelected)
                } else {
                    this.$http.deleteHeader('bk_biz_id')
                }
            },
            async handleApplyPermission () {
                try {
                    const permission = []
                    const operation = this.$tools.getValue(this.$route.meta, 'auth.operation', {})
                    if (Object.keys(operation).length) {
                        const translated = await translateAuth(Object.values(operation))
                        permission.push(...translated)
                    }
                    const url = await this.$store.dispatch('auth/getSkipUrl', { params: permission })
                    window.open(url)
                } catch (e) {
                    console.error(e)
                }
            },
            handleApplyCreate () {}
        }
    }
</script>

<style lang="scss" scoped>
    .business-extension {
        width: calc(100% + 32px);
        margin-left: -16px;
    }
    .extension-link {
        display: block;
        line-height: 38px;
        background-color: #FAFBFD;
        padding: 0 9px;
        font-size: 13px;
        color: #63656E;
        &:hover {
            opacity: .85;
        }
        .bk-icon {
            font-size: 18px;
            color: #979BA5;
        }
    }
</style>
