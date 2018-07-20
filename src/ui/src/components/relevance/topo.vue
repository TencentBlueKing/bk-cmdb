<template>
    <div class="relevance-topo-wrapper" :class="{'full-screen': isFullScreen}">
        <div class="loading-box" v-bkloading="{isLoading: isLoading}">
            <div id="topo" class="topo" ref="topo"></div>
            <ul class="model-list" v-if="filterList.length">
                <li class="model" :class="{'unselected': !filter.isShow}" v-for="(filter, index) in filterList" @click="changeModelDisplay(filter)" v-if="filter.count" :key="index">
                    <i class="icon" :class="filter['bk_obj_icon']"></i>
                    {{filter['bk_obj_name']}} {{filter.count}}
                </li>
            </ul>
            <span class="resize-btn" v-if="isFullScreen" @click="resizeCanvas(false)">
                <i class="icon-cc-resize-small"></i><span>{{$t('Common["退出"]')}}</span>
            </span>
            <div class="mask" v-if="attr.isShow" @click="attr.isShow = false"></div>
            <v-attribute
                ref="attribute"
                :isFullScreen="isFullScreen"
                :isShow.sync="attr.isShow"
                :instId="attr.instId"
                :objId="attr.objId"
                :instName="attr.instName"
                :objName="attr.objName"
            ></v-attribute>
        </div>
    </div>
</template>

