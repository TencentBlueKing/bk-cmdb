<template>
    <div>
        <div class="operate-menus" v-if="!editModel">
            <div class="menu-items menu-items-blue" @click="goRouter('business')">
                <div class="item-left">
                    <span>{{ navData.biz }}</span>
                    <span>{{$t('Operation["业务总数"]')}}</span>
                </div>
                <div class="item-right item-right-left">
                    <i class="icon icon-cc-op-biz"></i>
                </div>
            </div>
            <div class="menu-items menu-items-white" @click="goRouter('resource')">
                <div class="item-left">
                    <span>{{ navData.host }}</span>
                    <span>{{$t('Operation["主机总数"]')}}</span>
                </div>
                <div class="item-right item-right-right">
                    <i class="icon icon-cc-op-host"></i>
                </div>
            </div>
            <div class="menu-items menu-items-blue" @click="goRouter('model')">
                <div class="item-left">
                    <span>{{ navData.model }}</span>
                    <span>{{$t('Operation["模型总数"]')}}</span>
                </div>
                <div class="item-right item-right-left">
                    <i class="menu-icon icon-cc-nav-model"></i>
                </div>
            </div>
            <div class="menu-items menu-items-white" @click="goRouter('index')">
                <div class="item-left">
                    <span>{{ navData.inst }}</span>
                    <span>{{$t('Operation["实例总数"]')}}<bk-tooltip :content="tooltip" :width="'100px'" :placement="'right'"><i class="menu-icon icon-cc-attribute"></i></bk-tooltip></span>
                </div>
                <div class="item-right item-right-right">
                    <i class="icon icon-cc-op-example"></i>
                </div>
            </div>
        </div>
        <div class="operation-top">
            <span class="operation-edit" v-if="!editModel" @click="editModel = true">
                <i class="icon icon-cc-edit"></i>{{$t('Operation["编辑图表"]')}}</span>
            <div class="exit-edit" v-if="editModel">
                <span>{{$t('Operation["所有更改已实时保存"]')}}</span>
                <span @click="editCancel()"> {{$t('Operation["退出编辑"]')}}</span>
            </div>
            <span class="operation-title">{{$t('Operation["主机统计"]')}}</span>
        </div>
        <div v-for="(item, key) in hostData.disList"
            :key="item.report_type + item.config_id"
            :style="{ width: item.width + '%' }"
            class="operation-layout">
            <div class="chart-child">
                <div class="chart-title">
                    <span>{{item.name}}</span>
                </div>
                <div class="operation-charts" :id="item.report_type + item.config_id"></div>
                <div v-if="item.noData" class="null-data">
                    <span>暂无数据</span>
                </div>
                <div class="charts-options" v-if="editModel">
                    <i class="bk-icon icon-arrows-up icon-weight" :class="{ 'icon-disable': key === 0 }"
                        @click="moveChart('host', 'up', key, hostData.disList)"></i>
                    <i class="bk-icon icon-arrows-down icon-weight" :class="{ 'icon-disable': key === hostData.disList.length - 1 }"
                        @click="moveChart('host', 'down', key, hostData.disList)"></i>
                    <i class="icon icon-cc-edit-shape"
                        @click="openNew('edit', 'host', item, key)"></i>
                    <i class="icon icon-cc-tips-close"
                        @click="deleteChart('host', key, hostData.disList, item.config_id)"></i>
                </div>
            </div>
        </div>
        <div class="add-parent">
            <div class="operation-add" v-if="editModel" @click="openNew('add', 'host')">
                <span><i class="bk-icon icon-plus"></i>{{$t('Operation["新增图表"]')}}</span>
            </div>
        </div>
        <div class="operation-top">
            <span class="operation-title">{{$t('Operation["实例统计"]')}}</span>
        </div>
        <div v-for="(item, key) in instData.disList"
            :key="item.report_type + item.config_id"
            :style="{ width: item.width + '%' }"
            class="operation-layout">
            <div class="chart-child">
                <div class="chart-title">
                    <span>{{item.name}}</span>
                </div>
                <div class="operation-charts" :id="item.report_type + item.config_id"></div>
                <div v-if="item.noData" class="null-data">
                    <span>暂无数据</span>
                </div>
                <div class="charts-options" v-if="editModel">
                    <i class="bk-icon icon-arrows-up icon-weight" :class="{ 'icon-disable': key === 0 }"
                        @click="moveChart('inst', 'up', key, instData.disList)"></i>
                    <i class="bk-icon icon-arrows-down icon-weight" :class="{ 'icon-disable': key === instData.disList.length - 1 }"
                        @click="moveChart('inst', 'down', key, instData.disList)"></i>
                    <i class="icon icon-cc-edit-shape" @click="openNew('edit', 'inst', item, key)"></i>
                    <i class="icon icon-cc-tips-close"
                        @click="deleteChart('inst', key, instData.disList, item.config_id)"></i>
                </div>
            </div>
        </div>
        <div class="add-parent">
            <div class="operation-add" v-if="editModel" @click="openNew('add', 'inst')">
                <span><i class="bk-icon icon-plus"></i>{{$t('Operation["新增图表"]')}}</span>
            </div>
        </div>
        <v-detail v-if="isShow"
            :open-type="editType.openType"
            :host-type="editType.hostType"
            :chart-data="newChart"
            @transData="saveData"
            @closeChart="cancelData">
        </v-detail>
    </div>
