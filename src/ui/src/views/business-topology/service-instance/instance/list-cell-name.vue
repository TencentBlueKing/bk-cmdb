<template>
    <div :class="{ 'instance-name': true, disabled }">
        <span class="name-text" v-bk-overflow-tips>{{row.name}}</span>
        <cmdb-dot-menu class="instance-dot-menu" trigger="click" @click.native.stop>
            <ul class="menu-list">
                <cmdb-auth tag="li" class="menu-item"
                    v-if="!row.service_template_id"
                    :auth="{ type: $OPERATION.U_SERVICE_INSTANCE, bk_biz_id: bizId }"
                    @click="handleAddProcess">
                    {{$t('添加进程')}}
                </cmdb-auth>
                <cmdb-auth tag="li" class="menu-item"
                    v-if="!row.service_template_id"
                    :auth="{ type: $OPERATION.C_SERVICE_INSTANCE, bk_biz_id: bizId }"
                    @click="handleClone">
                    {{$t('克隆')}}
                </cmdb-auth>
                <cmdb-auth tag="li" class="menu-item"
                    :auth="{ type: $OPERATION.D_SERVICE_INSTANCE, bk_biz_id: bizId }"
                    @click="handleDelete">
                    {{$t('删除')}}
                </cmdb-auth>
            </ul>
        </cmdb-dot-menu>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    import { MENU_BUSINESS_DELETE_SERVICE } from '@/dictionary/menu-symbol'
    import createProcessMixin from './create-process-mixin'
    export default {
        name: 'list-cell-name',
        mixins: [createProcessMixin],
        props: {
            row: Object
        },
        computed: {
            ...mapGetters('objectBiz', ['bizId']),
            ...mapGetters('businessHost', ['selectedNode']),
            disabled () {
                return !this.row.process_count
            }
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
    .instance-name {
        display: flex;
        width: 100%;
        align-items: center;
        &.disabled {
            color: $textDisabledColor;
        }
        .name-text {
            max-width: calc(100% - 25px);
            font-weight: bold;
            @include ellipsis;
        }
        .instance-dot-menu {
            display: none;
        }
    }
    .menu-list {
        padding: 6px 0;
        font-size: 12px;
        .menu-item {
            padding: 0 12px;
            display: block;
            line-height: 32px;
            color: $textColor;
            font-size: 12px;
            cursor: pointer;
            &:hover {
                background-color: #E1ECFF;
                color: #3a84ff;
            }
            &.disabled {
                background-color: #fff;
                color: $textDisabledColor;
            }
        }
    }
</style>
