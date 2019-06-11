<template>
    <div>
        <div class="operation-top">
            <span class="operation-edit" v-if="!editModel" @click="editModel = !editModel">{{$t('Operation["编辑"]')}}</span>
            <div class="operation-edit" v-if="editModel">
                <span @click="editPage()">{{$t('Common["保存"]')}}</span>
                <span @click="editCancel()">{{$t('Common["取消"]')}}</span>
            </div>
            <span class="operation-title">{{$t('Operation["主机统计"]')}}</span>
        </div>
        <div class="operate-menus">
            <div class="menu-items">
                <div class="item-left">
                    <i class="icon icon-cc-business"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Common["业务"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：6</span>
                </div>
            </div>
            <div class="menu-items">
                <div class="item-left">
                    <i class="icon icon-cc-host"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Hosts["模块"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：23</span>
                </div>
            </div>
            <div class="menu-items">
                <div class="item-left">
                    <i class="icon icon-cc-host"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Nav["主机"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：44</span>
                </div>
            </div>
        </div>
        <div v-for="(item, key) in disCharts"
            :key="item.report_type + item.config_id"
            :style="{ width: item.position.width + '%' }"
            class="operation-layout">
            <div class="operation-charts" :id="item.report_type + item.config_id"></div>
            <div class="charts-options" v-if="editModel">
                <i class="bk-icon icon-arrows-down-shape" @click="moveChart('host', 'down', key)" disabled></i>
                <i class="bk-icon icon-arrows-up-shape" @click="moveChart('host', 'up', key)"></i>
                <i class="bk-icon icon-edit" @click="editCharts('host', item)"></i>
                <i class="icon icon-cc-tips-close" @click="deleteChart('host', key)"></i>
            </div>
        </div>
        <div class="operation-add" v-if="editModel" @click="openNew('host')">
            <i class="bk-icon icon-plus"></i>
        </div>
        <div class="operation-top">
            <span class="operation-title">{{$t('Operation["实例统计"]')}}</span>
        </div>
        <div class="operate-menus operate-menus-bottom">
            <div class="menu-items">
                <div class="item-left">
                    <i class="menu-icon icon-cc-nav-model"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Nav["模型"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：6</span>
                </div>
            </div>
            <div class="menu-items">
                <div class="item-left">
                    <i class="bk-cc-icon icon-cc-business"></i>
                </div>
                <div class="item-right">
                    <span>{{$t('Operation["实例"]')}}</span>
                    <span>{{$t('Index["数量"]')}}：23</span>
                </div>
            </div>
        </div>
        <!--<div v-for="item in charData"-->
        <!--:key="item.report_type + item.config_id"-->
        <!--:style="{ width: item.position.width + '%' }"-->
        <!--class="operation-layout">-->
        <!--<div class="operation-charts" :id="item.report_type + item.config_id"></div>-->
        <!--</div>-->
        <!--<div class="operation-add" v-if="editModel" @click="isShow = !isShow">-->
        <!--<i class="bk-icon icon-plus"></i>-->
        <!--</div>-->
        <div class="operation-add" v-if="editModel" @click="openNew('inst')">
            <i class="bk-icon icon-plus"></i>
        </div>
        <bk-dialog
            :is-show.sync="isShow"
            :quick-close="false"
            :close-icon="false"
            :has-header="false"
            :has-footer="false"
            :width="600"
            :padding="0"
            @cancel="isShow = !isShow">
            <div slot="content" class="dialog-content">
                <div class="model-header">
                    <p class="title">新增</p>
                    <i class="modal-close icon icon-cc-tips-close" @click="cancel()"></i>
                </div>
                <div class="content clearfix">
                    <div class="content-item">
                        <label class="label-text">
                            {{$t('Operation["图表类型"]')}}：
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="required" :value="true" v-model="newChart.type">
                            <span class="cmdb-radio-text">{{$t('ModelManagement["自定义"]')}}</span>
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="required" :value="false" v-model="newChart.type">
                            <span class="cmdb-radio-text">{{$t('ModelManagement["内置"]')}}</span>
                        </label>
                    </div>
                    <div v-if="newChart.type">
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('Operation["图表名称"]')}}：
                            </label>
                            <label class="cmdb-form-item">
                                <cmdb-form-bool-input v-model="newChart.name">
                                </cmdb-form-bool-input>
                            </label>
                        </div>
                        <div class="content-item" v-if="newChart.newType === 'inst'">
                            <label class="label-text">
                                {{$t('Operation["统计对象"]')}}：
                            </label>
                            <span class="cmdb-form-item">
                                <bk-selector
                                    setting-key="bk_obj_id"
                                    display-key="bk_obj_name"
                                    :list="staList"
                                    :selected.sync="newChart.static"
                                >
                                </bk-selector>
                            </span>
                        </div>
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('Operation["统计维度"]')}}：
                            </label>
                            <span class="cmdb-form-item">
                                <bk-selector
                                    setting-key="bk_property_id"
                                    display-key="bk_property_name"
                                    :list="filterList"
                                    :selected.sync="newChart.dim"
                                >
                                </bk-selector>
                            </span>
                        </div>
                        <div class="content-item">
                            <label class="label-text">
                                {{$t('Operation["图表展示"]')}}：
                            </label>
                            <label class="cmdb-form-radio cmdb-radio-small">
                                <input type="radio" name="present" value="pie" v-model="newChart.present">
                                <span class="cmdb-radio-text">{{$t('Operation["饼图"]')}}</span>
                            </label>
                            <label class="cmdb-form-radio cmdb-radio-small">
                                <input type="radio" name="present" value="Hist" v-model="newChart.present">
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
                                <bk-selector
                                    setting-key="repType"
                                    display-key="name"
                                    :list="seList.disList"
                                    :selected.sync="newChart.present"
                                >
                                </bk-selector>
                            </span>
                        </div>
                    </div>
                    <div class="content-item">
                        <label class="label-text">
                            {{$t('Operation["图表宽度"]')}}：
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="width" value="50" v-model="newChart.width">
                            <span class="cmdb-radio-text">50%</span>
                        </label>
                        <label class="cmdb-form-radio cmdb-radio-small">
                            <input type="radio" name="width" value="100" v-model="newChart.width">
                            <span class="cmdb-radio-text">100%</span>
                        </label>
                    </div>
                </div>
                <div class="footer">
                    <bk-button type="primary" @click="confirm">{{$t("Common['保存']")}}</bk-button>
                    <bk-button type="default" @click="cancel">{{$t("Common['取消']")}}</bk-button>
                </div>
            </div>
        </bk-dialog>
    </div>
