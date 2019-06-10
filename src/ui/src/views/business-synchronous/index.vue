<template>
    <div class="synchronous-wrapper">
        <template v-if="noFindData">
            <div class="no-content">
                <img src="../../assets/images/no-content.png" alt="no-content">
                <p>{{$t("BusinessSynchronous['找不到更新信息']")}}</p>
                <bk-button type="primary" @click="handleGoHome">{{$t("BusinessSynchronous['返回首页']")}}</bk-button>
            </div>
        </template>
        <template v-else>
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
                            @click="handleContentView(process['process_template_id'], index)">
                            <span>{{process['process_template_name']}}</span>
                            <i :class="['badge', { 'has-read': process['has_read'] }]">{{process['service_instance_count'] | badge}}</i>
                        </div>
                    </div>
                </div>
                <div class="tab-content">
                    <section class="tab-pane"
                        v-show="showContentId === process['process_template_id']"
                        v-for="(process, index) in list"
                        :key="index">
                        <div class="change-box">
                            <div class="title">
                                <h3>{{$t("BusinessSynchronous['变更内容']")}}</h3>
                                <span v-if="process['operational_type'] === 'changed'">（{{properties[process['process_template_id']].length}}）</span>
                            </div>
                            <div class="process-name"
                                v-show="process['operational_type'] === 'changed'">
                                {{$t("ProcessManagement['进程名称']")}}：{{process['process_template_name']}}
                            </div>
                            <div class="process-name mb50"
                                v-show="process['operational_type'] === 'added'">
                                {{$t("BusinessSynchronous['模板中新增进程']")}}
                                <span style="font-weight: bold;">{{process['process_template_name']}}</span>
                            </div>
                            <div class="process-name mb50"
                                v-show="process['operational_type'] === 'removed'">
                                <span style="font-weight: bold;">{{process['process_template_name']}}</span>
                                {{$t("BusinessSynchronous['从模版中删除']")}}
                            </div>
                            <div class="process-info clearfix" v-show="process['operational_type'] === 'changed'">
                                <div class="info-item fl"
                                    v-for="(attribute, attributeIndex) in properties[process['process_template_id']]"
                                    :key="attributeIndex">
                                    {{`${attribute['property_name']}：${attribute['template_property_value']}`}}
                                </div>
                            </div>
                        </div>
                        <div class="instances-box">
                            <div class="title">
                                <h3>{{$t("BusinessSynchronous['涉及实例']")}}</h3>
                                <span>（2）</span>
                            </div>
                            <div class="service-instances">
                                <div class="instances-item"
                                    v-for="(instance, instanceIndex) in process['service_instances']"
                                    :key="instanceIndex"
                                    @click="hanldeInstanceDetails(instance, process['operational_type'])">
                                    <h6>{{instance['service_instance']['name']}}</h6>
                                    <span v-if="process['operational_type'] === 'changed'">（{{instance['changed_attributes'].length}}）</span>
                                </div>
                            </div>
                        </div>
                    </section>
                </div>
            </div>
            <div class="btn-box">
                <bk-button
                    class="mr10"
                    :disabled="readNum !== list.length"
                    type="primary">
                    {{$t("BusinessSynchronous['确认并同步']")}}
                </bk-button>
                <bk-button>{{$t("Common['取消']")}}</bk-button>
            </div>
        </template>

        <cmdb-slider
            :width="676"
            :is-show.sync="slider.show"
            :title="slider.title">
            <template slot="content">
                <instance-details :attribute-list="slider.details"></instance-details>
            </template>
        </cmdb-slider>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    import instanceDetails from './children/details.vue'
    import imitationData from './data'
    export default {
        components: {
            instanceDetails
        },
        filters: {
            badge (value) {
                return value > 99 ? '99+' : value
            }
        },
        data () {
            return {
                slider: {
                    show: false,
                    title: '',
                    details: {}
                },
                noFindData: false,
                showContentId: null,
                readNum: 1,
                modelProperties: []
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'featureTipsParams']),
            list () {
                const formatList = []
                Object.keys(imitationData.data).forEach(key => {
                    formatList.push(...imitationData.data[key].map(info => {
                        return {
                            operational_type: key,
                            has_read: false,
                            ...info
                        }
                    }).filter(process => process['operational_type'] !== 'unchanged'))
                })
                return formatList
            },
            properties () {
                const changedList = this.list.filter(process => process['operational_type'] === 'changed')
                const attributesSet = {}
                changedList.forEach(process => {
                    const attributes = []
                    process['service_instances'].map(instance => {
                        instance['changed_attributes'].forEach(attribute => {
                            if (!attributes.filter(item => item['property_id'] === attribute['property_id']).length) {
                                attributes.push(attribute)
                            }
                        })
                    })
                    attributesSet[process['process_template_id']] = attributes
                })
                return attributesSet
            }
        },
        created () {
            this.showContentId = this.list[0]['process_template_id']
            this.$set(this.list[0], 'has_read', true)
            this.getModaelProperty()
            // console.log(this.list)
        },
        methods: {
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('businessSynchronous', [
                'searchServiceInstanceDifferences',
                'syncServiceInstanceByTemplate'
            ]),
            ...mapActions('processInstance', ['getServiceInstanceProcesses']),
            async getModaelProperty () {
                this.modelProperties = await this.searchObjectAttribute({
                    params: this.$injectMetadata({
                        bk_obj_id: 'process',
                        bk_supplier_account: this.supplierAccount
                    }),
                    config: {
                        requestId: `post_searchObjectAttribute_process`,
                        fromCache: false
                    }
                })
                console.log(this.modelProperties)
            },
            propertiesGroup () {

            },
            getServiceInstanceDifferences () {
                this.searchServiceInstanceDifferences({
                    params: this.$injectMetadata({
                        module_id: '',
                        service_template_id: ''
                    })
                })
            },
            handleContentView (id, index) {
                this.showContentId = id
                if (!this.list[index]['has_read']) {
                    this.$set(this.list[index], 'has_read', true)
                    this.readNum++
                }
            },
            hanldeInstanceDetails (instance, type) {
                if (type === 'changed') {
                    this.slider.title = instance['service_instance']['name']
                    this.slider.show = true
                    this.slider.details = instance['changed_attributes']
                } else {
                    this.getServiceInstanceProcesses({
                        params: this.$injectMetadata({
                            service_instance_id: 67
                        })
                    }).then(data => {
                        console.log(data[0])
                        const showProperty = data[0]['property']
                        console.log(showProperty)
                        this.propertiesGroup(showProperty)
                    })
                }
            },
            handleGoHome () {
                this.$router.push({ name: 'index' })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .synchronous-wrapper {
        position: relative;
        color: #63656e;
        .no-content {
            position: absolute;
            top: 50%;
            left: 50%;
            font-size: 16px;
            color: #63656e;
            text-align: center;
            transform: translate(-50%, -50%);
            img {
                width: 130px;
            }
            p {
                padding: 20px 0 30px;
            }
        }
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
                        margin-right: -14px;
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
                            &:hover {
                                color: #3a84ff;
                                border-color: #3a84ff;
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
