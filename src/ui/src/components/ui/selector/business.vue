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
    import { mapGetters } from 'vuex'
    export default {
        name: 'cmdb-business-selector',
        props: {
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
            showApplyPermission: Boolean,
            showApplyCreate: Boolean
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId', 'authorizedBusiness']),
            localSelected: {
                get () {
                    return this.bizId
                },
                set (value, oldValue) {
                    this.$emit('select', value, oldValue)
                }
            }
        },
        methods: {
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
            vertical-align: text-top;
        }
    }
</style>