</template>

<script>
    import Plotly from 'plotly.js'
    import { mapActions } from 'vuex'

    export default {
        name: 'index',
        data () {
            return {
                businessList: [],
                charData: [],
                disCharts: [],
                editModel: false,
                isShow: false,
                list: [
                    {
                        id: 1,
                        name: 'haha1'
                    },
                    {
                        id: 2,
                        name: 'haha2'
                    },
                    {
                        id: 3,
                        name: 'haha3'
                    }
                ],
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
                newChart: {
                    type: true,
                    name: '',
                    static: '',
                    dim: '',
                    present: '',
                    choose: '',
                    width: 50,
                    newType: 'host'
                },
                staList: [],
                demList: []
            }
        },
        computed: {
            filterList () {
                return this.demList.filter(item => {
                    return item.bk_property_type === 'enum'
                })
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Operation["统计报表"]'))
            this.getChartList()
        },
        mounted () {
            this.drawCharts()
        },
        methods: {
            ...mapActions('operationChart', [
                'getCountedCharts',
                'getCountedChartsData',
                'getStaticObj',
                'getStaticDimeObj'
            ]),
            async getChartList () {
                // const res = await this.getCountedCharts({})
                // this.charData = res.data.info
                // this.getChartData()
                // this.charData.forEach(item => {
                //    const data = this.getChartData(item.config_id).data
                //    item.chart_id = item.report_type + item.config_id
                // })
                this.charData = [{
                    'report_type': 'custom',
                    'name': '主机1',
                    'config_id': 1,
                    'chart_id': 'custom1',
                    'option': {
                        'bk_obj_id': 'host',
                        'chart_type': 'pie',
                        'field': 'bk_os_type'
                    },
                    'position': {
                        'width': '50',
                        'index': 3
                    },
                    'data': {
                        'windows': 10,
                        'linux': 20
                    }
                }, {
                    'report_type': 'custom',
                    'name': '主机222',
                    'config_id': 2,
                    'chart_id': 'custom1',
                    'option': {
                        'bk_obj_id': 'host',
                        'chart_type': 'pie',
                        'field': 'bk_os_type'
                    },
                    'position': {
                        'width': '50',
                        'index': 1
                    },
                    'data': {
                        'windows': 10,
                        'linux': 20
                    }
                }, {
                    'report_type': 'custom',
                    'name': '主机3333',
                    'config_id': 3,
                    'chart_id': 'custom1',
                    'option': {
                        'bk_obj_id': 'host',
                        'chart_type': 'pie',
                        'field': 'bk_os_type'
                    },
                    'position': {
                        'width': '100',
                        'index': 2
                    },
                    'data': {
                        'windows': 10,
                        'linux': 20
                    }
                }]
                this.disCharts = this.$tools.clone(this.charData)
            },
            async getChartData (id) {
                const res = await this.getCountedCharts({ id })
                return res
            },
            async getStaList () {
                this.staList = await this.getStaticObj({})
            },
            async getDemList (id) {
                this.demList = await this.getStaticDimeObj({
                    'bk_obj_id': id,
                    'bk_supplier_account': '0'
                })
            },
            drawCharts () {
                this.disCharts.forEach(item => {
                    const myDiv = document.getElementById(item.report_type + item.config_id)
                    const data = [{
                        values: [19, 26, 55],
                        labels: ['Residential', 'Non-Residential', 'Utility'],
                        type: 'pie'
                    }]
                    const layout = {
                        height: 400,
                        width: item.position.width === '50' ? 500 : 1000,
                        title: item.name
                    }
                    const options = {
                        displaylogo: false,
                        displayModeBar: false
                    }
                    Plotly.newPlot(myDiv, data, layout, options)
                })
            },
            editPage () {
                this.$bkInfo({
                    title: this.$tc('Cloud["确认保存更改?"]'),
                    confirmFn: () => {
                        this.editModel = !this.editModel
                        if (this.charData === this.disCharts) {
                            this.charData = this.$tools.clone(this.disCharts)
                            this.drawCharts()
                        }
                    }
                })
            },
<<<<<<< HEAD:src/ui/src/views/statistical-report/index.vue
            mutipChart () {
                const muDiv = document.getElementById('muDiv')
                const trace1 = {
                    x: ['giraffes', 'orangutans', 'monkeys'],
                    y: [20, 14, 23],
                    name: 'SF Zoo',
                    type: 'bar'
                }

                const trace2 = {
                    x: ['giraffes', 'orangutans', 'monkeys'],
                    y: [12, 18, 29],
                    name: 'LA Zoo',
                    type: 'bar'
=======
            editCancel () {
                this.disCharts = this.$tools.clone(this.charData)
                this.editModel = !this.editModel
                setTimeout(() => {
                    this.drawCharts()
                }, 100)
            },
            openNew (type) {
                this.newChart = {
                    type: true,
                    name: '',
                    static: '',
                    dim: '',
                    present: 'pie',
                    choose: '',
                    width: 50
                }
                this.isShow = !this.isShow
                this.newChart.newType = type
                if (type === 'host') {
                    this.seList.disList = this.$tools.clone(this.seList.host)
                    this.getDemList('host')
                } else {
                    this.seList.disList = this.$tools.clone(this.seList.inst)
                    this.getStaList()
>>>>>>> dcdf32d55... feature: operation-chart ui:src/ui/src/views/operation/index.vue
                }
            },
            moveChart (type, dire, key) {
                if (dire === 'up' && key !== 0) {
                    this.disCharts[key] = this.disCharts.splice(key - 1, 1, this.disCharts[key])[0]
                } else if (dire === 'down' && key !== this.disCharts.length - 1) {
                    this.disCharts[key + 1] = this.disCharts.splice(key, 1, this.disCharts[key + 1])[0]
                }
            },
            editCharts (type, data) {
                this.isShow = !this.isShow
                this.newChart = data
                if (type === 'host') this.seList.disList = this.$tools.clone(this.seList.host)
                else this.seList.disList = this.$tools.clone(this.seList.inst)
            },
            deleteChart (type, key) {
                this.disCharts.splice(key, 1)
            },
            confirm () {
                this.isShow = !this.isShow
            },
            cancel () {
                this.isShow = !this.isShow
            },
            selDem (data) {
                console.log(data)
                this.getDemList(data)
            }
        }
    }
</script>

<style scoped lang="scss">
    *{
        margin:0;
        padding:0
    }
    .operation-top{
        padding: 5px;
        margin: 5px;
        border-bottom: 2px solid #eee;

        .operation-title{
            display: block;
            font-size: 14px;
        }
        .operation-edit{
            display: block;
            text-align: right;
            font-size: 12px;
            color: $cmdbMainBtnColor;
            text-decoration: underline;
            cursor: pointer;
            span{
                margin-left: 10px;
            }
        }
    }
    .operate-menus{
        height: 110px;
        padding: 10px 5px;
        .menu-items{
            height: 100%;
            width: 25%;
            border: 1px solid #eee;
            float: left;
            padding: 15px 20px;
            box-shadow: 1px 2px 4px 0 rgba(51, 60, 72, 0.06);
            &:first-child{
                 border-right: none;
                .item-left{
                    background: $cmdbMainBtnColor;
                }
             }
            &:last-child{
                 border-left: none;
             }

             .item-left{
                 width: 60px;
                 height: 60px;
                 font-size: 32px;
                 text-align: center;
                 line-height: 60px;
                 float: left;
                 border-radius: 50%;
                 margin-right: 30px;
                 background-color: #f5e41f;
                 i{
                     color: white;
                     vertical-align: inherit;
                 }
             }
             .item-right{
                 height: 100%;
                 text-align: center;
                 width: calc(100% - 120px);
                 float: left;
                 span{
                     display: block;
                     text-align: left;
                     &:first-child{
                        font-weight: bold;
                        font-size: 20px;
                        margin-top: 5px;
                      }
                      &:last-child{
                        margin-top: 5px;
                        font-size: 12px;
                       }
                 }
             }
        }
    }
    .operate-menus-bottom{
        .menu-items{
            &:last-child{
                 border-left: 1px solid #eee;
                .item-left{
                    background: #f33c13;
                }
             }
        }
    }
    .operation-layout{
        padding: 10px 5px;
        display: inline-block;
        position: relative;
    }
    .operation-charts{
        border: 1px solid #eee;
        box-shadow: 1px 2px 4px 0 rgba(51, 60, 72, 0.06);
    }
    .charts-options{
        position: absolute;
        right: 15px;
        top: 10px;
        color: $cmdbMainBtnColor;
        i{
            cursor: pointer;
            font-weight: bold;
            margin-left: 3px;
        }
    }
    .operation-add{
        width: 200px;
        height: 120px;
        border: 1px solid #eee;
        font-size: 100px;
        margin-left: 5px;
        color: $cmdbMainBtnColor;
        box-shadow: 1px 2px 4px 0 rgba(51, 60, 72, 0.06);
        text-align: center;
        cursor: pointer;
        i{
            font-weight: bold;
            vertical-align: baseline;
        }
    }
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
                    vertical-align: top;
                }
                .cmdb-form-item {
                    display: inline-block;
                    margin-right: 10px;
                    width: 319px;
                    vertical-align: middle;
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
</style>