</template>

<script>
    import Plotly from 'plotly.js'
    import { mapActions } from 'vuex'
    import vDetail from './chart-detail'
    export default {
        name: 'index',
        components: {
            vDetail
        },
        data () {
            return {
                tooltip: '不包含业务、主机模型及实例',
                editModel: false,
                isShow: false,
                newChart: {},
                editType: {
                    openType: 'add',
                    hostType: 'host',
                    index: ''
                },
                hostData: {
                    disList: []
                },
                instData: {
                    disList: []
                },
                navData: {
                    biz: '',
                    host: '',
                    module: '',
                    inst: '',
                    model: ''
                }
            }
        },
        created () {
            this.$store.commit('setHeaderTitle', this.$t('Operation["统计报表"]'))
            this.getChartList()
        },
        methods: {
            ...mapActions('operationChart', [
                'getCountedCharts',
                'getCountedChartsData',
                'deleteOperationChart',
                'updateChartPosition'
            ]),
            async getChartList () {
                const res = await this.getCountedCharts({})
                this.hostData.disList = res.info.host
                this.instData.disList = res.info.inst
                res.info.nav.forEach(item => {
                    this.getNavData(item, 'nav')
                })
                this.hostData.disList.forEach((item) => {
                    this.getNavData(item, 'host')
                })
                this.instData.disList.forEach((item) => {
                    this.getNavData(item, 'inst')
                })
            },
            async getNavData (item, type) {
                const res = await this.getCountedChartsData({
                    params: {
                        config_id: item.config_id
                    },
                    config: {
                        globalError: false
                    }
                })
                if (type === 'nav') {
                    res.forEach(items => {
                        this.navData[items.id] = items.count
                    })
                } else {
                    if (JSON.stringify(res) === '{}' || JSON.stringify(res) === '[]' || res === null) {
                        item.hasData = false
                        item.data = {
                            data: [{
                                labels: [],
                                type: item.chart_type,
                                values: [],
                                x: [],
                                y: []
                            }]
                        }
                    } else {
                        item.hasData = true
                        item.data = await this.dataDeal(item, res)
                    }
                    this.drawCharts(item)
                }
            },
            dataDeal (data, res) {
                const returnData = {
                    data: []
                }
                if (res && Array.isArray(res)) {
                    const content = {
                        labels: [],
                        values: [],
                        x: [],
                        y: [],
                        type: data.chart_type,
                        hoverlabel: {
                            bgcolor: '#fff',
                            bordercolor: '#DCDEE5',
                            font: {
                                color: '#63656E',
                                size: '12'
                            }
                        },
                        marker: { size: 16, color: [] },
                        mode: 'markers',
                        textinfo: 'none',
                        hole: 0.7,
                        pull: 0.01,
                        width: 0.1
                    }
                    res.forEach(item => {
                        if (data.chart_type === 'pie') {
                            content.labels.push(item.id)
                            content.values.push(item.count)
                        } else {
                            const color = '#3A84FF'
                            content.marker.color.push(color)
                            content.x.push(item.id)
                            content.y.push(item.count)
                        }
                    })
                    if (data.chart_type !== 'pie') content.width = (0.015 * data.x_axis_count * (data.width === '50' ? 2 : 1))
                    returnData.data.push(content)
                } else if (res && !Array.isArray(res)) {
                    const barModel = {
                        chart: {
                            delete: {
                                name: 'delete',
                                type: 'bar',
                                y: [],
                                color: '#3A84FF',
                                marker: { size: 16, color: [] }
                            },
                            update: {
                                name: 'update',
                                type: 'bar',
                                y: [],
                                color: '#38C1E2',
                                marker: { size: 16, color: [] }
                            },
                            create: {
                                name: 'create',
                                type: 'bar',
                                y: [],
                                color: '#59D178',
                                marker: { size: 16, color: [] }
                            }
                        },
                        x: [],
                        isX: false
                    }
                    for (const item in res) {
                        const content = {
                            name: item,
                            x: [],
                            y: [],
                            hoverlabel: {
                                bgcolor: '#fff',
                                bordercolor: '#DCDEE5',
                                font: {
                                    color: '#63656E',
                                    size: '12'
                                }
                            },
                            textinfo: 'none'
                        }
                        if (Array.isArray(res[item])) {
                            barModel.isX = false
                            res[item].forEach(child => {
                                const time = this.$tools.formatTime(child.id, 'YYYY-MM-DD')
                                content.x.push(time)
                                content.y.push(child.count)
                                content.mode = 'line'
                            })
                            returnData.data.push(content)
                        } else {
                            barModel.x.push(item)
                            barModel.isX = true
                            for (const chartItem in barModel.chart) {
                                barModel.chart[chartItem].y.push(res[item][chartItem])
                                barModel.chart[chartItem].marker.color.push(barModel.chart[chartItem].color)
                            }
                        }
                    }
                    if (barModel.isX) {
                        for (const chartItem in barModel.chart) {
                            barModel.chart[chartItem].x = barModel.x
                            returnData.data.push(barModel.chart[chartItem])
                        }
                    }
                }
                return returnData
            },
            drawCharts (item) {
                const myDiv = item.report_type + item.config_id
                const data = item.data.data
                if (!item.hasData) this.$set(item, 'noData', true)
                const layout = {
                    title: ``,
                    barmode: 'stack',
                    height: 250,
                    margin: {
                        l: 50,
                        r: 50,
                        t: 50,
                        b: 50
                    },
                    yaxis: {
                        fixedrange: true,
                        rangemode: 'tozero'
                    },
                    xaxis: {
                        type: 'category',
                        range: [
                            -0.5,
                            item.x_axis_count - 0.5
                        ],
                        autorange: false,
                        tickformat: '%Y-%m-%d'
                    },
                    bargap: 0.1,
                    bargroupgap: 1.0044 - (0.0401 * item.x_axis_count),
                    dragmode: 'pan',
                    colorway: ['#3A84FF', '#A3C5FD', '#59D178', '#94F5A4', '#38C1E2', '#A1EDFF', '#4159F9', '#7888F0', '#EFAF4B', '#FEDB89', '#FF5656', '#FD9C9C']
                }
                console.log(layout)
                if (item.report_type !== 'model_inst_change_chart') layout.hovermode = 'closest'
                const options = {
                    displaylogo: false,
                    displayModeBar: false,
                    responsive: true
                }
                if (this.editType.openType === 'edit') {
                    Plotly.purge(myDiv)
                    setTimeout(() => {
                        Plotly.newPlot(myDiv, data, layout, options)
                        if (item.data.data[0].mode !== 'line') this.hoverConfig(document.getElementById(myDiv), myDiv)
                    }, 100)
                } else {
                    Plotly.newPlot(myDiv, data, layout, options)
                    if (item.data.data[0].mode !== 'line') this.hoverConfig(document.getElementById(myDiv), myDiv)
                }
            },
            hoverConfig (myPlot, myDiv) {
                myPlot.on('plotly_hover', (data) => {
                    const colors = ['#A3C5FD', '#B0E7F4', '#BDEDC9']
                    this.reStyle(myDiv, data, colors)
                })
                myPlot.on('plotly_unhover', (data) => {
                    const colors = ['#3A84FF', '#38C1E2', '#59D178']
                    this.reStyle(myDiv, data, colors)
                })
            },
            reStyle (myDiv, data, color) {
                let pn = ''
                let tn = ''
                let colors = []
                for (let i = 0; i < data.points.length; i++) {
                    pn = data.points[i].pointNumber
                    tn = data.points[i].curveNumber
                    colors = data.points[i].data.marker.color
                }
                if (data.points.length === 1) {
                    colors[pn] = color[0]
                    const update = {
                        'marker': {
                            color: colors
                        }
                    }
                    Plotly.restyle(myDiv, update, [tn])
                } else {
                    data.points[0].data.marker.color[pn] = color[0]
                    data.points[1].data.marker.color[pn] = color[1]
                    data.points[2].data.marker.color[pn] = color[2]
                    const update = {
                        data: {
                            maker: {
                                color: [data.points[0].data.marker.color, data.points[1].data.marker.color, data.points[2].data.marker.color]
                            }
                        }
                    }
                    Plotly.restyle(myDiv, update, [0, 2])
                }
            },
            editCancel () {
                this.editModel = false
            },
            moveChart (type, dire, key, list) {
                if (dire === 'up' && key !== 0) {
                    list[key] = list.splice(key - 1, 1, list[key])[0]
                } else if (dire === 'down' && key !== list.length - 1) {
                    list[key + 1] = list.splice(key, 1, list[key + 1])[0]
                }
                if (type === 'host') this.hostData.disList = list
                else this.instData.disList = list
                this.updatePosition()
            },
            deleteChart (type, key, list, id) {
                this.$bkInfo({
                    title: this.$tc('Operation["是否确认删除"]'),
                    confirmFn: () => {
                        list.splice(key, 1)
                        if (type === 'host') this.hostData.disList = list
                        else this.instData.disList = list
                        this.deleteOperationChart({
                            id: id
                        })
                        this.updatePosition()
                    }
                })
            },
            async openNew (type, host, data, key) {
                this.editType.hostType = host
                this.editType.key = key
                if (type === 'edit') {
                    this.newChart = this.$tools.clone(data)
                } else {
                    this.newChart = {
                        report_type: 'custom',
                        name: '',
                        config_id: null,
                        bk_obj_id: host,
                        chart_type: 'pie',
                        field: '',
                        width: '50',
                        x_axis_count: 10
                    }
                }
                this.editType.openType = type
                this.isShow = true
            },
            async saveData (data) {
                console.log(data)
                let editList = []
                if (this.editType.hostType === 'host') editList = this.hostData.disList
                else editList = this.instData.disList
                if (this.editType.openType === 'add') {
                    editList.push(data)
                } else {
                    editList[this.editType.key] = data
                }
                this.isShow = false
                await this.getNavData(data, this.editType.hostType)
                this.updatePosition()
                this.newChart = {}
            },
            cancelData () {
                this.isShow = false
                this.newChart = {}
            },
            updatePosition () {
                const data = {
                    'host': [],
                    'inst': []
                }
                this.hostData.disList.forEach(item => {
                    data.host.push(item.config_id)
                })
                this.instData.disList.forEach(item => {
                    data.inst.push(item.config_id)
                })
                this.updateChartPosition({
                    params: {
                        position: data
                    }
                })
            },
            goRouter (route) {
                this.$router.push(route)
            }
        }
    }
