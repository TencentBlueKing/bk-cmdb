/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template>
    <div class="topo-wrapper model" :class="{'no-edit': curClassify === 'bk_host_manage'}" id="topo-wrapper">
        <div id="topo" class="topo"></div>
        <button v-if="curClassify !== 'bk_host_manage'" class="bk-button vis-button vis-setup" @click="editModel" :title="$t('Common[\'编辑\']')">
            <i class="icon icon-cc-edit"></i>
        </button>
        <bk-button type="primary" class="bk-button vis-button vis-create" v-if="addModelAvailable && curClassify && modelNodes.length" @click="createModel">
            <span class="vis-button-text">{{$t('ModelManagement["新增"]')}}</span>
        </bk-button>
        <button class="bk-button vis-button vis-enable" v-if="addModelAvailable" @click="isShowDisableList = true">
            <i class="bk-icon icon-minus-circle-shape"></i>
            <span class="vis-button-text">{{disableModelList.length}}</span>
        </button>
        <transition name="topo-disable-list">
            <div class="topo-disable" v-show="isShowDisableList">
                <label class="disable-title">
                    <span>{{$t('ModelManagement["已停用模型"]')}}</span>
                    <i class="bk-icon icon-arrows-right" @click="isShowDisableList = false"></i>
                </label>
                <ul class="disable-list" ref="disableList">
                    <li class="disable-item" v-for="(model, index) in disableModelList" :key="index">
                        <a class="disable-item-link" href="javascript:void(0)" @click="disableNodeClick(model)">{{model['bk_obj_name']}}</a>
                    </li>
                </ul>
            </div>
        </transition>
        <bk-button type="danger" class="bk-button vis-button vis-del" :title="$t('ModelManagement[\'删除\']')" v-if="!isInnerType" @click="deleteClass">
            <i class="icon icon-cc-del"></i>
        </bk-button>
    </div>
</template>

