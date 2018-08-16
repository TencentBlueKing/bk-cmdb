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
    <div class="topo-box model" id="topo-box" v-bkloading="{isLoading: isTopoLoading}">
        <div class="has-model" v-show="modelList.length !== 0 && modelList.length !== disableModelList.length">
            <div id="topo" class="topo"></div>
            <bk-button type="primary" v-if="addModelAvailable" class="bk-button vis-button vis-create" @click="createModel">
                <span class="vis-button-text">新增</span>
            </bk-button>
        </div>
        <div class="no-model-prompting tc" v-show="disableModelList.length === modelList.length">
            <img src="../../common/images/no_model_prompting.png">
            <p>此分组下无模型</p>
            <bk-button type="primary" class="create-btn" @click="createModel">立即创建</bk-button>
        </div>
        <bk-button class="bk-button vis-button vis-enable" v-if="addModelAvailable && disableModelList.length" @click="isShowDisableList = true">
            <i class="bk-icon icon-minus-circle-shape"></i>
            <span class="vis-button-text">{{disableModelList.length}}</span>
        </bk-button>
        <transition name="topo-disable-list">
            <div class="topo-disable" v-show="isShowDisableList">
                <label class="disable-title">
                    <span>已禁用模型</span>
                    <i class="bk-icon icon-arrows-right" @click="isShowDisableList = false"></i>
                </label>
                <ul class="disable-list" ref="disableList">
                    <li class="disable-item" v-for="(model, index) in disableModelList" :key="index">
                        <a class="disable-item-link" href="javascript:void(0)" @click="">{{model['bk_obj_name']}}</a>
                    </li>
                </ul>
            </div>
        </transition>
    </div>
</template>

