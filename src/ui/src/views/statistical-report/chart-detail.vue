<template>
    <div class="chart-detail">
        <bk-dialog
            :is-show.sync="showDia"
            :quick-close="false"
            :close-icon="false"
            :has-header="false"
            :has-footer="false"
            :width="600"
            :padding="0">
            <div slot="content" class="dialog-content">
                <div class="model-header">
                    <p class="title">{{ editTitle }}</p>
                    <i class="modal-close icon icon-cc-tips-close" @click="closeChart()"></i>
                </div>
                <div class="content clearfix">
                    <div class="content-item">
                        <label class="label-text">
                            {{$t('Operation["图表类型"]')}}：
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="required" :value="true"
                                v-model="chartType" :disabled="openType === 'edit'">
                            <span class="cmdb-radio-text">{{$t('ModelManagement["自定义"]')}}</span>
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="required" :value="false"
                                v-model="chartType" :disabled="openType === 'edit'">
                            <span class="cmdb-radio-text">{{$t('ModelManagement["内置"]')}}</span>
                        </label>
                    </div>
                    <div v-if="chartType">
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('Operation["图表名称"]')}}：
                            </label>
                            <label class="cmdb-form-item">
                                <cmdb-form-bool-input v-model="chartData.name"
                                    v-validate="'required'"
                                    data-vv-name="chartName">
                                </cmdb-form-bool-input>
                                <span class="form-error">{{errors.first('chartName')}}</span>
                            </label>
                        </div>
                        <div class="content-item" v-if="chartData.bk_obj_id !== 'host'">
                            <label class="label-text">
                                {{$t('Operation["统计对象"]')}}：
                            </label>
                            <span class="cmdb-form-item">
                                <cmdb-selector
                                    setting-key="bk_obj_id"
                                    display-key="bk_obj_name"
                                    :list="staticFilter"
                                    v-model="chartData.bk_obj_id"
                                    v-validate="'required'"
                                    name="staticObj">
                                </cmdb-selector>
                                <span class="form-error">{{errors.first('staticObj')}}</span>
                            </span>
                        </div>
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('Operation["统计维度"]')}}：
                            </label>
                            <span class="cmdb-form-item">
                                <cmdb-selector
                                    setting-key="bk_property_id"
                                    display-key="bk_property_name"
                                    :list="filterList"
                                    v-validate="'required'"
                                    name="staticDim"
                                    v-model="chartData.field">
                                </cmdb-selector>
                                <span class="form-error">{{errors.first('staticDim')}}</span>
                            </span>
                        </div>
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('Operation["图表展示"]')}}：
                            </label>
                            <label class="cmdb-form-radio cmdb-radio-small">
                                <input type="radio" name="present" value="pie" v-model="chartData.chart_type">
                                <span class="cmdb-radio-text">{{$t('Operation["饼图"]')}}</span>
                            </label>
                            <label class="cmdb-form-radio cmdb-radio-small">
                                <input type="radio" name="present" value="bar" v-model="chartData.chart_type">
                                <span class="cmdb-radio-text">{{$t('Operation["柱状图"]')}}</span>
                            </label>
                        </div>
                    </div>
                    <div v-else>
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('Operation["图表选择"]')}}：
                            </label>
                            <span class="cmdb-form-item">
                                <cmdb-selector
                                    setting-key="repType"
                                    display-key="name"
                                    :list="seList.disList"
                                    v-model="chartData.report_type"
                                    v-validate="'required'"
                                    name="present"
                                    :disabled="openType === 'edit'">
                                </cmdb-selector>
                                <span class="form-error">{{errors.first('present')}}</span>
                            </span>
                        </div>
                    </div>
                    <div class="content-item">
                        <label class="label-text">
                            {{$t('Operation["图表宽度"]')}}：
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="width" value="50" v-model="chartData.width">
                            <span class="cmdb-radio-text">50%</span>
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="width" value="100" v-model="chartData.width">
                            <span class="cmdb-radio-text">100%</span>
                        </label>
                    </div>
                </div>
                <div class="footer">
                    <bk-button type="primary" @click="confirm">{{$t("Common['保存']")}}</bk-button>
                    <bk-button type="default" @click="closeChart">{{$t("Common['取消']")}}</bk-button>
                </div>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import { mapActions } from 'vuex'

    export default {
        name: 'chart-detail',
        props: {
            openType: {
                type: String,
                default: 'add'
            },
            chartData: {
                type: Object,
                default () {
                    return {
                        report_type: 'custom',
                        name: '',
                        config_id: null,
                        bk_obj_id: 'host',
                        chart_type: 'pie',
                        field: 'bk_os_type',
                        width: '50'
                    }
                }
            }
        },
        data () {
            return {
                seList: {
                    host: [
                        {
                            name: '按操作系统类型统计',
                            repType: 'host_os_chart'
                        },
                        {
                            name: '按业务统计',
                            repType: 'host_biz_chart'
                        },
                        {
                            name: '按云区域统计',
                            repType: 'host_cloud_chart'
                        },
                        {
                            name: '主机数量变化趋势',
                            repType: 'host_change_biz_chart'
                        }
                    ],
                    inst: [
                        {
                            name: '实例数量统计',
                            repType: 'model_inst_chart'
                        },
                        {
                            name: '实例变更统计',
                            repType: 'model_inst_change_chart'
                        }
                    ],
                    disList: []
                },
                demList: [],
                staList: [],
                chartType: true,
                showDia: true,
                hostFilter: ['host', 'module', 'biz'],
                editTitle: this.openType === 'add' ? this.$t('Common["新增"]') : this.$t('Common["编辑"]')
            }
        },
        computed: {
            filterList () {
                return this.demList.filter(item => {
                    return item.bk_property_type === 'enum'
                })
            },
            staticFilter () {
                return this.staList.filter(item => {
                    return this.hostFilter.indexOf(item.bk_obj_id) === -1
                })
            }
        },
        watch: {
            'chartType' () {
                if (this.chartData.bk_obj_id === 'host') this.seList.disList = this.$tools.clone(this.seList.host)
                else this.seList.disList = this.$tools.clone(this.seList.inst)
            },
            'chartData.bk_obj_id' () {
                this.getDemList(this.chartData.bk_obj_id)
            }
        },
        created () {
            this.chartType = this.chartData.report_type === 'custom'
            if (this.chartType && this.chartData.bk_obj_id === 'host') this.getDemList('host')
            else if (this.chartType && this.chartData.bk_obj_id !== 'host') this.getStaList()
        },
        methods: {
            ...mapActions('operationChart', [
                'getStaticObj',
                'getStaticDimeObj',
                'newStatisticalCharts',
                'updateStatisticalCharts'
            ]),
            async getStaList () {
                this.staList = await this.getStaticObj({})
            },
            async getDemList (id) {
                this.demList = await this.getStaticDimeObj({
                    params: {
                        bk_obj_id: id
                    }
                })
            },
            confirm () {
                this.$validator.validateAll().then(result => {
                    if (result) {
                        if (this.openType === 'add') {
                            this.newStatisticalCharts({ params: this.chartData }).then(res => {
                                this.transData(res)
                            })
                        } else {
                            this.updateStatisticalCharts({ params: this.chartData }).then(res => {
                                this.transData(res)
                            })
                        }
                    }
                })
            },
            transData (res) {
                this.chartData.config_id = res
                this.showDia = false
                setTimeout(() => {
                    this.$emit('transData', this.chartData)
                }, 300)
            },
            closeChart () {
                this.showDia = false
                setTimeout(() => {
                    this.$emit('closeChart')
                }, 300)
            }
        }
    }
</script>

<style scoped lang="scss">
    .dialog-content {
        position: relative;
        .model-header{
            padding: 10px;
            background: #eee;
            .modal-close{
                position: absolute;
                right: 5px;
                top: 5px;
                color: $cmdbMainBtnColor;
                cursor: pointer;
            }
            .title {
                font-size: 16px;
                color: #333948;
                line-height: 1;
                text-align: center;
            }
        }
        .content{
            padding: 10px 20px;
            .content-item{
                padding: 10px;
                .label-text {
                    width: 150px;
                    margin-right: 20px;
                }
                .cmdb-form-radio {
                    width: 114px;
                }
                .cmdb-form-item {
                    display: inline-block;
                    margin-right: 10px;
                    width: 319px;
                    vertical-align: middle;
                    position: relative;
                }
            }
        }
    }
    .footer {
        padding: 10px 24px;
        font-size: 0;
        text-align: center;
        button {
            margin-right: 10px;
        }
    }
    .form-error {
        position: absolute;
        top: 100%;
        left: 0;
        line-height: 14px;
        font-size: 12px;
        color: #ff5656;
    }
</style>