<script>
    import vis from 'vis'
    import vAttribute from './attribute'
    import { getImgUrl } from '@/utils/util'
    import { mapGetters } from 'vuex'
    const LEVEL = 500
    export default {
        props: {
            instId: {
                type: Number
            },
            objId: {
                type: String
            }
        },
        data () {
            return {
                topoStruct: {},
                nodeId: 0,                  // 前端节点ID 递增
                isFullScreen: false,        // 是否全屏显示
                scale: 0.8,
                network: {},
                attr: {
                    isShow: false,
                    instId: '',
                    objId: '',
                    instName: '',
                    objName: ''
                },
                container: '',
                isLoading: false,
                popBox: {
                    isShow: false,
                    rand: '',
                    showPopTimer: 0,
                    timer: 0
                },
                edges: [],
                options: {
                    physics: false,
                    interaction: {
                        dragNodes: false,
                        navigationButtons: true,
                        hover: true
                    },
                    edges: {
                        color: {
                            color: '#c3cdd7',
                            highlight: '#c3cdd7',
                            hover: '#c3cdd7'
                        },
                        smooth: {           // 线的动画
                            enabled: false
                        },
                        arrows: 'middle'
                    },
                    nodes: {
                        font: {
                            color: '#737987',
                            vadjust: -5
                        },
                        shape: 'image',
                        scaling: {
                            min: 15,
                            max: 25
                        },
                        widthConstraint: {
                            maximum: 120
                        }
                    },
                    layout: {
                        hierarchical: {
                            direction: 'LR',
                            nodeSpacing: 90
                        }
                    }
                },
                activeNode: {}
            }
        },
        computed: {
            ...mapGetters('object', [
                'attribute'
            ]),
            graphData () {
                return {
                    nodes: new vis.DataSet(this.nodes),
                    edges: new vis.DataSet(this.edges)
                }
            },
            nodes () {
                return this.getNodes(this.topoStruct, true)
            },
            filterList () {
                let data = this.activeNode.children
                let filterList = []
                for (let key in data) {
                    if (key !== 'curr') {
                        data[key].map(inst => {
                            let current = filterList.find(({bk_inst_id: bkInstId, bk_obj_id: bkObjId}) => {
                                return bkObjId === inst['bk_obj_id']
                            })
                            if (!current) {
                                filterList.push({
                                    bk_obj_id: inst['bk_obj_id'],
                                    bk_obj_name: inst['bk_obj_name'],
                                    bk_obj_icon: inst['bk_obj_icon'],
                                    isShow: inst.isShow,
                                    count: 1
                                })
                            } else {
                                current.count++
                            }
                        })
                    }
                }
                return filterList
            }
        },
        methods: {
            resizeCanvas (isFullScreen) {
                this.isLoading = true
                this.isFullScreen = isFullScreen
                this.scale = this.network.getScale()
                this.$nextTick(() => {
                    this.network.moveTo({scale: this.scale})
                    this.$refs.attribute.resetAttributeBox()
                })
            },
            /**
             * 获取模型Id对应的key
             * @param objId {String} - 模型id
             */
            getInstanceIdKey (objId) {
                if (objId === 'host') {
                    return 'bk_host_id'
                } else if (objId === 'biz') {
                    return 'bk_biz_id'
                }
                return 'bk_inst_id'
            },
            /**
             * 设置筛选
             * @param filter {Object} - 当前点击的筛选项
             */
            changeModelDisplay (filter) {
                let {
                    activeNode
                } = this
                if (activeNode !== null && activeNode.children) {
                    for (let key in activeNode.children) {
                        activeNode.children[key].map(inst => {
                            if (inst['bk_obj_id'] === filter['bk_obj_id']) {
                                inst.isShow = !inst.isShow
                            }
                        })
                    }
                }
                this.initTopo()
            },
            /**
             * 获取node
             * @param data {Object} - 拓扑结构
             * @param isRoot {Boolean} - 是否为根节点
             * @return {Array} - 节点列表
             */
            getNodes (data, isRoot = false) {
                let nodes = []
                if (isRoot) {
                    nodes.push({
                        id: data.id,
                        label: data['bk_inst_name'],
                        value: 25,
                        image: data.image,
                        level: LEVEL,
                        isLoad: data.isLoad,
                        bk_obj_id: data['bk_obj_id'],
                        bk_obj_name: data['bk_obj_name'],
                        bk_obj_icon: data['bk_obj_icon'],
                        bk_inst_id: data['bk_inst_id'],
                        bk_inst_name: data['bk_inst_name'],
                        selectedUrl: data.selectedUrl,
                        unselectedUrl: data.unselectedUrl
                    })
                }
                for (let key in data.children) {
                    data.children[key].map(inst => {
                        if (!inst.isShow) {
                            return
                        }
                        nodes.push(inst)
                        if (inst.hasOwnProperty('children')) {
                            let res = this.getNodes(inst)
                            nodes = nodes.concat(res)
                        }
                    })
                }
                return nodes
            },
            /**
             * 把十六位色值转换为rgb
             * @param color {String} - 十六位色值 例：#123456 / #123
             * @return {Object} - {r: r, g: g, b: b}
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
            /**
             * 获取图片base64
             * @param img {String} - 图片地址
             * @param rgb {Object} - 颜色 rgb: {r: r, g: g, b:b}
             * @return {String} - base64
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
            /**
             * 通过图标类名获取icon路径
             * @param iconClass {String} - icon-xxx-xxx
             * @return {String} - xxx-xxx
             */
            getIconByClass (iconClass) {
                return iconClass.substr(5)
            },
            /**
             * 获取指定节点
             * @param nodeId {Number} - 指定节点id
             * @param topoStruct {Object} - 拓扑结构
             * @param isRoot {Boolean} - 是否为根节点
             */
            getActiveNode (nodeId, topoStruct = this.topoStruct, isRoot = true) {
                let activeNode = null
                if (isRoot) {
                    if (topoStruct.id === nodeId) {
                        activeNode = topoStruct
                    }
                }
                if (!activeNode) {
                    for (let key in topoStruct.children) {
                        for (let index in topoStruct.children[key]) {
                            let inst = topoStruct.children[key][index]
                            if (inst.id === nodeId) {
                                activeNode = inst
                            } else if (inst.hasOwnProperty('children')) {
                                let res = this.getActiveNode(nodeId, inst, false)
                                if (res) {
                                    activeNode = res
                                }
                            }
                            if (activeNode) {
                                break
                            }
                        }
                        if (activeNode) {
                            break
                        }
                    }
                }
                return activeNode
            },
            /**
             * 设置拓扑树形结构
             * @param data {Object} - 从查询关联关系接口直接返回的内容
             * @param isRoot {Boolean} - 页面第一次调用时为true
             */
            async setTopoStruct (data, isRoot) {
                let count = 0
                let competedNum = 0

                let image = await getImgUrl(`./static/svg/${this.getIconByClass(data['curr']['bk_obj_icon'])}.svg`)
                let selectedUrl = this.initImg(image, '#3c96ff')
                let unselectedUrl = this.initImg(image, '#62687f')
                let topoStruct = {
                    prev: [],
                    next: []
                }
                if (isRoot) {
                    this.topoStruct = {
                        isRoot: isRoot,
                        isLoad: true,
                        isShow: true,
                        image: {
                            selected: selectedUrl,
                            unselected: unselectedUrl
                        },
                        bk_inst_id: data['curr']['bk_inst_id'],
                        bk_inst_name: data['curr']['bk_inst_name'],
                        bk_obj_id: data['curr']['bk_obj_id'],
                        bk_obj_name: data['curr']['bk_obj_name'],
                        bk_obj_icon: data['curr']['bk_obj_icon'],
                        id: this.nodeId++,
                        parentId: null,
                        level: isRoot ? LEVEL : this.activeNode.level,
                        selectedUrl: selectedUrl,
                        unselectedUrl: unselectedUrl
                    }
                    this.activeNode = this.topoStruct
                }
                let currentNodeId = this.activeNode !== null ? this.activeNode.id : this.count
                for (let key in data) {
                    if (key !== 'curr') {
                        data[key].map(async model => {
                            if (model.children !== null) {
                                model.children.map(async inst => {
                                    if (inst.id === '') {
                                        return
                                    }
                                    count++
                                    let level = 0
                                    if (isRoot) {
                                        level = key === 'prev' ? LEVEL - 1 : LEVEL + 1
                                    } else {
                                        level = this.activeNode.level - LEVEL < 0 ? this.activeNode.level - 1 : this.activeNode.level + 1
                                    }
                                    let nodeId = this.nodeId++

                                    this.edges.push({
                                        to: key === 'prev' ? nodeId : currentNodeId,
                                        from: key === 'prev' ? currentNodeId : nodeId
                                    })

                                    // 处理nodes
                                    let image = await getImgUrl(`./static/svg/${this.getIconByClass(model['bk_obj_icon'])}.svg`)
                                    let selectedUrl = this.initImg(image, '#3c96ff')
                                    let unselectedUrl = this.initImg(image, '#62687f')
                                    
                                    let parentNode = null
                                    if (this.activeNode.id) {
                                        parentNode = this.getActiveNode(this.activeNode.parentId)
                                    }
                                    // 不重复插入父节点 子网区域
                                    if ((parentNode === null || parentNode['bk_obj_id'] !== model['bk_obj_id'] || parentNode['bk_inst_id'] !== inst['bk_inst_id']) && model['bk_obj_id'] !== 'plat') {
                                        topoStruct[key].push({
                                            isLoad: false,
                                            label: inst['bk_inst_name'],
                                            value: 15,
                                            isShow: true,
                                            image: {
                                                selected: selectedUrl,
                                                unselected: unselectedUrl
                                            },
                                            bk_inst_id: inst['bk_inst_id'],
                                            bk_inst_name: inst['bk_inst_name'],
                                            bk_obj_id: model['bk_obj_id'],
                                            bk_obj_name: model['bk_obj_name'],
                                            bk_obj_icon: model['bk_obj_icon'],
                                            fromId: key === 'prev' ? currentNodeId : nodeId,
                                            id: nodeId,
                                            parentId: this.activeNode.id,
                                            level: level,
                                            selectedUrl: selectedUrl,
                                            unselectedUrl: unselectedUrl
                                        })
                                    }
                                    competedNum++
                                })
                            }
                        })
                    }
                }
                let timer = setInterval(() => {
                    if (count === competedNum) {
                        clearInterval(timer)
                        this.$set(this.activeNode, 'children', topoStruct)
                        this.initTopo()
                    }
                }, 200)
            },
            /**
             * 获取节点关联关系
             * @param objId {String} - 模型ID
             * @param instId {Number} - 实例ID
             * @param isRoot {Boolean} - 页面首次调用时为true
             */
            async getRelationInfo (objId, instId, isRoot = false) {
                this.isLoading = true
                try {
                    const res = await this.$axios.post(`inst/association/topo/search/owner/0/object/${objId}/inst/${instId}`)
                    if (isRoot) {
                        this.$emit('handleAssociationLoaded', res.data[0])
                    }
                    await this.setTopoStruct(res.data[0], isRoot)
                } catch (e) {
                    this.isLoading = false
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            /**
             * 获取图片地址
             * @param image {Object} - Image对象
             * @param color {String} - 颜色色值 例 #123456 / #123
             */
            initImg (image, color) {
                let base64 = this.getBase64Img(image, this.parseColor(color))
                let svg = `<svg xmlns="http://www.w3.org/2000/svg" stroke="" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <rect x="" height="100" width="100" style="fill: #f9f9f9; border: none"/>
                    <image width="100%" xlink:href="${base64}"></image>
                </svg>`
                return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
            },
            /**
             * 绘制拓扑图
             */
            initTopo () {
                this.network = new vis.Network(this.container, this.graphData, this.options)
                this.network.focus(this.activeNode.id)
                this.network.moveTo({scale: 0.8})
                window.network = this.network
                // 绑定事件
                let networkCanvas = this.container.getElementsByTagName('canvas')[0]
                this.network.on('hoverNode', (params) => {
                    let id = params.node
                    networkCanvas.style.cursor = 'pointer'
                    this.initPopBox(id, event)
                })
                this.network.on('blurNode', () => {
                    networkCanvas.style.cursor = 'default'
                })
                this.network.on('dragStart', () => {
                    this.removePop()
                })
                this.network.on('resize', () => {
                    this.isLoading = false
                    this.$nextTick(() => {
                        this.network.moveTo({scale: this.scale})
                    })
                })
                this.network.on('click', (params) => {
                    this.removePop()
                    // 点击了具体某个节点
                    if (params.nodes.length) {
                        let id = params.nodes[0]
                        if (this.activeNode) {
                            this.activeNode.image.unselected = this.activeNode.unselectedUrl
                        }
                        this.activeNode = this.getActiveNode(id)
                        this.activeNode.image.unselected = this.activeNode.selectedUrl

                        // 当前节点没有点击过时 加载其关联内容
                        if (!this.activeNode.isLoad) {
                            this.activeNode.isLoad = true
                            this.getRelationInfo(this.activeNode['bk_obj_id'], this.activeNode['bk_inst_id'])
                        } else {
                            this.initTopo()
                        }
                    }
                })
                this.isLoading = false
            },
            /**
             * 删除关联
             * @param activeNode {Object} - 当前节点
             */
            async deleteRelation (activeNode) {
                let associated = []
                let id = 0
                let parentNode = this.getActiveNode(activeNode.parentId)
                let toNode = activeNode.fromId === activeNode.id ? parentNode : activeNode
                let fromNode = activeNode.fromId === activeNode.id ? activeNode : parentNode
                this.removePop()
                this.isLoading = true
                try {
                    const res = await this.$axios.post(`inst/association/topo/search/owner/0/object/${toNode['bk_obj_id']}/inst/${toNode['bk_inst_id']}`)
                    for (let key in res.data[0]) {
                        if (key !== 'curr') {
                            res.data[0][key].map(model => {
                                if (model['bk_obj_id'] === fromNode['bk_obj_id'] && model.children !== null) {
                                    model.children.map(inst => {
                                        associated.push(inst['bk_inst_id'])
                                    })
                                }
                            })
                        }
                    }
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }

                await this.$store.dispatch('object/getAttribute', {objId: toNode['bk_obj_id']})
                let toNodeAttr = this.attribute[toNode['bk_obj_id']].find(({bk_asst_obj_id: bkAsstObjId}) => {
                    return fromNode['bk_obj_id'] === bkAsstObjId
                })
                id = toNodeAttr ? toNodeAttr['bk_property_id'] : ''
                let params = {
                    updateType: 'remove',
                    objId: toNode['bk_obj_id'],             // 父节点bk_obj_id
                    associated: associated,                 // 已关联的inst_id
                    id: id,                                 // 父节点bk_property_id
                    multiple: associated.length > 1,        // 是否为多关联
                    value: fromNode['bk_inst_id'],          // 子节点bk_inst_id
                    params: {}
                }
                if (toNode['bk_obj_id'] === 'host') {
                    params.params['bk_host_id'] = toNode['bk_inst_id'].toString()
                } else {
                    params[this.getInstanceIdKey(toNode['bk_obj_id'])] = toNode['bk_inst_id']
                }
                await this.$store.dispatch({
                    type: 'association/updateAssociation',
                    ...params
                })
                for (let key in parentNode.children) {
                    let index = parentNode.children[key].findIndex(({id}) => {
                        return id === activeNode.id
                    })
                    if (index > -1) {
                        parentNode.children[key].splice(index, 1)
                        break
                    }
                }
                this.initTopo()
                this.isLoading = false
                this.$emit('handleUpdate')
            },
            /**
             * 显示详情
             * @param objId {String} - 模型ID
             * @param instId {Number} - 实例ID 兼容biz_id host_id 等
             */
            showInstDetail (objId, instId) {
                this.attr.objId = objId
                this.attr.instId = instId
                let model = this.nodes.find(node => {
                    return node['bk_obj_id'] === objId
                })
                this.attr.objName = model ? model['bk_obj_name'] : ''
                let inst = this.nodes.find(node => {
                    return node['bk_inst_id'] === instId
                })
                this.attr.instName = inst ? inst['bk_inst_name'] : ''
                this.removePop()
                this.attr.isShow = true
            },
            /**
             * 初始化pop
             * @param id {Number} - 当前节点id
             * @param event {Object} - 点击事件对象
             * @param time {Number} - pop自动消失时间 单位：ms
             */
            initPopBox (id, event, time = 5000) {
                this.removePop()
                let activeNode = this.getActiveNode(id)
                // 创建popBox
                this.popBox.rand = Math.random().toString(36).substr(2)
                let X = event.clientX
                let Y = event.clientY
                let div = document.createElement('div')
                div.setAttribute('class', 'topo-pop-box')
                div.setAttribute('id', this.popBox.rand)
                div.style.top = `${Y + 12}px`
                div.style.left = `${X + 60}px`
                div.innerHTML = Math.abs(activeNode.level - LEVEL) === 1 ? `<div class="detail" id="instDetail">${this.$t('Common["详情信息"]')}</div><div class="color-danger" id="deleteRelation">${this.$t('Common["删除关联"]')}</div>` : `<div class="detail" id="instDetail">${this.$t('Common["详情信息"]')}</div>`
                document.body.appendChild(div)

                // 监听事件
                document.getElementById('instDetail').addEventListener('click', e => {
                    e.stopPropagation()
                    this.showInstDetail(activeNode['bk_obj_id'], activeNode['bk_inst_id'])
                }, false)
                let deleteElem = document.getElementById('deleteRelation')
                if (deleteElem) {
                    deleteElem.addEventListener('click', e => {
                        e.stopPropagation()
                        this.deleteRelation(activeNode)
                    }, false)
                }
                document.body.addEventListener('click', this.removePop, false)

                clearTimeout(this.popBox.timer)
                this.popBox.timer = setTimeout(() => {
                    this.removePop()
                    clearTimeout(this.popBox.timer)
                }, time)
            },
            /**
             * 删除pop
             */
            removePop () {
                if (this.popBox.rand) {
                    let div = document.getElementById(this.popBox.rand)
                    document.body.removeChild(div)
                    this.popBox.rand = ''
                }
            }
        },
        mounted () {
            this.container = this.$refs.topo
        },
        async created () {
            await this.getRelationInfo(this.objId, this.instId, true)
        },
        components: {
            vAttribute
        }
    }
</script>

<style lang="scss" scoped>
    .relevance-topo-wrapper {
        position: relative;
        height: calc(100% - 64px);
        background: #f9f9f9;
        &.full-screen {
            position: fixed;
            left: 0;
            top: 0;
            bottom: 0;
            right: 0;
            height: 100%;
            .model-list {
                padding: 20px 0 0 20px;
            }
        }
        .loading-box {
            height: 100%;
        }
        .topo {
            height: 100%;
        }
        .model-list {
            position: absolute;
            padding: 10px 0 0 10px;
            left: 0;
            top: 0;
            background: #f9f9f9;
            .model {
                cursor: pointer;
                &.unselected {
                    color: #c3cdd7;
                }
            }
            .icon {
                position: relative;
                top: 1px;
                vertical-align: bottom;
            }
        }
        .resize-btn {
            position: absolute;
            width: auto;
            top: 20px;
            right: 20px;
            height: 24px;
            line-height: 22px;
            padding: 0 10px;
            i {
                margin-right: 5px;
                font-size: 12px;
            }
            span {
                font-size: 12px;
                vertical-align: bottom;
            }
        }
        .mask {
            position: fixed;
            top: 0;
            bottom: 0;
            left: 0;
            right: 0;
            opacity: 0;
        }
    }
</style>

<style lang="scss">
    .relevance-topo-wrapper {
        .attribute-group {
            width: 580px;
            &:first-child {
                padding-top: 8px;
            }
        }
    }
</style>
