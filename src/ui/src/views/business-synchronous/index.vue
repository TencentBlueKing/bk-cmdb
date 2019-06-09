<template>
    <div class="synchronous-wrapper">
        <p class="tips">请确认以下模版更新信息：</p>
        <div class="info-tab">
            <div class="tab-head">
                <div class="tab-nav">
                    <div :class="['nav-item', {
                             'delete-item': process['operational_type'] === 'removed',
                             'active': showContentId === process['process_template_id']
                         }]"
                        v-for="(process, index) in list"
                        :key="index"
                        @click="handleContentView(process['process_template_id'])">
                        <span>{{process['process_template_name']}}</span>
                        <i class="badge has-read">{{process['service_instance_count'] | badge}}</i>
                    </div>
                </div>
            </div>
            <div class="tab-content">
                <section class="tab-pane"
                    v-if="showContentId === process['process_template_id']"
                    v-for="(process, index) in list"
                    :key="index">
                    <div class="change-box">
                        <div class="title">
                            <h3>变更内容</h3>
                            <span>（2）</span>
                        </div>
                        <div class="process-name"
                            v-if="process['operational_type'] === 'changed'">
                            进程名称：{{process['process_template_name']}}
                        </div>
                        <div class="process-info clearfix">
                            <div class="info-item fl"
                                v-for="(instance, instanceIndex) in process['service_instances']"
                                :key="instanceIndex">
                                <!-- 进程名称：模版进程三 -->
                            </div>
                        </div>
                    </div>
                    <div class="instances-box">
                        <div class="title">
                            <h3>涉及实例</h3>
                            <span>（2）</span>
                        </div>
                        <div class="service-instances">
                            <div class="instances-item"
                                v-for="(instance, instanceIndex) in process['service_instances']"
                                :key="instanceIndex">
                                <h6>{{instance['service_instance']['name']}}</h6>
                                <span>（{{instance['changed_attributes'].length}}）</span>
                            </div>
                        </div>
                    </div>
                </section>
            </div>
        </div>
        <div class="btn-box">
            <bk-button class="mr10" type="primary">{{$t("BusinessSynchronous['确认并同步']")}}</bk-button>
            <bk-button>{{$t("Common['取消']")}}</bk-button>
        </div>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'
    import imitationData from './data'
    export default {
        filters: {
            badge (value) {
                return value > 99 ? '99+' : value
            }
        },
        data () {
            return {
                originData: {},
                showContentId: null
            }
        },
        computed: {
            list () {
                const formatList = []
                Object.keys(imitationData.data).forEach(key => {
                    formatList.push(...imitationData.data[key].map(info => {
                        return {
                            operational_type: key,
                            ...info
                        }
                    }).filter(process => process['operational_type'] !== 'unchanged'))
                })
                return formatList
            }
        },
        created () {
            this.showContentId = this.list[0]['process_template_id']
            console.log(this.list)
        },
        methods: {
            ...mapActions('businessSynchronous', [
                'searchServiceInstanceDifferences',
                'syncServiceInstanceByTemplate'
            ]),
            getFormatList () {
                const formatList = []
                Object.keys(imitationData.data).forEach(key => {
                    formatList.push(...imitationData.data[key].map(info => {
                        return {
                            operational_type: key,
                            ...info
                        }
                    }))
                })
                console.log(formatList)
            },
            getServiceInstanceDifferences () {
                this.searchServiceInstanceDifferences({
                    params: this.$injectMateData({
                        bk_module_id: '',
                        service_template_id: ''
                    })
                })
            },
            handleContentView (id) {
                this.showContentId = id
            }
        }
    }
</script>

<style lang="scss" scoped>
    .synchronous-wrapper {
        color: #63656e;
        .tips {
            padding-bottom: 20px;
        }
        .info-tab {
            @include space-between;
            height: 500px;
            border: 1px solid #c3cdd7;
            .tab-head {
                height: 100%;
                .tab-nav {
                    @include scrollbar-y;
                    width: 200px;
                    height: 100%;
                    background-color: #fafbfd;
                    padding-bottom: 20px;
                    border-right: 1px solid #c3cdd7;
                }
                .nav-item {
                    @include space-between;
                    position: relative;
                    height: 60px;
                    padding: 0px 14px;
                    border-bottom: 1px solid #c3cdd7;
                    cursor: pointer;
                    &.delete-item span {
                        text-decoration: line-through;
                    }
                    span {
                        @include ellipsis;
                        flex: 1;
                        padding-right: 10px;
                        font-size: 14px;
                    }
                    .badge {
                        display: inline-block;
                        width: 56px;
                        height: 36px;
                        line-height: 36px;
                        font-size: 20px;
                        font-style: normal;
                        font-weight: bold;
                        text-align: center;
                        background-color: #ff5656;
                        color: #ffffff;
                        border-radius: 20px;
                        transform: scale(.5);
                        &.has-read {
                            color: #ffffff;
                            background-color: #c4c6cc;
                        }
                    }
                    &.active {
                        color: #3a84ff;
                        background-color: #ffffff;
                        span {
                            font-weight: bold;
                        }
                        &::after {
                            content: '';
                            position: absolute;
                            top: 0;
                            right: -1px;
                            width: 1px;
                            height: 100%;
                            background-color: #ffffff;
                        }
                        &.delete-item {
                            color: #ff5656;
                        }
                    }
                }
            }
            .tab-content {
                flex: 1;
                height: 100%;
                overflow: hidden;
                .tab-pane {
                    font-size: 14px;
                    padding: 20px 20px 20px 38px;
                    .title {
                        display: flex;
                        align-items: center;
                        padding-bottom: 24px;
                        h3 {
                            font-size: 14px;
                        }
                        span {
                            color: #c4c6cc;
                        }
                    }
                    .change-box {
                        color: #313238;
                        .process-info {
                            padding-top: 20px;
                            padding-bottom: 30px;
                            .info-item {
                                @include ellipsis;
                                width: 33.333%;
                                padding-right: 20px;
                                padding-bottom: 20px;
                            }
                        }
                    }
                    .service-instances {
                        @include scrollbar-y;
                        max-height: 186px;
                        display: flex;
                        flex-wrap: wrap;
                        .instances-item {
                            @include space-between;
                            width: 240px;
                            font-size: 12px;
                            padding: 2px 6px;
                            margin-bottom: 16px;
                            margin-right: 14px;
                            border: 1px solid #dcdee5;
                            background-color: #fafbfd;
                            cursor: pointer;
                            h6 {
                                @include ellipsis;
                                flex: 1;
                                padding-right: 4px;
                                font-weight: normal;
                            }
                        }
                    }
                }
            }
        }
        .btn-box {
            padding-top: 20px;
        }
    }
</style>
