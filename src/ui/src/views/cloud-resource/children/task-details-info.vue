<template>
    <cmdb-sticky-layout>
        <bk-form class="info-form clearfix" :label-width="85">
            <bk-form-item class="form-item clearfix fl" :label="$t('任务名称')">
                <span class="form-value">{{task.bk_task_name}}</span>
            </bk-form-item>
            <bk-form-item class="form-item clearfix fl" :label="$t('账户名称')">
                <span class="form-value">{{task.bk_account_id}}</span>
            </bk-form-item>
            <bk-form-item class="form-item clearfix fl" :label="$t('资源类型')">
                <span class="form-value">{{task.bk_resource_type}}</span>
            </bk-form-item>
            <bk-form-item class="form-item clearfix" :label="$t('云区域设定')"></bk-form-item>
        </bk-form>
        <div class="info-table">
            <bk-table :data="task.bk_sync_vpcs">
                <bk-table-column label="VPC" prop="bk_vpc_id" width="200" :formatter="vpcFormatter"></bk-table-column>
                <bk-table-column :label="$t('地域')" prop="bk_region_name"></bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop="bk_host_count"></bk-table-column>
                <bk-table-column :label="$t('主机录入到')" prop="directory" width="250">
                    <template slot-scope="{ row }">
                        <task-directory-selector class="form-table-selector"
                            display="info"
                            :value="row.bk_sync_dir">
                        </task-directory-selector>
                    </template>
                </bk-table-column>
            </bk-table>
        </div>
        <div class="info-options" slot="footer" slot-scope="{ sticky }"
            :class="{ 'is-sticky': sticky }">
            <bk-button theme="primary" @click="handleEdit">{{$t('编辑')}}</bk-button>
            <bk-button class="ml10" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
    </cmdb-sticky-layout>
</template>

<script>
    import TaskForm from './task-form.vue'
    import TaskDirectorySelector from './task-directory-selector.vue'
    export default {
        name: 'task-details-info',
        components: {
            TaskDirectorySelector
        },
        props: {
            task: {
                type: Object,
                default: null
            },
            container: {
                type: Object,
                default: null
            }
        },
        data () {
            return {}
        },
        methods: {
            vpcFormatter (row) {
                const vpcId = row.bk_vpc_id
                const vpcName = row.bk_vpc_name
                if (vpcId !== vpcName) {
                    return `${vpcId}(${vpcName})`
                }
                return vpcId
            },
            handleEdit () {
                this.container.show({
                    detailsComponent: TaskForm.name,
                    props: {
                        mode: 'edit',
                        mission: this.$tools.clone(this.mission)
                    }
                })
            },
            handleCancel () {
                this.container.hide()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .info-form {
        margin: 0 30px;
        padding: 10px 0;
    }
    .form-item {
        width: 300px;
        margin: 5px 15px 0 0;
        /deep/ {
            .bk-label {
                position: relative;
                text-align: left;
                padding: 0 10px 0 0;
                @include ellipsis;
                &:after {
                    content: ":";
                    position: absolute;
                    right: 8px;
                    top: 0;
                    line-height: 30px;
                }
            }
        }
    }
    .form-value {
        font-size: 14px;
        color: #313238;
        line-height: 30px;
    }
    .info-table {
        margin: 0 30px;
    }
    .info-options {
        font-size: 0;
        margin-top: 20px;
        padding: 0 30px;
        &.is-sticky {
            margin-top: 0;
            padding: 15px 30px;
            border-top: 1px solid $borderColor;
            background-color: #FAFBFD;
        }
    }
</style>
