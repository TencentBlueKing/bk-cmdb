<template>
    <div class="relevance-topo-wrapper" v-bkloading="{isLoading: isLoading}">
        <div id="topo" class="topo"></div>
        <ul class="model-list" v-if="filterList.length">
            <li class="model" :class="{'unselected': !filter.isShow}" v-for="filter in filterList" @click="changeModelDisplay(filter)" v-if="filter.count">
                <i class="icon" :class="filter['bk_obj_icon']"></i>
                {{filter['bk_obj_name']}} {{filter.count}}
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
                            color: '#737987',
                            vadjust: -5
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
            getInstanceIdKey (objId) {
                if (objId === 'host') {
                    return 'bk_host_id'
                } else if (objId === 'biz') {
                    return 'bk_biz_id'
                }
                return 'bk_inst_id'
            },
            changeModelDisplay (filter) {
                let {
                    activeNode
                } = this
                // activeNode.isShow = !this.activeNode.isShow
                if (activeNode !== null && activeNode.children) {
                    for (let key in activeNode.children) {
                        activeNode.children[key].map(inst => {
                            if (inst['bk_obj_id'] === filter['bk_obj_id']) {
                                inst.isShow = !inst.isShow
                            }
                        })
                        // let inst = activeNode.children[key].find(inst => {
                        //     return inst['bk_obj_id'] === filter['bk_obj_id']
                        // })
                        // if (inst) {
                        //     inst.isShow = !inst.isShow
                        // }
                    }
                }
                this.initTopo()
            },
            getNodes (data, level, isRoot, direction) {
                let nodes = []
                if (isRoot) {
                    nodes.push({
                        id: data.id,
                        label: data['bk_inst_name'],
                        value: 25,
                        image: data.image,
                        level: LEVEL,
                        isLoad: data.isLoad,
                        objId: data['bk_obj_id'],
                        objName: data['bk_obj_name'],
                        objIcon: data['bk_obj_icon'],
                        instId: data['bk_inst_id'],
                        instName: data['bk_inst_name'],
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
            setFilterList () {
                let data = this.activeNode.children
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
            getActiveNode (topoStruct, nodeId, isRoot = false) {
                let activeNode = null
                if (isRoot) {
                    if (topoStruct.id === nodeId) {
                        activeNode = topoStruct
                    }
                }
                if (!activeNode) {
                    for (let key in topoStruct.children) {
                        topoStruct.children[key].map(inst => {
                            if (inst.id === nodeId) {
                                activeNode = inst
                            } else if (inst.hasOwnProperty('children')) {
                                let res = this.getActiveNode(inst, nodeId)
                                if (res) {
                                    activeNode = res
                                }
                            }
                        })
                    }
                }
                return activeNode
            },
            async setTopoStruct (data, isRoot) {
                let count = 0
                let competedNum = 0
                let currentTid = this.activeNode !== null ? this.activeNode.nodeId : `${data['curr']['bk_obj_id']}|${data['curr']['bk_inst_id']}|${LEVEL}|${Math.random().toString(36).substr(2)}`

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
            async getRelationInfo (objId, instId, isRoot = false) {
                this.isLoading = true
                try {
                    const res = await this.$axios.post(`inst/association/topo/search/owner/0/object/${objId}/inst/${instId}`)
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
                this.network.on('click', (params) => {
                    // 点击了具体某个节点
                    if (params.nodes.length) {
                        let id = params.nodes[0]
                        if (this.activeNode) {
                            this.activeNode.image.unselected = this.activeNode.unselectedUrl
                        }
                        this.activeNode = this.getActiveNode(this.topoStruct, id, true)
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
            async deleteRelation (activeNode) {
                let associated = []
                let id = 0
                let parentNode = this.getActiveNode(this.topoStruct, activeNode.parentId, true)
                let toNode = activeNode.fromId === activeNode.id ? parentNode : activeNode
                let fromNode = activeNode.fromId === activeNode.id ? activeNode : parentNode
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

                await this.$store.dispatch('object/getAttribute', toNode['bk_obj_id'])
                let toNodeAttr = this.attribute[toNode['bk_obj_id']].find(({bk_asst_obj_id: bkAsstObjId}) => {
                    return fromNode['bk_obj_id'] === bkAsstObjId
                })
                id = toNodeAttr ? toNodeAttr['bk_property_id'] : ''
                let params = {
                    updateType: 'remove',
                    objId: activeNode['bk_obj_id'],
                    associated: associated,
                    id: id,
                    multiple: !!associated.length,
                    value: fromNode['bk_inst_id'],
                    params: {}
                }
                if (activeNode['bk_obj_id'] === 'host') {
                    params.params['bk_host_id'] = activeNode['bk_inst_id']
                } else {
                    params[this.getInstanceIdKey(activeNode['bk_obj_id'])] = activeNode['bk_inst_id']
                }
                console.log(params)
                // await this.$store.dispatch({
                //     type: 'association/updateAssociation',
                //     ...params
                // })
                // for (let key in parentNode.children) {
                //     let index = parentNode.children[key].findIndex(({id}) => {
                //         return id === activeNode.id
                //     })
                //     if (index > -1) {
                //         parentNode.children[key].splice(index, 1)
                //         break
                //     }
                // }
                this.initTopo()
            },
            showInstDetail (objId, instId) {
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
                let activeNode = this.getActiveNode(this.topoStruct, id, true)
                // 创建popBox
                this.popBox.rand = Math.random().toString(36).substr(2)
                let X = event.clientX
                let Y = event.clientY
                let div = document.createElement('div')
                div.setAttribute('class', 'topo-pop-box')
                div.setAttribute('id', this.popBox.rand)
                div.style.top = `${Y - 40}px`
                div.style.left = `${X}px`
                div.innerHTML = Math.abs(activeNode.level - LEVEL) === 1 ? '<span class="detail" id="instDetail">详情</span> | <span class="color-danger" id="deleteRelation">删除关联</span>' : '<span class="detail" id="instDetail">详情</span>'
                document.body.appendChild(div)

                // 监听事件
                document.getElementById('instDetail').addEventListener('click', (e) => {
                    e.stopPropagation()
                    this.showInstDetail(activeNode['bk_obj_id'], activeNode['bk_inst_id'])
                }, false)
                let deleteElem = document.getElementById('deleteRelation')
                if (deleteElem) {
                    addEventListener('click', (e) => {
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
            padding-left: 30px;
            left: 0;
            top: 0;
            background: #fff;
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
