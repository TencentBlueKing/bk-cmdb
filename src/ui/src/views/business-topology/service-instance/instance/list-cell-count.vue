<template>
    <span v-if="row.process_count">{{row.process_count}}</span>
    <span class="process-count-tips" v-else-if="row.service_template_id">
        <i class="tips-icon bk-icon icon-exclamation-circle"></i>
        <i18n class="tips-content" path="模板服务实例无进程提示">
            <cmdb-auth class="tips-link" place="link"
                :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
                @click.native.stop
                @click="redirectToTemplate">
                {{$t('跳转添加并同步')}}
            </cmdb-auth>
        </i18n>
    </span>
    <span class="process-count-tips" v-else>
        <i class="tips-icon bk-icon icon-exclamation-circle"></i>
        <i18n class="tips-content" path="普通服务实例无进程提示">
            <cmdb-auth class="tips-link" place="link"
                :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
                @click.native.stop
                @click="handleAddProcess">
                {{$t('立即添加')}}
            </cmdb-auth>
        </i18n>
    </span>
</template>

<script>
    import { mapGetters } from 'vuex'
    import createProcessMixin from './create-process-mixin'
    export default {
        name: 'list-cell-count',
        mixins: [createProcessMixin],
        props: {
            row: Object
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId'])
        },
        methods: {
            redirectToTemplate () {
                this.$routerActions.redirect({
                    name: 'operationalTemplate',
                    params: {
                        bizId: this.bizId,
                        templateId: this.row.service_template_id,
                        moduleId: this.row.bk_module_id
                    },
                    history: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .process-count-tips {
        display: flex;
        align-items: center;
        .tips-icon {
            color: $warningColor;
            font-size: 14px;
        }
        .tips-content {
            padding: 0 4px;
            color: $textDisabledColor;
            .tips-link {
                color: $primaryColor;
                cursor: pointer;
                &.disabled {
                    color: $textDisabledColor;
                }
            }
        }
    }
</style>
