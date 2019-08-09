<template>
    <div>
        <bk-dialog v-model="showDia"
            class="bk-dialog-no-padding bk-dialog-no-tools"
            :close-icon="false"
            :mask-close="false"
            :show-footer="false"
            :width="720">
            <div class="dialog-content">
                <div class="model-header">
                    <p class="title">{{ editTitle }}</p>
                    <i class="modal-close bk-icon icon-close" @click="closeChart()"></i>
                </div>
                <div class="content clearfix">
                    <div class="content-item">
                        <label class="label-text">
                            {{$t('图表类型')}}
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-big">
                            <input type="radio" name="required" :value="true"
                                v-model="chartType" :disabled="openType === 'edit'">
                            <span class="cmdb-radio-text">{{$t('自定义')}}</span>
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-big">
                            <input type="radio" name="required" :value="false"
                                v-model="chartType" :disabled="openType === 'edit'">
                            <span class="cmdb-radio-text">{{$t('内置')}}</span>
                        </label>
                    </div>
                    <div v-if="!chartType">
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('图表名称')}}
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
                                {{$t('图表名称')}}
                            </label>
                            <span class="cmdb-form-item">
                                <input class="cmdb-form-input" :placeholder="$t('请输入图表名称')" v-model="chartData.name" name="collectionName" v-validate="'required'">
                                <span class="form-error">{{errors.first('collectionName')}}</span>
                            </span>
                        </div>
                        <div class="content-item" v-if="chartData.bk_obj_id !== 'host'">
                            <label class="label-text">
                                {{$t('统计对象')}}
                            </label>
                            <span class="cmdb-form-item">
                                <cmdb-selector
                                    :disabled="openType === 'edit'"
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
                                {{$t('统计维度')}}
                            </label>
                            <span class="cmdb-form-item">
                                <cmdb-selector
                                    :disabled="openType === 'edit'"
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
                                {{$t('图表类型')}}
                            </label>
                            <label class="cmdb-form-radio cmdb-radio-big">
                                <input type="radio" name="present" value="pie" v-model="chartData.chart_type">
                                <span class="cmdb-radio-text cmdb-radio-text-icon">{{$t('饼图')}}</span>
                            </label>
                            <label class="cmdb-form-radio cmdb-radio-big">
                                <input type="radio" name="present" value="bar" v-model="chartData.chart_type">
                                <span class="cmdb-radio-text cmdb-radio-text-icon">{{$t('柱状图')}}</span>
                            </label>
                        </div>
                    </div>
                    <div class="content-item">
                        <label class="label-text">
                            {{$t('图表宽度')}}
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
                            <bk-popover placement="right" :content="$t('该图表可视区横轴坐标数量')">
                                <span class="text-underline">{{$t('横轴坐标数量')}}</span>
                            </bk-popover>
                        </label>
                        <label class="cmdb-form-item">
                            <div class="axis-picker">
                                <input class="cmdb-form-input form-input"
                                    v-validate="'required|number'" name="chartNumber"
                                    v-model="chartData.x_axis_count">
                                <i class="bk-icon icon-angle-down" @click="calculate('down')"></i>
                                <i class="bk-icon icon-angle-up" @click="calculate('up')" v-if="maxNum !== chartData.x_axis_count"></i>
                                <bk-popover class="tool-tip" placement="right" :content="$t('已经超出可显示的最大数量')" v-if="maxNum <= chartData.x_axis_count">
                                    <i class="bk-icon icon-angle-up" @click="calculate('up')"></i>
                                </bk-popover>
                            </div>
                            <span class="tips">{{$t('考虑显示效果，上限为25个，100%宽度建议显示20个，50%宽度10个')}}
                                <span class="form-error">{{errors.first('chartNumber')}}</span>
                            </span>
                        </label>
                    </div>
                </div>
                <div class="footer" slot="footer">
                    <bk-button theme="primary" @click="confirm">{{$t('确定')}}</bk-button>
                    <bk-button theme="default" @click="closeChart">{{$t('取消')}}</bk-button>
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
                        field: '',
                        width: '50',
                        x_axis_count: 10
                    }
                }
            }
        },
        data () {
            return {
                seList: {
                    host: [
                        {
                            name: this.$t('按操作系统类型统计'),
                            repType: 'host_os_chart'
                        },
                        {
                            name: this.$t('按业务统计'),
                            repType: 'host_biz_chart'
                        },
                        {
                            name: this.$t('按云区域统计'),
                            repType: 'host_cloud_chart'
                        },
                        {
                            name: this.$t('主机数量变化趋势'),
                            repType: 'host_change_biz_chart'
                        }
                    ],
                    inst: [
                        {
                            name: this.$t('实例数量统计'),
                            repType: 'model_inst_chart'
                        },
                        {
                            name: this.$t('实例变更统计'),
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
                editTitle: '',
                maxNum: 0
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
            },
            typeFilter () {
                const data = this.chartData.bk_obj_id === 'host' ? this.seList.host : this.seList.inst
                return data.filter(item => {
                    return item.name === this.chartData.name
                })
            }
        },
        watch: {
            'chartType' () {
                this.$validator.reset()
                if (this.chartData.bk_obj_id === 'host') this.seList.disList = this.$tools.clone(this.seList.host)
                else this.seList.disList = this.$tools.clone(this.seList.inst)
                if (this.chartType) this.chartData.name = ''
            },
            'chartData.bk_obj_id' () {
                this.getDemList(this.chartData.bk_obj_id)
            }
        },
        created () {
            if (this.openType !== 'add') this.maxNum = this.chartData.chart_type === 'pie' ? this.chartData.data.data[0].labels.length : this.chartData.data.data[0].x.length
            else this.maxNum = 25
            this.initTitle()
            this.chartType = this.chartData.report_type === 'custom'
            this.getDemList(this.chartData.bk_obj_id)
            if (this.chartType && this.chartData.bk_obj_id !== 'host') this.getStaList()
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
                    this.chartData.x_axis_count += 1
                    this.maxNum = parseInt(this.maxNum) >= 25 ? 25 : this.maxNum
                    if (this.chartData.x_axis_count > this.maxNum) {
                        this.chartData.x_axis_count = this.maxNum
                    }
                } else {
                    this.chartData.x_axis_count -= 1
                    if (this.chartData.x_axis_count < 1) {
                        this.chartData.x_axis_count = 1
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
                this.$validator.reset()
                if (this.openType === 'add') this.chartData.field = ''
            },
            confirm () {
                this.chartData.report_type = this.chartType ? 'custom' : this.typeFilter[0].repType
                this.chartData.x_axis_count = parseInt(this.chartData.x_axis_count)
                const data = this.$tools.clone(this.chartData)
                this.$validator.validateAll().then(result => {
                    if (result) {
                        if (this.openType === 'add') {
                            if (!this.chartType) this.delKeys(data, ['bk_obj_id', 'config_id', 'field', 'name', 'chart_type'])
                            this.newStatisticalCharts({ params: data }).then(res => {
                                this.transData(res.info)
                            })
                        } else {
                            this.delKeys(data, ['data', 'hasData', 'create_time', 'title'])
                            this.updateStatisticalCharts({ params: data }).then(res => {
                                this.chartData.config_id = res
                                this.transData(this.chartData)
                            })
                        }
                    }
                })
            },
            transData (res) {
                this.showDia = false
                setTimeout(() => {
                    this.$emit('transData', res)
                }, 300)
            },
            closeChart () {
                this.showDia = false
                setTimeout(() => {
                    this.$emit('closeChart')
                }, 300)
            },
            delKeys (obj, keys) {
                keys.map((key) => {
                    delete obj[key]
                })
                return obj
            },
            initTitle () {
                if (this.openType !== 'add') this.editTitle = this.$t('编辑') + '【' + this.chartData.title + '】'
                else this.editTitle = this.chartData.bk_obj_id === 'host' ? this.$t('新增主机统计图表') : this.$t('新增实例统计图表')
            }
        }
    }
</script>

<style scoped lang="scss">
    .dialog-content {
        position: relative;
        .model-header{
            padding: 15px;
            background:white;
            .modal-close{
                position: absolute;
                right: 10px;
                top: 15px;
                color:#D8D8D8;
                cursor: pointer;
                width: 26px;
                height: 26px;
                line-height: 26px;
                text-align: center;
                border-radius: 50%;
                font-weight: 700;
                &:hover {
                     background-color: #f0f1f5;
                 }
            }
            .title {
                display: inline-block;
                float:left;
                font-size:24px;
                color:rgba(68,68,68,1);
                margin-left: 14px;
            }
        }
        .content{
            padding: 10px 20px;
            margin: 25px 40px;
            .content-item{
                padding: 10px;
                .label-text {
                    width: 150px;
                    margin-right: 50px;
                }
                .label-text-x{
                    width: 150px;
                    margin-right: 20px;
                    .text-underline {
                        border-bottom: 1px dashed #63656E;
                    }
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
                    div {
                        width: 100%;
                    }
                    input {
                        border: 1px solid #c4c6cc;
                        border-radius: 2px;
                        height: 30px!important;
                    }
                    .axis-picker{
                        position:relative;
                        width:120px;
                        height:32px;
                        .tool-tip {
                            font-size:12px;
                            position:absolute;
                            right: -120px;
                            top: -11px;
                        }
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
                    color:rgba(151,155,165,1);
                    line-height:15px;
                    margin-top: 6px;
                    position: absolute;
                }
            }
        }
    }
    .footer {
        border-top: 1px solid rgba(220,222,229,1);
        height: 50px;
        line-height: 50px;
        font-size: 0;
        text-align: right;
        padding-right: 14px;
        background:rgba(250,251,253,1);
        button {
            vertical-align: middle;
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
