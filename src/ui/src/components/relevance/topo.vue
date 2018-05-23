<template>
    <div class="relevance-topo-wrapper" v-bkloading="{isLoading: isLoading}">
        <div id="topo" class="topo"></div>
        <ul class="model-list">
            <li class="model" v-for="filter in filterList">
                <i class="icon icon-cc-biz" :class="filter['bk_obj_icon']"></i>
                {{filter['bk_obj_name']}}
            </li>
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
                tid: 0,

                filterList: [],
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
                nodes: [],
                edges: [],
                options: {
                    physics: false,
                    interaction: {
                        // dragNodes: false,
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
            getCurrentNode (topoStruct) {
                let {
                    activeNode
                } = this
                let currentNode = null
                for (let key in topoStruct) {
                    if (key !== 'curr') {
                        for (let i = 0; i < topoStruct[key].length; i++) {
                            let model = topoStruct[key][i]
                            if (model.children !== null) {
                                for (let j = 0; j < model.children.length; j++) {
                                    let inst = model.children[j]
                                    if (activeNode.id === inst.tid) {
                                        currentNode = inst
                                    } else {
                                        if (inst.hasOwnProperty('children')) {
                                            let res = this.getCurrentNode(inst.children)
                                            currentNode = res !== null ? res : null
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
                return currentNode
            },
            setTopoStruct (data) {
                let currentNode = this.getCurrentNode(this.TopoStruct)
                for (let key in data) {
                    if (key !== 'curr') {
                        data[key].map(model => {
                            if (model.children !== null) {
                                model.children.map(inst => {
                                    inst.tid = `${model['bk_obj_id']}|${inst['bk_inst_id']}|${this.tid++}`
                                })
                            }
                        })
                    }
                }
                console.log(currentNode)
                if (currentNode) {
                    this.$set(currentNode, 'children', data)
                } else {
                    this.topoStruct = data
                }
            },
            async getRelationInfo (objId, instId) {
                this.isLoading = true
                try {
                    const res = await this.$axios.post(`inst/association/topo/search/owner/0/object/${objId}/inst/${instId}`)
                    this.setFilterList(res.data[0])
                    this.formatTopo(res.data[0])
                    this.setTopoStruct(res.data[0])
                } catch (e) {
                    this.isLoading = false
                    this.$alertMsg(e.message || e.data['bk_error_msg'] || e.statusText)
                }
            },
            formatTopo (relationInfo) {
                let {
                    nodes,
                    edges
                } = this
                let insertNode = []
                for (let key in relationInfo) {
                    if (key !== 'curr') {
                        relationInfo[key].map(model => {
                            if (model.children !== null) {
                                model.children.map(inst => {
                                    // 处理edges
                                    let isEdgeExist = false
                                    if (key === 'prev') {
                                        isEdgeExist = edges.findIndex(({to, from}) => {
                                            return from === relationInfo['curr']['bk_inst_id'] && to === inst['bk_inst_id']
                                        }) > -1
                                    } else {
                                        isEdgeExist = edges.findIndex(({to, from}) => {
                                            return to === relationInfo['curr']['bk_inst_id'] && from === inst['bk_inst_id']
                                        }) > -1
                                    }
                                    if (!isEdgeExist) {
                                        edges.push({
                                            to: key === 'prev' ? `${model['bk_obj_id']}|${inst['bk_inst_id']}` : `${relationInfo['curr']['bk_obj_id']}|${relationInfo['curr']['bk_inst_id']}`,
                                            from: key === 'prev' ? `${relationInfo['curr']['bk_obj_id']}|${relationInfo['curr']['bk_inst_id']}` : `${model['bk_obj_id']}|${inst['bk_inst_id']}`
                                        })
                                    }
    
                                    // 处理nodes
                                    let isNodeExist = nodes.findIndex(({instId, objId}) => {
                                        return instId === inst['bk_inst_id'] && objId === model['bk_obj_id']
                                    }) > -1
                                    
                                    if (!isNodeExist) {
                                        insertNode.push({
                                            bk_obj_id: model['bk_obj_id'],
                                            bk_obj_name: model['bk_obj_name'],
                                            bk_inst_id: inst['bk_inst_id'],
                                            bk_inst_name: inst['bk_inst_name'],
                                            bk_obj_icon: model['bk_obj_icon'],
                                            id: inst['id'],
                                            relation: key
                                        })
                                    }
                                })
                            }
                        })
                    } else {
                        let isNodeExist = nodes.findIndex(({instId, objId}) => {
                            return instId === relationInfo[key]['bk_inst_id'] && objId === relationInfo[key]['bk_obj_id']
                        }) > -1
                        if (!isNodeExist) {
                            insertNode.push({
                                bk_obj_id: relationInfo[key]['bk_obj_id'],
                                bk_obj_name: relationInfo[key]['bk_obj_name'],
                                bk_inst_id: relationInfo[key]['bk_inst_id'],
                                bk_inst_name: relationInfo[key]['bk_inst_name'],
                                bk_obj_icon: relationInfo[key]['bk_obj_icon'],
                                id: relationInfo[key]['id']
                            })
                        }
                    }
                }

                let count = 0
                insertNode.map(async node => {
                    let src = `./static/svg/${this.getIconByClass(node['bk_obj_icon'])}.svg`
                    let image = await getImgUrl(src)
                    let selectedUrl = this.initImg(image, '#3c96ff')
                    let unselectedUrl = this.initImg(image, '#6c7bb2')

                    let isNodeExist = nodes.findIndex(({instId, objId}) => {
                        return instId === node['bk_inst_id'] && objId === node['bk_obj_id']
                    }) > -1
                    if (!isNodeExist && node['id'] !== '') {
                        let level = 500
                        if (node['bk_obj_id'] === this.objId && node['bk_inst_id'] === this.instId) {
                            level = 500
                        } else {
                            let activeNodeLevel = this.activeNode['level'] ? this.activeNode['level'] : level
                            level = node['relation'] === 'prev' ? activeNodeLevel - 1 : activeNodeLevel + 1
                        }
                        nodes.push({
                            objId: node['bk_obj_id'],
                            instId: node['bk_inst_id'],
                            objName: node['bk_obj_name'],
                            instName: node['bk_inst_name'],
                            id: `${node['bk_obj_id']}|${node['bk_inst_id']}`,
                            label: node['bk_inst_name'],
                            value: this.instId === node['bk_inst_id'] && this.objId === node['bk_obj_id'] ? 25 : 15,  // 设置大小
                            image: unselectedUrl,
                            // image: {
                            //     selected: selectedUrl,
                            //     unselected: unselectedUrl
                            // },
                            selectedUrl: selectedUrl,
                            unselectedUrl: unselectedUrl,
                            level: level
                        })
                    }
                    count++
                })
                let timer = setInterval(() => {
                    if (count === insertNode.length) {
                        clearInterval(timer)
                        this.initTopo()
                    }
                }, 200)
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
                let data = {
                    nodes: this.nodes,
                    edges: this.edges
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
                    if (this.activeNode['image']) {
                        this.activeNode.image = this.activeNode.unselectedUrl
                    }
                })
                this.network.on('click', (params) => {
                    // 点击了某一根线
                    if (params.edges.length) {
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
                        if (this.activeNode['image']) {
                            this.activeNode.image = this.activeNode.unselectedUrl
                        }
                        this.activeNode = this.nodes.find(node => {
                            return id === node.id
                        })
                        this.activeNode.image = this.activeNode.selectedUrl

                        this.getRelationInfo(id.split('|')[0], Number(id.split('|')[1]))
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
                // console.log(params)
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
            }
            .icon {
                position: relative;
                top: 1px;
                vertical-align: bottom;
            }
        }
    }
</style>