</script>

<style scoped lang="scss">
    *{
        margin: 0;
        padding: 0
    }
    .operation-top{
        margin-top: 33px;
        margin-left: 10px;
        .operation-title{
            display: inline-block;
            width: 64px;
            height: 21px;
            font-size: 16px;
            font-weight: bold;
            color: rgba(49,50,56,1);
            line-height: 21px;
        }
        .exit-edit{
            width: calc(100% - 60px);
            height: 46px;
            text-align: center;
            background: rgba(225,236,255,1);
            border: 1px solid rgba(220,222,229,1);
            line-height: 46px;
            position: fixed;
            z-index: 1;
            top: 60px;
            left: 60px;
            span{
                &:first-child{
                    width: 126px;
                    height: 19px;
                    font-size: 14px;
                    color: rgba(49,50,56,1);
                    line-height: 19px;
                 }
                &:last-child{
                    display: inline-block;
                    position: absolute;
                    top: 7px;
                    right: 15px;
                    width: 86px;
                    height: 32px;
                    background: rgba(58,132,255,1);
                    border-radius: 2px;
                    font-size: 14px;
                    font-weight: 400;
                    color: rgba(255,255,255,1);
                    line-height: 32px;
                 }
            }
        }
        .operation-edit{
            display: inline-block;
            cursor: pointer;
            float: right;
            margin-right: 10px;
            height: 19px;
            font-size: 14px;
            color: rgba(58,132,255,1);
            line-height: 19px;
            i{
                margin-right: 5px;
            }
        }
    }
    .operate-menus{
        height: 105px;
        display: flex;
        width: 100%;
        .menu-items-blue{
            background: linear-gradient(165deg,rgba(108,186,255,1) 0%,rgba(58,132,255,1) 100%);
        }
        .menu-items-white{
            background: linear-gradient(165deg,rgba(77,212,240,1) 0%,rgba(27,167,207,1) 100%);
        }
        .menu-items{
            display: flex;
            align-items: center;
            width: 25%;
            height: 100%;
            margin: 10px;
            float: left;
            box-shadow: 0 2px 4px 0 rgba(220,222,229,1);
            border-radius: 4px;
            position: relative;
            .item-left{
                width: 60px;
                height: 45px;
                float: left;
                margin-left: 26px;
                line-height: 45px;
                span{
                    display: block;
                    text-align: left;
                    &:first-child{
                         width: 29px;
                         height: 29px;
                         font-size: 24px;
                         font-weight: bold;
                         color: rgba(255,255,255,1);
                         line-height: 29px;
                     }
                    &:last-child{
                         height: 16px;
                         width: 70px;
                         font-size: 12px;
                         font-weight: bold;
                         color: rgba(255,255,255,1);
                         line-height: 16px;
                         i{
                             display: inline-block;
                             font-size: 16px;
                             margin-left: 5px;
                         }
                     }
                }
            }
            .item-right-left{
                background-color: #377AE9;
             }
            .item-right-right{
                background-color: #19A2CA;
             }
            .item-right{
                 display: flex;
                 justify-content: center;
                 align-items: center;
                 height: 60px;
                 width: 60px;
                 position: absolute;
                 right: 15%;
                 border-radius: 50%;
               i{
                   color: white;
                   font-size: 34px;
               }
            }
        }
    }
    .operation-layout{
        display: inline-block;
        position: relative;
        padding: 10px;
        .chart-child{
             background: rgba(255,255,255,1);
             border-radius: 2px;
             border: 1px solid rgba(220,222,229,1);
             position: relative;
             width: 100%;
             .operation-charts{
                width: 100%;
             }
        }
        .chart-title{
            display: block;
            height: 50px;
            border-bottom: 1px solid rgba(240,241,245,1);
            span{
            line-height: 50px;
            font-size: 14px;
            font-weight: bold;
            color: rgba(99,101,110,1);
            margin-left: 21px;
            }
        }
        .charts-options{
            position: absolute;
            right: 15px;
            top: 10px;
            color: $cmdbMainBtnColor;
            i{
                box-sizing: border-box;
                cursor: pointer;
                margin-left: 3px;
                border: 1px solid rgba(99,101,110,0);
                padding: 1px;
                &:hover{
                     border: 1px dashed #3A84FF;
                 }
            }
        }
        .icon-disable{
            color: grey;
        }
        .icon-weight{
            padding: 0!important;
            font-size: 18px;
            font-weight: bold;
        }
        .null-data {
            position: absolute;
            top: 50%;
            width: 100%;
            text-align: center;
            height: 19px;
            span {
                margin: 0 auto;
                display: block;
                font-size: 14px;
                color: rgba(151,155,165,1);
                line-height: 19px;
            }
        }
    }
    .add-parent{
        width: 100%;
        padding: 10px;
        .operation-add{
            width: 100%;
            height: 49px;
            line-height: 49px;
            border: 1px dashed rgba(58,132,255,1);
            font-size: 15px;
            color: rgba(58,132,255,1);
            box-shadow: 1px 2px 4px 0 rgba(51, 60, 72, 0.06);
            text-align: center;
            cursor: pointer;
            &:hover{
                background: rgba(225,236,255,1);
             }
            span{
                display: inline-block;
                i{
                    font-weight: bold;
                    vertical-align: baseline;
                    margin-right: 10px;
                }
            }
        }
    }
</style>
