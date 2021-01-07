<template>
    <div class="instance-operation">
        <cmdb-auth tag="span" class="operation-item"
            v-if="!row.service_template_id"
            :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, relation: [bizId] }"
            @click.native.stop
            @click="handleAddProcess">
            {{$t('添加进程')}}
        </cmdb-auth>
        <cmdb-auth tag="span" class="operation-item"
            v-if="!row.service_template_id"
            :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, relation: [bizId] }"
            @click.native.stop
            @click="handleClone">
            {{$t('克隆')}}
        </cmdb-auth>
        <cmdb-auth tag="span" class="operation-item"
            :auth="{ type: $OPERATION.D_SERVICE_INSTANCE, relation: [bizId] }"
            @click.native.stop
            @click="handleDelete">
            {{$t('删除')}}
        </cmdb-auth>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import createProcessMixin from './create-process-mixin'
    import { MENU_BUSINESS_DELETE_SERVICE } from '@/dictionary/menu-symbol'
    export default {
        mixins: [createProcessMixin],
        props: {
            row: Object
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['selectedNode'])
        },
        methods: {
            handleClone () {
                this.$routerActions.redirect({
                    name: 'cloneServiceInstance',
                    params: {
                        instanceId: this.row.id,
                        hostId: this.row.bk_host_id,
                        setId: this.selectedNode.parent.data.bk_inst_id,
                        moduleId: this.selectedNode.data.bk_inst_id
                    },
                    query: {
                        title: this.row.name,
                        node: this.selectedNode.id
                    },
                    history: true
                })
            },
            handleDelete () {
                this.$routerActions.redirect({
                    name: MENU_BUSINESS_DELETE_SERVICE,
                    params: {
                        ids: this.row.id,
                        moduleId: this.selectedNode.data.bk_inst_id
                    },
                    history: true
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .instance-operation {
        display: flex;
        align-items: center;
    }
    .operation-item {
        display: inline-block;
        line-height: 32px;
        color: $textColor;
        font-size: 12px;
        cursor: pointer;
        color: $primaryColor;
        &:hover {
            opacity: .7;
        }
        &.disabled {
            color: $textDisabledColor;
            opacity: 1;
        }
        & ~ .operation-item {
            margin-left: 10px;
        }
    }
</style>
