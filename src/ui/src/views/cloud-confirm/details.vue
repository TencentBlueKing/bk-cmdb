<template>
    <div class="detail-wrapper">
        <div class="detail-box">
            <ul class="event-form" v-model="curPush">
                <li class="form-item">
                    <label for="" class="label-name">
                        任务名称<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <label>{{curPush.bk_task_name}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        账号类型<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <label>{{curPush.bk_account_type}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        自动同步
                    </label>
                    <div class="item-content">
                        <label>{{curPush.bk_period_type}}</label>
                        <label>{{curPush.bk_period}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        ID<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <label>{{curPush.bk_secret_id}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">
                        Key<span class="color-danger">*</span>
                    </label>
                    <div class="item-content">
                        <label>*************</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">同步资源</label>
                    <div class="item-content">
                        <label>{{curPush.bk_obj_id}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">账号管理员</label>
                    <div class="item-content">
                        <label>{{curPush.bk_account_admin}}</label>
                    </div>
                </li>
                <li class="form-item">
                    <label for="" class="label-name">资源确认</label>
                    <div class="item-content">
                        <input type="checkbox" v-model="curPush.bk_confirm"/>
                        <label>新增需要确认</label>
                        <input type="checkbox" v-model="curPush.bk_attr_confirm" style="margin-left: 30px"/>
                        <label>属性变化需要确认</label>
                    </div>
                </li>
            </ul>
        </div>
        <footer class="footer">
            <bk-button type="primary" :loading="$loading('savePush')" class="btn" @click="edit">{{$t('Common["编辑"]')}}</bk-button>
            <bk-button type="default" class="btn vice-btn" @click="cancel">{{$t('Common["关闭"]')}}</bk-button>
        </footer>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    export default {
        props: {
            curPush: {
                type: Object
            }
        },
        methods: {
            ...mapActions('cloudDiscover', ['addCloudTask']),
            getPeriodSelect () {

            },
            async save () {
                let res = null
                res = await this.addCloudTask({params: this.params, config: {requestId: 'savePush'}})
                this.$emit('saveSuccess')
                this.$success(this.$t('EventPush["保存成功"]'))
            },
            edit () {
                this.$emit('edit', this.curPush)
            },
            cancel () {
                this.$emit('cancel')
            }
        }
    }
</script>

<style lang="scss" scoped>
    .detail-wrapper{
        height: 100%;
        .detail-box{
            padding: 20px 30px;
            height: calc(100% - 63px);
            overflow-y: auto;
            @include scrollbar;
        }
        .event-form{
            .form-item{
                width: 100%;
                margin-bottom: 20px;
                &:after{
                    display: block;
                    content: "";
                    clear: both;
                }
            }
            .label-name{
                position: relative;
                float: left;
                width: 85px;
                text-align: right;
                line-height: 36px;
                font-size: 14px;
                .color-danger{
                    position: absolute;
                    top: 2px;
                    right: -10px;
                }
            }
            .item-content{
                margin-left: 25px;
                width: calc(100% - 110px);
                line-height: 36px;
                float: left;
                span {
                    font-size: 14px;
                }
                &.length-short{
                    position: relative;
                    input{
                        width: 97px;
                        margin-right: 10px;
                    }
                }
            }
        }
        .footer{
            height: 63px;
            line-height: 63px;
            background: #f9f9f9;
            font-size: 0;
            padding-left: 130px;
            .btn{
                margin-right: 11px;
            }
        }
    }
</style>
