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
                        <label class="cmdb-form-radio cmdb-radio-big">
                            <input type="radio" name="required" :value="true"
                                v-model="chartType" :disabled="openType === 'edit'">
                            <span class="cmdb-radio-text">{{$t('ModelManagement["自定义"]')}}</span>
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-big">
                            <input type="radio" name="required" :value="false"
                                v-model="chartType" :disabled="openType === 'edit'">
                            <span class="cmdb-radio-text">{{$t('ModelManagement["内置"]')}}</span>
                        </label>
                    </div>
                    <div v-if="!chartType">
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('Operation["图表名称"]')}}：
                            </label>
                            <span class="cmdb-form-item">
                                <cmdb-selector
                                    setting-key="name"
                                    display-key="name"
                                    :list="seList.disList"
                                    v-model="chartData.name"
                                    v-validate="'required'"
                                    name="present"
                                    :disabled="openType === 'edit'">
                                </cmdb-selector>
                                <span class="form-error">{{errors.first('present')}}</span>
                            </span>
                        </div>
                    </div>
                    <div v-if="chartType">
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('Operation["图表名称"]')}}：
                            </label>
                            <span class="cmdb-form-item">
                                <input class="cmdb-form-input" placeholder="请输入图表名称" v-model="chartData.name">
                                <span class="form-error">{{errors.first('present')}}</span>
                            </span>
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
                                {{$t('Operation["图表类型"]')}}：
                            </label>
                            <label class="cmdb-form-radio cmdb-radio-big">
                                <input type="radio" name="present" value="pie" v-model="chartData.chart_type">
                                <span class="cmdb-radio-text cmdb-radio-text-icon"><i class="icon icon-cc-op-pie"></i>{{$t('Operation["柱状图"]')}}</span>
                            </label>
                            <label class="cmdb-form-radio cmdb-radio-big">
                                <input type="radio" name="present" value="bar" v-model="chartData.chart_type">
                                <span class="cmdb-radio-text cmdb-radio-text-icon"><i class="icon icon-cc-op-bar"></i>{{$t('Operation["柱状图"]')}}</span>
                            </label>
                        </div>
                    </div>
                    <div class="content-item">
                        <label class="label-text">
                            {{$t('Operation["图表宽度"]')}}：
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-big">
                            <input type="radio" name="width" value="50" v-model="chartData.width">
                            <span class="cmdb-radio-text">50%</span>
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-big">
                            <input type="radio" name="width" value="100" v-model="chartData.width">
                            <span class="cmdb-radio-text">100%</span>
                        </label>
                    </div>
                    <div class="content-item">
                        <label class="label-text-x">
                            {{$t('Operation["横轴坐标数量"]')}}：
                        </label>
                        <label class="cmdb-form-item">
                            <div class="axis-picker">
                                <input class="cmdb-form-input form-input" v-model="cal">
                                <i class="bk-icon icon-angle-down" @click="calculate('down')"></i>
                                <i class="bk-icon icon-angle-up" @click="calculate('up')"></i>
                            </div>
                            <span class="form-error">{{errors.first('chartName')}}</span>
                        </label>
                        <span class="tips">{{$t('Operation["考虑显示效果，上限为25个，100%宽度建议显示20个，50%宽度10个"]')}}</span>
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
            hostType: {
                type: String,
                default: 'host'
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
                cal: 10,
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
                hostFilter: ['host', 'module', 'biz', 'set', 'process'],
                editTitle: this.openType === 'add' ? this.$t('Operation["新增图表"]') : this.$t('Operation["编辑图表"]')
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
            calculate (flag) {
                if (flag === 'up') {
                    this.cal += 1
                    if (this.cal > 25) {
                        this.cal = 25
                    }
                } else {
                    this.cal -= 1
                    if (this.cal < 1) {
                        this.cal = 1
                    }
                }
            },
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
            background:white;
            margin-bottom:25px;
            .modal-close{
                position: absolute;
                right: 10px;
                top: 15px;
                color:#D8D8D8;
                cursor: pointer;
            }
            .title {
                display: inline-block;
                float:left;
                font-size:24px;
                font-family:MicrosoftYaHei;
                color:rgba(68,68,68,1);
                margin-left: 14px;
            }
        }
        .content{
            padding: 10px 20px;
            margin-left:40px;
            .content-item{
                padding: 10px;
                .label-text {
                    width: 150px;
                    margin-right: 50px;
                }
                .label-text-x{
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
                    .axis-picker{
                        position:relative;
                        width:120px;
                        height:32px;
                        i{
                            font-size:12px;
                            position:absolute;
                            right:8px;
                            &:nth-child(2){
                                bottom:4px
                             }
                            &:nth-child(3){
                                 top: 4px;
                             }
                        }
    }
                }
                .form-input{
                    float:left;
                    width:120px;
                    height:32px;
                }
                .cmdb-radio-text-icon{
                    i{
                        vertical-align: middle;
                        line-height: 19px;
                        font-size:16px;
                        margin-left:5px;
                        border:1px dashed grey;
                    }
                }
                .tips{
                    display:block;
                    width:380px;
                    height:15px;
                    font-size:11px;
                    font-family:MicrosoftYaHei;
                    color:rgba(151,155,165,1);
                    line-height:15px;
                    margin-left: 122px;
                    margin-top: 6px;
                }
            }
        }
    }
    .footer {
        border-top: 1px solid rgba(220,222,229,1);
        height:50px;
        line-height: 50px;
        font-size: 0;
        text-align: center;
        button {
            float:right;
            margin-right: 24px;
            margin-top:7px;
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