<script>
    import vis from 'vis'
    import { mapGetters } from 'vuex'
    export default {
        props: {
            activeClassify: {
                type: Object,
                required: true
            }
        },
        data () {
            return {
                isShowDisableList: false,
                isTopoLoading: false,   // 拓扑图loading
                modelList: [],          // 该分类下所有模型
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
                            color: '#6b7baa',
                            highlight: '#6b7baa',
                            hover: '#6b7baa'
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
                container: ''
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ]),
            addModelAvailable () {
                let notAvailable = ['bk_biz_topo', 'bk_host_manage', 'bk_organization']
                return notAvailable.indexOf(this.activeClassify['bk_classification_id']) === -1
            },
            isCreateBoxShow () {
                return this.disableModelList.length === this.modelList.length && this.modelList.length
            },
            graphData () {
                return {
                    nodes: new vis.DataSet(this.nodes),
                    edges: new vis.DataSet(this.edges)
                }
            },
            /*
                停用的模型列表
            */
            disableModelList () {
                let disableModelList = []
                this.modelList.map(model => {
                    if (model['bk_ispaused']) {
                        disableModelList.push(model)
                    }
                })
                return disableModelList
            }
        },
        watch: {
            'activeClassify' () {
                this.isShowDisableList = false
                this.init()
            }
        },
        methods: {
            /*
                编辑模型
            */
            editModel () {
                this.$emit('editModel')
            },
            /*
                创建模型
            */
            createModel () {
                this.$emit('createModel')
            },
             /*
                查询模型拓扑
            */
            getTopoStructure () {
                let params = {
                    bk_classification_id: this.activeClassify['bk_classification_id']
                }
                return this.$axios.post('objects/topo', params).then(res => {
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    return res.data || []
                })
            },
            /*
                获取模型
            */
            getTopoModel () {
                let params = {
                    bk_classification_id: this.activeClassify['bk_classification_id']
                }
                return this.$axios.post(`object/classification/${this.bkSupplierAccount}/objects`, params).then(res => {
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    return res.data || []
                })
            },
            initEdges (edges) {
                let curClsId = this.activeClassify['bk_classificatin_id']
                this.edges = []
                edges.map(edge => {
                    this.edges.push({
                        dashes: !(edge['from']['bk_classification_id'] === curClsId && edge['to']['bk_classification_id'] === curClsId),
                        label: edge['label'],
                        arrows: edge['arrows'],
                        to: edge['to']['bk_obj_id'],
                        from: edge['from']['bk_obj_id']
                    })
                })
            },
            initNodes () {
                this.nodes = []
                let status = 0
                // 当前状态下的节点数量  已启用/未启用
                let curTypeNum = 0
                let svgForUser = {
                    bgColor: '#fff',
                    iconColor: '#498fe0',
                    color: '#6b7baa'
                }
                let svgForPre = {
                    bgColor: '#6b7baa',
                    iconColor: '#fff',
                    color: '#fff'
                }
                let svgForOther = {
                    bgColor: '#fff',
                    iconColor: '#d6d8df',
                    color: '#d6d8df'
                }
                this.modelList.map(node => {
                    if (!node['bk_ispaused']) {
                        curTypeNum++
                        let svgColor = {}
                        if (node['ispre']) {
                            svgColor = svgForPre
                        } else {
                            svgColor = svgForUser
                        }
                        if (!this.isBelongtoCurclassify(node['bk_obj_id'])) {
                            svgColor = svgForOther
                        }
                        // 没有图标的话就设置一个默认图标
                        if (!node.hasOwnProperty('bk_obj_icon')) {
                            node['bk_obj_icon'] = 'icon-cc-business'
                        }
                        
                        let img = `./static/svg/${this.getIconByClass(node['bk_obj_icon'])}.svg`
                        let image = new Image()
                        image.onload = () => {
                            var base64 = this.getBase64Img(image, this.parseColor(svgColor.iconColor))
                            let svg = `<svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, .1)" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                            <circle cx="50" cy="50" r="49" fill="${svgColor.bgColor}"/>
                                <svg xmlns="http://www.w3.org/2000/svg" stroke="rgba(0, 0, 0, 0)" viewBox="0 0 18 18" x="35" y="-12" fill="${svgColor.iconColor}" width="35" >
                                <image width="15" xlink:href="${base64}"></image>
                                </svg>
                                <foreignObject x="0" y="58" width="100%" height="100%">
                                    <div xmlns="http://www.w3.org/1999/xhtml" style="font-size:14px">
                                        <div style="color:${svgColor.color};text-align: center;width: 50px;overflow:hidden;white-space:nowrap;text-overflow:ellipsis;margin:0 auto">${node['bk_obj_name']}</div>
                                    </div>
                                </foreignObject>
                            </svg>`
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
                                if (position.hasOwnProperty(this.activeClassify['bk_classification_id'])) {
                                    temp.x = position[this.activeClassify['bk_classification_id']].x
                                    temp.y = position[this.activeClassify['bk_classification_id']].y
                                }
                            }
                            this.nodes.push(temp)
                            status++
                        }
                        image.src = img
                    }
                })
                let timer = setInterval(() => {
                    if (status === curTypeNum) {
                        clearInterval(timer)
                        this.initTopo()
                    }
                }, 200)
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
                let dataURL = canvas.toDataURL(`image/${ext}`)
                return dataURL
            },
            /*
                判断模型是否属于当前分类
            */
            isBelongtoCurclassify (objId) {
                for (let i = 0; i < this.modelList.length; i++) {
                    if (this.modelList[i]['bk_obj_id'] === objId) {
                        return this.modelList[i]['bk_classification_id'] === this.activeClassify['bk_classification_id']
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
                for (let i = 0, pos = {}; i < this.modelList.length; i++) {
                    if (this.modelList[i]['bk_obj_id'] === node) {
                        index = i
                        params['bk_classification_id'] = this.modelList[i]['bk_classification_id']
                        params['bk_obj_name'] = this.modelList[i]['bk_obj_name']
                        params['bk_supplier_account'] = this.modelList[i]['bk_supplier_account']

                        if (this.modelList[i].hasOwnProperty('position') && this.modelList[i]['position'] !== '') {
                            pos = JSON.parse(this.modelList[i]['position'])
                        }
                        pos[this.activeClassify['bk_classification_id']] = position[node]
                        params.position = JSON.stringify(pos)
                        id = this.modelList[i]['id']
                        break
                    }
                }
                this.$axios.put(`object/${id}`, params).then(res => {
                    if (res.result) {
                        this.$set(this.modelList[index], 'position', params['position'])
                        this.$store.commit('updateClassifyPosition', this.modelList[index])
                    } else {
                        this.$alertMsg('更新位置信息失败')
                    }
                })
            },
            /*
                获取当前模型
            */
            getModelById (id) {
                return this.modelList.find(({bk_obj_id: bkObjId}) => {
                    return bkObjId === id
                })
            },
            /*
                初始化拓扑图
            */
            initTopo () {
                let network = new vis.Network(this.container, this.graphData, this.options)
                window.network = network
                let networkCanvas = this.container.getElementsByTagName('canvas')[0]

                // 设置按钮title
                let visBtnSetting = [{
                    'title': '放大',
                    'icon': 'icon-plus'
                }, {
                    'title': '缩小',
                    'icon': 'icon-minus'
                }, {
                    'title': '还原',
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
                        if (this.isBelongtoCurclassify(id)) {
                            this.$emit('editModel', this.getModelById(id))
                        }
                    }
                })
                network.on('dragEnd', (params) => {
                    // 拖拽了具体某个节点
                    if (params.nodes.length) {
                        this.savePosition(params.nodes[0])
                    }
                })
                this.isTopoLoading = false
            },
            async init () {
                this.isTopoLoading = true
                let modelList = this.$deepClone(this.activeClassify['bk_objects'])
                for (let key in this.activeClassify['bk_asst_objects']) {
                    let object = this.activeClassify['bk_asst_objects'][key]
                    object.map(asstModel => {
                        let model = modelList.find(({bk_classification_id: bkClassificationId}) => {
                            return bkClassificationId === asstModel['bk_classification_id']
                        })
                        if (!model) {
                            modelList.push(asstModel)
                        }
                    })
                }
                this.modelList = modelList
                this.initEdges(await this.getTopoStructure())
                this.initNodes()
            }
        },
        mounted () {
            this.container = document.getElementById('topo')
            this.$refs.disableList.style.maxHeight = `${Math.ceil(document.body.getBoundingClientRect().height * 2 / 3)}px`
        },
        created () {
            // this.init()
        }
    }
</script>

<style lang="scss" scoped>
    .topo-box{
        position: relative;
        width: 100%;
        height: 100%;
        .topo{
            height: 100%;
        }
    }
    .has-model{
        height: 100%;
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
    .no-model-prompting{
        padding-top: 218px;
        >img{
            display: inline-block;
            width: 200px;
            margin-left: 10px;
        }
        .create-btn{
            width: 208px;
        }
        p{
            font-size: 14px;
            color: #6b7baa;
            margin: 0;
            line-height: 14px;
            margin-top: 23px;
            margin-bottom: 20px;
        }
    }
</style>


<style lang="scss">
    .model#topo-box{
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
            font-size: 0;
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
            .bk-icon,
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
</style>
