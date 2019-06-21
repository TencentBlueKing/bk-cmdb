<template>
    <div class="synchronous-wrapper">
        <template v-if="noFindData">
            <div class="no-content">
                <img src="../../assets/images/no-content.png" alt="no-content">
                <p>{{$t("BusinessSynchronous['找不到更新信息']")}}</p>
                <bk-button type="primary" @click="handleGoBackModule">{{$t("Common['返回']")}}</bk-button>
            </div>
        </template>
        <template v-else-if="isLatsetData">
            <div class="no-content">
                <img src="../../assets/images/latset-data.png" alt="no-content">
                <p>{{$t("BusinessSynchronous['最新数据']")}}</p>
                <bk-button type="primary" @click="handleGoBackModule">{{$t("Common['返回']")}}</bk-button>
            </div>
        </template>
        <template v-else>
            <feature-tips
                :show-tips="showFeatureTips"
                :desc="$t('BusinessSynchronous[\'功能提示\']')">
            </feature-tips>
            <p class="tips" :style="{ 'padding-top': showFeatureTips ? '24px' : '0' }">
                {{$t("BusinessSynchronous['请确认']")}}
                <span>{{treePath}}</span>
                {{$t("BusinessSynchronous['模版更新信息']")}}
            </p>
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
                                    {{`${attribute['property_name']}：${attribute['show_value'] ? attribute['show_value'] : '--'}`}}
                                </div>
                            </div>
                        </div>
                        <div class="instances-box">
                            <div class="title">
                                <h3>{{$t("BusinessSynchronous['涉及实例']")}}</h3>
                                <span>（{{process['service_instances'].length}}）</span>
                            </div>
                            <div class="service-instances">
                                <div class="instances-item"
                                    v-for="(instance, instanceIndex) in process['service_instances']"
                                    :key="instanceIndex"
                                    @click="hanldeInstanceDetails(instance, process['operational_type'], process['process_template_id'])">
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
                    type="primary"
                    @click="handleSubmitSync">
                    {{$t("BusinessSynchronous['确认并同步']")}}
                </bk-button>
                <bk-button @click="handleGoBackModule">{{$t("Common['取消']")}}</bk-button>
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
    import { mapGetters, mapActions, mapMutations } from 'vuex'
    import instanceDetails from './children/details.vue'
    import featureTips from '@/components/feature-tips/index'
    export default {
        components: {
            instanceDetails,
            featureTips
        },
        filters: {
            badge (value) {
                return value > 99 ? '99+' : value
            }
        },
        data () {
            return {
                showFeatureTips: true,
                viewsTitle: '',
                noFindData: true,
                isLatsetData: true,
                showContentId: null,
                readNum: 1,
                serviceTemplateId: '',
                differenData: {},
                modelProperties: [],
                changedData: {
                    instanceDetails: {},
                    type: 'changed',
                    current: {}
                },
                slider: {
                    show: false,
                    title: '',
                    details: {}
                }
            }
        },
        computed: {
            ...mapGetters(['supplierAccount', 'featureTipsParams']),
            business () {
                return this.$store.getters['objectBiz/bizId']
            },
            routerParams () {
                return this.$route.params
            },
            treePath () {
                return this.$route.query.path
            },
            list () {
                const formatList = []
                Object.keys(this.differenData).forEach(key => {
                    formatList.push(...this.differenData[key].map(info => {
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
                                const property = this.modelProperties.find(property => property['bk_property_id'] === attribute['property_id'])
                                if (['enum'].includes(property['bk_property_type'])) {
                                    attribute['show_value'] = property['option'].find(option => option['id'] === attribute['template_property_value']['value'])['name']
                                } else if (['bool'].includes(property['bk_property_type'])) {
                                    attribute['show_value'] = attribute['template_property_value']['value'] ? '是' : '否'
                                } else {
                                    attribute['show_value'] = attribute['template_property_value']['value']
                                }
                                attributes.push(attribute)
                            }
                        })
                    })
                    attributesSet[process['process_template_id']] = attributes
                })
                return attributesSet
            },
            instanceIds () {
                const ids = []
                this.list.forEach(item => {
                    item['service_instances'].forEach(instance => {
                        ids.push(instance['service_instance']['id'])
                    })
                })
                return ids
            },
            instanceMap () {
                return this.$store.state.businessSync.instanceMap
            }
        },
        async created () {
            this.$store.commit('setHeaderTitle', '')
            await this.getModaelProperty()
            await this.getModuleInstance()
            if (this.list.length) {
                this.isLatsetData = false
                this.showContentId = this.list[0]['process_template_id']
                this.$set(this.list[0], 'has_read', true)
            } else {
                this.isLatsetData = true
            }
        },
        methods: {
            ...mapMutations('businessSynchronous', ['setInstance']),
            ...mapActions('objectModelProperty', ['searchObjectAttribute']),
            ...mapActions('processInstance', ['getServiceInstanceProcesses']),
            ...mapActions('processTemplate', ['getBatchProcessTemplate']),
            ...mapActions('businessSynchronous', [
                'searchServiceInstanceDifferences',
                'syncServiceInstanceByTemplate'
            ]),
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
            },
            async getModuleInstance () {
                const data = await this.$store.dispatch('objectModule/searchModule', {
                    bizId: this.business,
                    setId: Number(this.routerParams.setId),
                    params: {
                        page: { start: 0, limit: 1 },
                        fields: [],
                        condition: {
                            bk_module_id: Number(this.routerParams.moduleId),
                            bk_supplier_account: this.supplierAccount
                        }
                    },
                    config: {
                        requestId: 'getNodeInstance',
                        cancelPrevious: true
                    }
                })
                if (data.info.length) {
                    this.noFindData = false
                    const instance = data.info[0]
                    this.serviceTemplateId = instance['service_template_id']
                    this.viewsTitle = instance['bk_module_name']
                    await this.getServiceInstanceDifferences()
                } else {
                    this.noFindData = true
                }
            },
            async getServiceInstanceDifferences () {
                try {
                    this.differenData = await this.searchServiceInstanceDifferences({
                        params: this.$injectMetadata({
                            bk_module_id: Number(this.routerParams.moduleId),
                            service_template_id: this.serviceTemplateId
                        })
                    })
                    this.$store.commit('setHeaderTitle', `${this.$t("BusinessSynchronous['同步模板']")}【${this.viewsTitle}】`)
                } catch (error) {
                    console.error(error)
                    this.noFindData = true
                }
            },
            propertiesGroup () {
                const instance = this.changedData.instanceDetails
                return Object.keys(instance).filter(propertyKey => this.modelProperties.find(property => property['bk_property_id'] === propertyKey))
                    .map(key => {
                        const property = this.modelProperties.find(property => property['bk_property_id'] === key)
                        let propertyValue = ''
                        if (['enum'].includes(property['bk_property_type'])) {
                            propertyValue = property['option'].find(option => option['id'] === instance[key])['name']
                        } else if (['bool'].includes(property['bk_property_type'])) {
                            propertyValue = instance[key] ? '是' : '否'
                        } else {
                            propertyValue = instance[key]
                        }
                        return {
                            id: property['id'],
                            property_id: property['bk_property_id'],
                            property_name: property['bk_property_name'],
                            before_value: this.changedData.type === 'added' ? '--' : propertyValue,
                            show_value: this.changedData.type === 'removed' ? this.$t("BusinessSynchronous['该进程已删除']") : propertyValue
                        }
                    })
            },
            filterShowList () {
                const list = this.propertiesGroup()
                if (this.changedData.type === 'added') {
                    return list.filter(property => property['show_value'])
                } else {
                    return list.filter(property => property['before_value'])
                }
            },
            handleContentView (id, index) {
                this.showContentId = id
                if (!this.list[index]['has_read']) {
                    this.$set(this.list[index], 'has_read', true)
                    this.readNum++
                }
            },
            getTableShowList (list) {
                return list.map(item => {
                    const result = item
                    const property = this.modelProperties.find(property => property['bk_property_id'] === item['property_id'])
                    if (['enum'].includes(property['bk_property_type'])) {
                        result['before_value'] = property['option'].find(option => option['id'] === item['property_value'])['name']
                        result['show_value'] = property['option'].find(option => option['id'] === item['template_property_value']['value'])['name']
                    } else if (['bool'].includes(property['bk_property_type'])) {
                        result['before_value'] = item['property_value'] ? '是' : '否'
                        result['show_value'] = item['template_property_value']['value'] ? '是' : '否'
                    } else {
                        result['before_value'] = item['property_value']
                        result['show_value'] = item['template_property_value']['value'] ? item['template_property_value']['value'] : '--'
                    }
                    return result
                })
            },
            async hanldeInstanceDetails (instance, type, processId) {
                const instanceId = instance['service_instance']['id']
                this.slider.title = instance['service_instance']['name']
                this.changedData.current = instance['changed_attributes']
                this.changedData.type = type
                if (type === 'changed') {
                    this.slider.details = this.getTableShowList(instance['changed_attributes'])
                } else if (type === 'remove') {
                    if (this.instanceMap.hasOwnProperty(instanceId)) {
                        this.changedData.instanceDetails = this.instanceMap[instanceId]
                    } else {
                        try {
                            const result = await this.getServiceInstanceProcesses({
                                params: this.$injectMetadata({
                                    service_instance_id: instanceId
                                })
                            })
                            this.changedData.instanceDetails = result[0]['property']
                            this.$store.commit('businessSync/setInstance', {
                                id: instanceId,
                                instanceProperty: result[0]['property']
                            })
                        } catch (e) {
                            console.error(e)
                        }
                    }
                    this.slider.details = this.filterShowList()
                } else {
                    try {
                        const result = await this.getBatchProcessTemplate({
                            params: this.$injectMetadata({
                                service_template_id: instance['service_instance']['service_template_id']
                            })
                        })
                        const processProperties = result.info.find(process => process['id'] === processId)['property']
                        const instanceDetails = {}
                        Object.keys(processProperties).forEach(key => {
                            instanceDetails[key] = processProperties[key]['value']
                        })
                        this.changedData.instanceDetails = instanceDetails
                    } catch (e) {
                        console.error(e)
                    }
                    this.slider.details = this.filterShowList()
                }
                this.slider.show = true
            },
            handleSubmitSync () {
                this.syncServiceInstanceByTemplate({
                    params: this.$injectMetadata({
                        service_template_id: this.serviceTemplateId,
                        bk_module_id: Number(this.routerParams.moduleId),
                        service_instances: this.instanceIds
                    })
                }).then(() => {
                    this.$success(this.$t('BusinessSynchronous["同步成功"]'))
                    this.handleGoBackModule()
                })
            },
            handleGoBackModule () {
                this.$router.replace({
                    name: 'topology',
                    query: {
                        module: this.routerParams.moduleId
                    }
                })
            }
        }
    }
</script>

<style lang="scss" scoped>
    .synchronous-wrapper {
        position: relative;
        color: #63656e;
        padding-top: 10px;
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
            span {
                font-weight: bold;
            }
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
                        font-size: 16px;
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
                            font-size: 16px;
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
                            font-size: 14px;
                            padding: 2px 6px 4px;
                            margin-bottom: 16px;
                            margin-right: 14px;
                            border: 1px solid #dcdee5;
                            background-color: #fafbfd;
                            cursor: pointer;
                            h6 {
                                @include ellipsis;
                                flex: 1;
                                font-size: 14px;
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
