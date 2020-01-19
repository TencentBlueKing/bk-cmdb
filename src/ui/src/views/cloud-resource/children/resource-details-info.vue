<template>
    <cmdb-sticky-layout>
        <bk-form class="info-form clearfix" :label-width="85">
            <bk-form-item class="form-item clearfix fl" :label="$t('任务名称')">
                <span class="form-value">{{mission.mission_name}}</span>
            </bk-form-item>
            <bk-form-item class="form-item clearfix fl" :label="$t('账户名称')">
                <span class="form-value">{{mission.account_name}}</span>
            </bk-form-item>
            <bk-form-item class="form-item clearfix fl" :label="$t('资源类型')">
                <span class="form-value">{{mission.resource_type}}</span>
            </bk-form-item>
            <bk-form-item class="form-item clearfix" :label="$t('云区域设定')"></bk-form-item>
        </bk-form>
        <div class="info-table">
            <bk-table :data="mission.list">
                <bk-table-column label="VPC" prop="vpc" width="200"></bk-table-column>
                <bk-table-column :label="$t('地域')" prop="location"></bk-table-column>
                <bk-table-column :label="$t('主机数量')" prop="host_count"></bk-table-column>
                <bk-table-column :label="$t('主机录入到')" prop="folder" width="250"></bk-table-column>
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
    import CloudResourceForm from './resource-form.vue'
    export default {
        name: 'cloud-resource-details-info',
        props: {
            mission: {
                type: Object,
                default: () => ({
                    mission_name: '发现云主机',
                    account_name: '王者荣耀专用账户',
                    resource_type: '主机',
                    list: [{ vpc: 'vpc', location: '广东三区', host_count: 100, folder: '资源池/lol' }, {}, {}, {}, {}]
                })
            },
            container: {
                type: Object,
                required: true
            }
        },
        data () {
            return {}
        },
        methods: {
            handleEdit () {
                this.container.show({
                    detailsComponent: CloudResourceForm.name,
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
