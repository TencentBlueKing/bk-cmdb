<template>
    <div class="relevance-topo-wrapper" v-bkloading="{isLoading: isLoading}">
        <div id="topo" class="topo"></div>
        <ul class="model-list" v-if="filterList.length">
            <li class="model" :class="{'unselected': !filter.model.isShow}" v-for="filter in filterList" @click="changeModelDisplay(filter)" v-if="filter.count">
                <i class="icon" :class="filter.model['bk_obj_icon']"></i>
                {{filter.model['bk_obj_name']}} {{filter.count}}
            </li>
            <!-- <li class="model" v-for="filter in filterList">
                <i class="icon icon-cc-biz" :class="filter['bk_obj_icon']"></i>
                {{filter['bk_obj_name']}}
            </li> -->
        </ul>
        <v-attribute 
            :isShow.sync="attr.isShow"
            :instId="attr.instId"
            :objId="attr.objId"
            :instName="attr.instName"
            :objName="attr.objName"
        ></v-attribute>
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
            isShow: {
                type: Boolean,
                default: false
            },
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
                nodeId: 0,

                network: {},
                position: {},
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
                            type: 'curvedCW',
                            roundness: 0
                        },
                        arrows: 'middle'
                    },
                    nodes: {
                        font: {
                            color: '#737987'
                        },
                        shape: 'image',
                        scaling: {
                            min: 15,
                            max: 25
                        }
                    },
                    layout: {
                        hierarchical: {
                            direction: 'LR'
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
                return this.getNodes(this.topoStruct, 500, true)
            },
            filterList () {
                let {
                    activeNode
                } = this
                let modelList = []
                if (activeNode !== null && activeNode.children) {
                    for (let key in activeNode.children) {
                        if (key !== 'curr') {
                            console.log(activeNode)
                            activeNode.children[key].map(model => {
                                modelList.push({
                                    model: model,
                                    count: model.children ? model.children.length : 0
                                })
                            })
                        }
                    }
                }
                return modelList
            }
        },
        watch: {
            isShow (isShow) {
                if (isShow) {
                    this.nodes = []
                    this.edges = []
                    this.getRelationInfo(this.objId, this.instId)
                } else {
                    this.activeNode = {}
                }
            }
        },
        methods: {
            changeModelDisplay (filter) {
                let {
                    activeNode
                } = this
                if (activeNode !== null && activeNode.children) {
                    for (let key in activeNode.children) {
                        if (key !== 'curr') {
                            let model = activeNode.children[key].find(model => {
                                return model['bk_obj_id'] === filter.model['bk_obj_id']
                            })
                            if (model) {
                                model.isShow = !model.isShow
                            }
                        }
                    }
                }
                this.initTopo()
            },
            getNodes (data, level, isRoot, direction) {
                let nodes = []
                let localLevel = level
                for (let key in data) {
                    if (key !== 'curr') {
                        if (isRoot) {
                            direction = key === 'prev' ? 'left' : 'right'
                            localLevel = key === 'prev' ? level - 1 : level + 1
                        }
                        data[key].map(model => {
                            if (model.children !== null && model.isShow) {
                                model.children.map(inst => {
                                    nodes.push({
                                        id: inst.nodeId,
                                        label: inst['bk_inst_name'],
                                        value: 15,
                                        image: inst.unselectedUrl,
                                        level: localLevel,
                                        isLoad: inst['isLoad'],
                                        objId: model['bk_obj_id'],
                                        objName: model['bk_obj_name'],
                                        instId: inst['bk_inst_id'],
                                        instName: model['bk_inst_name'],
                                        selectedUrl: inst.selectedUrl,
                                        unselectedUrl: inst.unselectedUrl
                                    })
                                    if (inst.hasOwnProperty('children')) {
                                        let childLevel = direction === 'left' ? localLevel - 1 : localLevel + 1
                                        nodes = nodes.concat(this.getNodes(inst.children, childLevel, false, direction))
                                    }
                                })
                            }
                        })
                    } else {
                        if (isRoot) {
                            let current = data[key]
                            nodes.push({
                                id: current.nodeId,
                                label: current['bk_inst_name'],
                                value: 25,
                                image: current.unselectedUrl,
                                level: localLevel,
                                isLoad: current['isLoad'],
                                objId: current['bk_obj_id'],
                                objName: current['bk_obj_name'],
                                instId: current['bk_inst_id'],
                                instName: current['bk_inst_name'],
                                selectedUrl: current.selectedUrl,
                                unselectedUrl: current.unselectedUrl
                            })
                        }
                    }
                }
                return nodes
            },
            setFilterList (data) {
                let filterList = []
                for (let key in data) {
                    if (key !== 'curr') {
                        data[key].map(model => {
                            if (model.children !== null) {
                                model.children.map(inst => {
                                    let current = filterList.find(({bk_inst_id: bkInstId, bk_obj_id: bkObjId}) => {
                                        return bkObjId === inst['bk_obj_id'] && bkInstId === inst['bk_inst_id']
                                    })
                                    if (!current) {
                                        filterList.push({
                                            bk_obj_id: model['bk_obj_id'],
                                            bk_obj_name: model['bk_obj_name'],
                                            bk_obj_icon: model['bk_obj_icon'],
                                            count: 1
                                        })
                                    } else {
                                        current.count++
                                    }
                                })
                            }
                        })
                    }
                }
                this.filterList = filterList
            },
            getPosition () {
                this.position = this.network.getPositions()
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
                let dataURL = canvas.toDataURL('image/' + ext)
                return dataURL
            },
            /*
                通过图标类名获取icon路径
                class: icon-xxx-xxx
                return xxx-xxx
            */
            getIconByClass (iconClass) {
                return iconClass.substr(5)
            },
            getDeleteModel (topoStruct, nodeId) {
                let deleteModel = null
                for (let key in topoStruct) {
                    if (key !== 'curr') {
                        topoStruct[key].map(model => {
                            if (model.children !== null) {
                                model.children.map(inst => {
                                    if (inst.nodeId === nodeId) {
                                        deleteModel = model
                                    }
                                    if (inst.hasOwnProperty('children') && deleteModel === null) {
                                        let res = this.getDeleteModel(inst.children, nodeId)
                                        if (res) {
                                            deleteModel = res
                                        }
                                    }
                                })
                            }
                        })
                    }
                }
                return deleteModel
            },
            getActiveNode (topoStruct, nodeId) {
                if (!nodeId) {
                    nodeId = this.nodeId
                }
                let activeNode = null
                for (let key in topoStruct) {
                    if (key !== 'curr') {
                        for (let i = 0; i < topoStruct[key].length; i++) {
                            let model = topoStruct[key][i]
                            if (model.children !== null) {
                                for (let j = 0; j < model.children.length; j++) {
                                    let inst = model.children[j]
                                    if (nodeId === inst.nodeId) {
                                        activeNode = inst
                                    } else {
                                        if (inst.hasOwnProperty('children') && activeNode === null) {
                                            let res = this.getActiveNode(inst.children)
                                            if (res) {
                                                activeNode = res
                                            }
                                        }
                                    }
                                }
                            }
                        }
                    } else {
                        if (topoStruct[key].nodeId === nodeId) {
                            activeNode = topoStruct[key]
                        }
                    }
                }
                return activeNode
            },
            async setTopoStruct (data, isRoot) {
                this.activeNode = this.getActiveNode(this.topoStruct)
                let count = 0
                let competedNum = 0
                let currentTid = this.activeNode !== null ? this.activeNode.nodeId : `${data['curr']['bk_obj_id']}|${data['curr']['bk_inst_id']}|${LEVEL}|${Math.random().toString(36).substr(2)}`
                // let topoStruct = {}
                for (let key in data) {
                    // topo
                    if (key !== 'curr') {
                        data[key].map(async model => {
                            // this.$set(model, 'isShow', true)
                            model.isShow = true
                            if (model.children !== null) {
                                model.children.map(async inst => {
                                    count++
                                    let level = 0
                                    if (isRoot) {
                                        level = key === 'prev' ? LEVEL - 1 : LEVEL + 1
                                    } else {
                                        level = this.activeNode.level - LEVEL < 0 ? this.activeNode.level - 1 : this.activeNode.level + 1
                                    }
                                    let nodeId = `${model['bk_obj_id']}|${inst['bk_inst_id']}|${level}|${Math.random().toString(36).substr(2)}`

                                    // 处理edges
                                    this.edges.push({
                                        to: key === 'prev' ? nodeId : currentTid,
                                        from: key === 'prev' ? currentTid : nodeId
                                    })

                                    // 处理nodes
                                    let image = await getImgUrl(`./static/svg/${this.getIconByClass(model['bk_obj_icon'])}.svg`)
                                    let selectedUrl = this.initImg(image, '#3c96ff')
                                    let unselectedUrl = this.initImg(image, '#6c7bb2')
                                    inst.isLoad = false
                                    inst.selectedUrl = selectedUrl
                                    inst.unselectedUrl = unselectedUrl
                                    inst.nodeId = nodeId
                                    inst.level = level
                                    competedNum++
                                })
                            }
                        })
                    } else {
                        // if (!this.activeNode) {
                        count++
                        let image = await getImgUrl(`./static/svg/${this.getIconByClass(data[key]['bk_obj_icon'])}.svg`)
                        let selectedUrl = this.initImg(image, '#3c96ff')
                        let unselectedUrl = this.initImg(image, '#6c7bb2')
                        data[key].selectedUrl = selectedUrl
                        data[key].unselectedUrl = unselectedUrl
                        data[key].isLoad = true
                        data[key].nodeId = currentTid
                        data[key].level = isRoot ? LEVEL : this.activeNode.level
                        competedNum++
                        // }
                    }
                }
                let timer = setInterval(() => {
                    if (count === competedNum) {
                        clearInterval(timer)
                        if (this.activeNode) {
                            this.$set(this.activeNode, 'children', data)
                        } else {
                            this.topoStruct = data
                        }
                        this.initTopo()
                    }
                }, 200)
            },
            async getRelationInfo (objId, instId, isRoot = false) {
                this.isLoading = true
                try {
                    const res = await this.$axios.post(`inst/association/topo/search/owner/0/object/${objId}/inst/${instId}`)
                    // this.setFilterList(res.data[0])
                    await this.setTopoStruct(res.data[0], isRoot)
                } catch (e) {
                    this.isLoading = false
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            initImg (image, color) {
                let base64 = this.getBase64Img(image, this.parseColor(color))
                let svg = `<svg xmlns="http://www.w3.org/2000/svg" stroke="" xmlns:xlink="http://www.w3.org/1999/xlink" width="100" height="100">
                    <rect x="" height="100" width="100" style="fill: #fff; border: none"/>
                    <image width="100%" xlink:href="${base64}"></image>
                </svg>`
                return `data:image/svg+xml;charset=utf-8,${encodeURIComponent(svg)}`
            },
            initTopo () {
                for (let key in this.position) {
                    let node = this.nodes.find(({id}) => {
                        return id === key
                    })
                    if (node) {
                        node.x = this.position[key].x
                        node.y = this.position[key].y
                    }
                }
                this.network = new vis.Network(this.container, this.graphData, this.options)
                // let network = window.network

                // 绑定事件
                let networkCanvas = this.container.getElementsByTagName('canvas')[0]
                this.network.on('hoverNode', (params) => {
                    let id = params.node
                    this.initPopBox(id, event)
                    networkCanvas.style.cursor = 'pointer'
                })
                this.network.on('blurNode', () => {
                    networkCanvas.style.cursor = 'default'
                })
                this.network.on('dragging', () => {
                    if (this.activeNode) {
                        this.activeNode.image = this.activeNode.unselectedUrl
                    }
                })
                this.network.on('click', (params) => {
                    // 点击了某一根线
                    if (!params.nodes.length && params.edges.length) {
                        let edgeId = params.edges[0]
                        let edge = this.edges.find(({id}) => {
                            return id === edgeId
                        })
                        this.deleteRelation(edge)
                    }
                    // 点击了具体某个节点
                    if (params.nodes.length) {
                        this.getPosition()
                        let id = params.nodes[0]
                        if (this.activeNode) {
                            this.activeNode.image = this.activeNode.unselectedUrl
                        }
                        this.nodeId = id
                        this.activeNode = this.getActiveNode(this.topoStruct)
                        this.activeNode.image = this.activeNode.selectedUrl

                        // 当前节点没有点击过时 加载其关联内容
                        if (!this.activeNode.isLoad) {
                            this.activeNode.isLoad = true
                            this.getRelationInfo(id.split('|')[0], Number(id.split('|')[1]))
                        }
                    }
                })
                this.isLoading = false
            },
            async deleteRelation (edge) {
                let associated = []
                let id = 0
                try {
                    const res = await this.$axios.post(`inst/association/topo/search/owner/0/object/${edge.to.split('|')[0]}/inst/${edge.to.split('|')[1]}`)
                    if (res.data[0]['next']) {
                        res.data[0]['next'].map(model => {
                            if (model['bk_obj_id'] === edge.to.split('|')[0] && model.children !== null) {
                                model.map(inst => {
                                    associated.push(inst['bk_inst_id'])
                                })
                            }
                        })
                    }
                } catch (e) {
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }

                await this.$store.dispatch('object/getAttribute', edge.to.split('|')[0])
                let parentAttrList = this.attribute[edge.to.split('|')[0]]
                let parentAttr = parentAttrList.find(({bk_asst_obj_id: bkAsstObjId}) => {
                    return edge.to.split('|')[0] === bkAsstObjId
                })
                if (parentAttr) {
                    id = parentAttr['bk_property_id']
                }
                let params = {
                    updateType: 'update',
                    objId: edge.to.split('|')[0],
                    associated: associated,
                    id: id,
                    value: edge.from.split('|')[1],
                    params: {

                    }
                }

                // this.activeNode.children = null
                // this.initTopo()
                // this.topoStruct['next'][1].children = null
                let parentNode = this.getActiveNode(this.topoStruct, edge.to)
                let childNode = this.getActiveNode(this.topoStruct, edge.from)
                console.log(parentNode, childNode, edge.to, edge.from, 111)

                let deleteObj = Math.abs(parentNode - LEVEL) > Math.abs(childNode - LEVEL) ? childNode : parentNode
                let targetObj = Math.abs(parentNode - LEVEL) > Math.abs(childNode - LEVEL) ? parentNode : childNode
                console.log(deleteObj.nodeId)
                // let deleteModel = this.getDeleteModel(this.topoStruct, deleteObj.curr['nodeId'])
                // deleteObj = null
                // targetObj = null
                console.log(deleteObj, 'deleteObj')
                // console.log(deleteModel)
                // for (let key in deleteModel) {
                //     if (key !== 'curr') {
                //         deleteModel[key].map(model => {
                //             if (model.children !== null) {
                //                 model.children.map(inst => {
                //                 })
                //             }
                //         })
                //     }
                // }
                // deleteModel.children = null
                this.initTopo()
                // this.$store.dispatch('association/updateAssociation', )
            },
            showInstDetail (id) {
                let objId = id.split('|')[0]
                let instId = Number(id.split('|')[1])
                this.attr.objId = objId
                this.attr.instId = instId

                let model = this.nodes.find(node => {
                    return node.objId === objId
                })
                this.attr.objName = model ? model.objName : ''
                let inst = this.nodes.find(node => {
                    return node.instId === instId
                })
                this.attr.instName = inst ? inst.instName : ''

                this.removePop()
                this.attr.isShow = true
            },
            initPopBox (id, event, time = 5000) {
                this.removePop()

                // 创建popBox
                this.popBox.rand = Math.random().toString(36).substr(2)
                let X = event.clientX
                let Y = event.clientY
                let div = document.createElement('div')
                div.setAttribute('class', 'topo-pop-box')
                div.setAttribute('id', this.popBox.rand)
                div.style.top = `${Y - 40}px`
                div.style.left = `${X}px`
                div.innerHTML = '<span class="detail" id="instDetail">详情</span> | <span class="color-danger" id="deleteRelation">删除关联</span>'
                document.body.appendChild(div)

                // 监听事件
                document.getElementById('instDetail').addEventListener('click', (e) => {
                    e.stopPropagation()
                    this.showInstDetail(id)
                }, false)
                document.getElementById('deleteRelation').addEventListener('click', (e) => {
                    e.stopPropagation()
                    this.deleteRelation(id)
                }, false)
                document.body.addEventListener('click', this.removePop, false)

                clearTimeout(this.popBox.timer)
                this.popBox.timer = setTimeout(() => {
                    this.removePop()
                    clearTimeout(this.popBox.timer)
                }, time)
            },
            removePop () {
                if (this.popBox.rand) {
                    let div = document.getElementById(this.popBox.rand)
                    document.body.removeChild(div)
                    this.popBox.rand = ''
                }
            }
        },
        mounted () {
            this.container = document.getElementById('topo')
        },
        created () {
            this.getRelationInfo(this.objId, this.instId, true)
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
        .topo {
            height: 100%;
        }
        .model-list {
            position: absolute;
            right: 30px;
            top: 0;
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
    }
</style>