<script>
    import vis from 'vis'
    export default {
        props: {
            isTopoLoading: {
                type: Boolean
            },
            isInnerType: {
                type: Boolean,
                default: false
            },
            addModelAvailable: {
                type: Boolean,
                required: true
            },
            disableModelList: {
                type: Array,
                default () {
                    return []
                }
            },
            /*
                当前模型分类ID
            */
            curClassify: '',
            /*
                节点之间的关联关系
            */
            topo: {
                type: Object,
                default: () => {
                    return {
                        nodes: [],
                        edges: []
                    }
                }
            },
            /*
                所有模型节点
            */
            modelNodes: {
                type: Array,
                default: () => {
                    return []
                }
            },
            /*
                true: 已启用模型   false: 未启用模型
            */
            isPaused: {
                default: false
            },
            /*
                是否点击的切换分组
                true: 点击的切换分组 false: 点击的已启用/未启用
            */
            isChangeClassify: false
        },
        data () {
            return {
                nodes: [],
                edges: [],
                otherNodes: [],         // 其他模型分类下的节点
                options: {
                    interaction: {
                        navigationButtons: true,
                        hover: true
                    },
                    edges: {
                        color: {
                            color: '#c3cdd7',
                            highlight: '#3c96ff',
                            hover: '#3c96ff'
                        },
                        smooth: {           // 线的动画
                            type: 'curvedCW',
                            roundness: 0
                        }
                    },
                    nodes: {
                        shadow: {
                            x: 2,
                            y: 2,
                            color: 'rgba(0, 0, 0, .1)'
                        },
                        color: {
                            border: '#eee'
                        },
                        shapeProperties: {
                            borderDashes: [5, 5]
                        }
                    }
                },
                container: '',
                network: '',
                isShowDisableList: false
            }
        },
        computed: {
            graph_data () {
                return {
                    nodes: new vis.DataSet(this.nodes),
                    edges: new vis.DataSet(this.edges)
                }
            }
        },
        watch: {
            topo: {
                handler: function (val) {
                    if (this.isChangeClassify) {
                        this.formatTopo(val)
                    }
                },
                deep: true
            },
            /*
                已启用模型/未启用模型 切换时
            */
            isPaused () {
                this.formatTopo()
            }
        },
        methods: {
            /*
                删除分组
            */
            deleteClass () {
                this.$emit('deleteClass')
            },
            createModel () {
                this.$emit('createModel')
            },
            editModel () {
                this.$emit('editModel')
            },
            disableNodeClick (model) {
                this.$emit('nodeClick', model, true)
            },
            /*
                把十六位色值转换为rgb
                return {
                    r: '111',
                    g: '222',
                    b: '123'
                }
            */
            parseColor (color) {
                let r = ''
                let g = ''
                let b = ''
                let len = color.length
                // 非简写模式 #123456
                if (len === 7) {
                    r = parseInt(color.slice(1, 3), 16)
                    g = parseInt(color.slice(3, 5), 16)
                    b = parseInt(color.slice(5, 7), 16)
                } else if (len === 4) {   // 简写模式 #6cf
                    r = parseInt(color.charAt(1) + color.charAt(1), 16)
                    g = parseInt(color.charAt(2) + color.charAt(2), 16)
                    b = parseInt(color.charAt(3) + color.charAt(3), 16)
                }
                return {
                    r: r,
                    g: g,
                    b: b
                }
            },
            /*
                通过图标类名获取icon路径
                class: icon-xxx-xxx
                return xxx-xxx
            */
            getIconByClass (iconClass) {
                return iconClass.substr(5)
            },
            /*
                获取图片base64
                img: 图片地址
                rgb: {
                    r: r,
                    g: g,
                    b: b
                }
            */
            getBase64Img (img, rgb) {
                let canvas = document.createElement('canvas')
                canvas.width = img.width
                canvas.height = img.height
                let ctx = canvas.getContext('2d')
                
                ctx.drawImage(img, 0, 0, img.width, img.height)

                let dataL = ctx.getImageData(0, 0, canvas.width, canvas.height)
                
                // 设置颜色
                for (let i = 0; i < dataL.data.length; i += 4) {
                    dataL.data[i] = rgb.r
                    dataL.data[i + 1] = rgb.g
                    dataL.data[i + 2] = rgb.b
                }
                ctx.putImageData(dataL, 0, 0)

                let ext = img.src.substring(img.src.lastIndexOf('.') + 1).toLowerCase()
                let dataURL = canvas.toDataURL('image/' + ext)
                return dataURL
            },
            /*
                格式化拓扑图
            */
            formatTopo (topo) {
                if (topo === undefined) {
                    topo = this.topo
                }
                // 关联关系
                this.edges = []
                topo.edges.map(edge => {
                    this.edges.push({
                        dashes: !(this.curClassify === edge.from['bk_classification_id'] && this.curClassify === edge.to['bk_classification_id']),
                        label: edge['label'],
                        arrows: edge['arrows'],
                        to: edge['to']['bk_obj_id'],
                        from: edge['from']['bk_obj_id']
                    })
                })
                this.nodes = []
                let status = 0
                // 当前状态下的节点数量  已启用/未启用
                let curTypeNum = 0
                topo.nodes.map(node => {
                    if (node['bk_ispaused'] === this.isPaused) {
                        curTypeNum++
                        let svgForUser = {
                            bgColor: '#fff',
                            iconColor: '#868b97',
                            color: '#868b97'
                        }
                        let svgForPre = {
                            bgColor: '#fff',
                            iconColor: '#6894c8',
                            color: '#6894c8'
                        }
                        let svgForOther = {
                            bgColor: '#fff',
                            iconColor: '#c3cdd7',
                            color: '#c3cdd7'
                        }
                        let svgColor = {}
                        if (node['ispre']) {
                            svgColor = svgForPre
                        } else {
                            svgColor = svgForUser
                        }
                        if (!this.isBelongtoCurclassify(node['bk_obj_id']) || node['bk_obj_id'] === 'plat') {
                            svgColor = svgForOther
                        }
                        // 没有图标的话就设置一个默认图标
                        if (!node.hasOwnProperty('bk_obj_icon') || node['bk_obj_icon'] === '') {
                            node['bk_obj_icon'] = 'icon-cc-business'
                        }

                        let img = `./static/svg/${this.getIconByClass(node['bk_obj_icon'])}.svg`
                        let image = new Image()
                        image.onload = () => {
                            var base64 = this.getBase64Img(image, this.parseColor(svgColor.iconColor))
                            let svg = `<svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                            <circle cx="50" cy="50" r="49" fill="${svgColor.bgColor}"/>
                                <svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, 0)" viewBox="0 0 18 18" x="35" y="-12" fill="${svgColor.iconColor}" width="35" >
                                <image width="15" xlink:href="${base64}"></image>
                                </svg>
                                <foreignObject x="0" y="58" width="100%" height="100%">
                                    <div xmlns="http://www.w3.org/1999/xhtml" style="font-size:14px">
                                        <div style="color:${svgColor.color};text-align: center;width: 60px;overflow:hidden;white-space:nowrap;text-overflow:ellipsis;margin:0 auto">${node['bk_obj_name']}</div>
                                    </div>
                                </foreignObject>
                            </svg>`
                            this.initImage(svg, node)
                            status++
                        }
                        image.onerror = () => {
                            let svg = `<svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, .1)" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                            <circle cx="50" cy="50" r="49" fill="${svgColor.bgColor}"/>
                                <svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, 0)" viewBox="0 0 18 18" x="35" y="-12" fill="${svgColor.iconColor}" width="35" >
                                </svg>
                                <foreignObject x="0" y="58" width="100%" height="100%">
                                    <div xmlns="http://www.w3.org/1999/xhtml" style="font-size:14px">
                                        <div style="color:${svgColor.color};text-align: center;width: 50px;overflow:hidden;white-space:nowrap;text-overflow:ellipsis;margin:0 auto">${node['bk_obj_name']}</div>
                                    </div>
                                </foreignObject>
                            </svg>`
                            this.initImage(svg, node)
                            status++
                        }
                        image.src = img
                    }
                })
                let timer = setInterval(() => {
                    if (status === curTypeNum) {
                        clearInterval(timer)
                        this.init()
                    }
                }, 200)
            },
            initImage (svg, node) {
                let url = `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
                let temp = {
                    id: node['bk_obj_id'],
                    size: 55,
                    physics: false,
                    image: url,
                    shape: 'image'
                }
                if (node.hasOwnProperty('position') && node['position'] !== '') {
                    let position = JSON.parse(node['position'])
                    if (position.hasOwnProperty(this.curClassify)) {
                        temp.x = position[this.curClassify].x
                        temp.y = position[this.curClassify].y
                    }
                }
                this.nodes.push(temp)
            },
            /*
                获取当前模型item
            */
            getModelById (id) {
                for (let i = 0; i < this.modelNodes.length; i++) {
                    if (this.modelNodes[i]['bk_obj_id'] === id) {
                        return this.modelNodes[i]
                    }
                }
            },
            /*
                保存位置信息
                node 拖拽的节点的ObjId
            */
            savePosition (node) {
                let position = window.network.getPositions()
                let params = {
                    bk_classification_id: '',
                    bk_obj_name: '',
                    bk_supplier_account: '',
                    position: ''
                }
                let id = ''
                let index = ''
                let pos = {}
                for (let i = 0; i < this.modelNodes.length; i++) {
                    if (this.modelNodes[i]['bk_obj_id'] === node) {
                        index = i
                        params['bk_classification_id'] = this.modelNodes[i]['bk_classification_id']
                        params['bk_obj_name'] = this.modelNodes[i]['bk_obj_name']
                        params['bk_supplier_account'] = this.modelNodes[i]['bk_supplier_account']
                        if (this.modelNodes[i].hasOwnProperty('position') && this.modelNodes[i]['position'] !== '') {
                            pos = JSON.parse(this.modelNodes[i]['position'])
                        }
                        pos[this.curClassify] = position[node]
                        params.position = JSON.stringify(pos)
                        id = this.modelNodes[i]['id']
                        break
                    }
                }
                this.$axios.put(`object/${id}`, params).then(res => {
                    if (res.result) {
                        this.$set(this.modelNodes[index], 'position', params['position'])
                    } else {
                        this.$alertMsg(this.$t('["更新位置信息失败"]'))
                    }
                })
            },
            /*
                根据ObjId判断是否属于当前分类
            */
            isBelongtoCurclassify (objId) {
                for (let i = 0; i < this.modelNodes.length; i++) {
                    if (this.modelNodes[i]['bk_obj_id'] === objId) {
                        return this.modelNodes[i]['bk_classification_id'] === this.curClassify
                    }
                }
            },
            /*
                初始化拓扑图
            */
            init () {
                let network = new vis.Network(this.container, this.graph_data, this.options)
                window.network = network
                let networkCanvas = this.container.getElementsByTagName('canvas')[0]

                // 设置按钮title
                let visBtnSetting = [{
                    'title': this.$t('ModelManagement["放大"]'),
                    'icon': 'icon-plus'
                }, {
                    'title': this.$t('ModelManagement["缩小"]'),
                    'icon': 'icon-minus'
                }, {
                    'title': this.$t('ModelManagement["还原"]'),
                    'icon': 'icon-full-screen'
                }]
                document.querySelectorAll('.vis-zoomIn,.vis-zoomOut,.vis-zoomExtends').forEach((visBtn, index) => {
                    visBtn.setAttribute('title', visBtnSetting[index]['title'])
                    visBtn.classList.add('bk-icon')
                    visBtn.classList.add(visBtnSetting[index]['icon'])
                })

                // 自适应大小
                let arr = []
                this.nodes.map(val => {
                    arr.push(val.id)
                })
                window.network.fit({nodes: arr})
                window.network.moveTo({scale: 0.8})

                // 绑定事件
                let self = this
                network.on('hoverNode', () => {
                    networkCanvas.style.cursor = 'pointer'
                })
                network.on('blurNode', () => {
                    networkCanvas.style.cursor = 'default'
                })
                network.on('click', (params) => {
                    // 点击了具体某个节点
                    if (params.nodes.length) {
                        let id = params.nodes[0]
                        if (this.isBelongtoCurclassify(id) && id !== 'plat') {
                            self.$emit('nodeClick', self.getModelById(id))
                        }
                    }
                })
                network.on('dragEnd', (params) => {
                    // 拖拽了具体某个节点
                    if (params.nodes.length) {
                        self.$emit('updateIsChangeClassify', false)
                        self.savePosition(params.nodes[0])
                    }
                })
                this.$emit('update:isTopoLoading', false)
            }
        },
        mounted () {
            this.container = document.getElementById('topo')
            this.$refs.disableList.style.maxHeight = `${Math.ceil(document.body.getBoundingClientRect().height * 2 / 3)}px`
        }
    }
</script>

<style lang="scss" scoped>
    .topo-wrapper{
        width: 100%;
        height: calc(100% - 4px);
        .topo{
            height: 100%;
        }
    }
</style>

<style lang="scss">
    .model#topo-wrapper{
        .vis-network{
            overflow: visible !important;
            outline: none;
        }
        .vis-button{
            position: absolute;
            width: 30px;
            height: 30px;
            line-height: 30px;
            top: 9px;
            padding: 0;
            cursor: pointer;
            border-radius: 50%;
            box-shadow: 0px 1px 5px 0px rgba(12, 34, 59, 0.2);
            border: none;
            text-align: center;
            // font-size: 0;
            &:not(.vis-create){
                background-color: #fff;
            }
            [class^="icon-cc-"],
            .vis-button-text{
                height: 100%;
                line-height: 30px;
                display: inline-block;
                vertical-align: middle;
                font-size: 12px;
                margin: 0 3px;
            }
            [class^="icon-cc-"],
            &.bk-icon{
                font-weight: bold;
                font-size: 14px;
            }
            &:not(.vis-create):hover{
                color: #6eb1ff;
            }
            &:not(.vis-create):active{
                color: #3188ed;
            }
        }
        .vis-up,.vis-down,.vis-left,.vis-right{
            display: none;
        }
        &.no-edit {
            .vis-zoomIn{
                left: 54px;
            }
            .vis-zoomOut{
                left: 92px;
            }
            .vis-zoomExtends{
                left: 15px;
            }
        }
        .vis-zoomIn{
            left: 92px;
        }
        .vis-zoomOut{
            left: 132px;
        }
        .vis-zoomExtends{
            left: 54px;
        }
        .vis-setup{
            left: 15px;
            .icon{
                position: relative;
                top: -1px;
                font-weight: normal;
            }
        }
        .vis-create{
            left: 170px;
            font-size: 18px;
            border-radius: 14px;
            width: 60px;
            font-size: 0;
            [class^="icon-cc-"]{
                position: relative;
                top: 1px;
            }
        }
        .vis-enable{
            width: auto;
            right: 49px;
            padding: 0 10px;
            border-radius: 15px;
            font-weight: normal;
            text-align: left;
            font-size: 0;
            .bk-icon{
                font-weight: normal;
                vertical-align: middle;
                font-size: 14px;
            }
        }
        .vis-del{
            right: 9px;
            color: #737987;
            .icon{
                font-weight: normal;
            }
            &:hover{
                color: #ef4c4c;
            }
        }
    }
    .topo-disable{
        position: absolute;
        right: 49px;
        top: 9px;
        width: 169px;
        background-color: #ffffff;
        box-shadow: 0px 1px 5px 0px 
            rgba(12, 34, 59, 0.2);
        border-radius: 2px;
        line-height: 40px;
        font-size: 14px;
        transform-origin: right top;
        transition: transform .3s ease-in-out;
        .disable-title{
            display: block;
            background-color: #fafbfd;
            padding: 0 0 0 22px;
            position: relative;
            .bk-icon{
                position: absolute;
                top: 12px;
                right: 12px;
                display: inline-block;
                font-size: 16px;
                transform: rotate(-45deg);
                cursor: pointer;
            }
        }
        .disable-list{
            padding: 6px;
            overflow: auto;
            min-height: 52px;
            @include scrollbar;
            .disable-item{
                .disable-item-link{
                    display: block;
                    padding: 0 0 0 16px;
                    color: inherit;
                    text-decoration: none;
                    transition: unset !important;
                    @include ellipsis;
                    &:hover{
                        background-color: #f1f7ff;
                        color: #3c96ff;
                    }
                    &:active{
                        background-color: #e2efff;
                        color: #3c96ff
                    }
                }
            }
        }
    }
    .topo-disable-list-enter,
    .topo-disable-list-leave-to{
        transform: scaleX(0) scaleY(0);
    }
</style>
