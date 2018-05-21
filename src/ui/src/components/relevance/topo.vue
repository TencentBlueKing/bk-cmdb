<template>
    <div class="relevance-topo-wrapper" v-bkloading="{isLoading: isLoading}">
        <div id="topo" class="topo"></div>
        <ul class="model-list">
            <li class="model" @click="">
                <i class="icon icon-cc-biz"></i>
                业务
            </li>
            <li class="model">
                <i class="icon icon-cc-cpu"></i>
                业务
            </li>
            <li class="model">
                <i class="icon icon-cc-win"></i>
                业务
            </li>
            <li class="model">
                <i class="icon icon-cc-firewall"></i>
                业务
            </li>
        </ul>
    </div>
</template>

<script>
    import vis from 'vis'
    import { getImgUrl } from '@/utils/util'
    export default {
        props: {
            isShow: {
                type: Boolean,
                default: false
            },
            instId: {
                type: Number
            }
        },
        data () {
            return {
                container: '',
                isLoading: false,
                popBox: {
                    isShow: false,
                    rand: ''
                },
                nodes: [],
                edges: [],
                options: {
                    physics: false,
                    interaction: {
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
                    }
                }
            }
        },
        computed: {
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
                    this.getRelationInfo()
                }
            }
        },
        methods: {
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
            getRelationInfo (instId) {
                this.isLoading = true
                let relationInfo = [
                    {
                        curr: {
                            bk_obj_icon: 'icon-cc-host',
                            bk_obj_id: 'host',
                            bk_obj_name: '主机',
                            bk_inst_id: 21,
                            bk_inst_name: '192.168.1.1'
                        },
                        prev: [
                            {
                                bk_obj_icon: 'icon-cc-subnet',
                                bk_obj_id: 'plat',
                                bk_obj_name: '主机',
                                insts: [
                                    {
                                        bk_inst_id: 1,
                                        bk_inst_name: '父级'
                                    }
                                ],
                                count: 1
                            }
                        ],
                        next: [
                            {
                                bk_obj_icon: 'icon-cc-subnet',
                                bk_obj_id: 'plat',
                                bk_obj_name: '子网区域',
                                insts: [
                                    {
                                        bk_inst_id: 2,
                                        bk_inst_name: '子级'
                                    }
                                ],
                                count: 1
                            }
                        ]
                    }
                ]
                this.formatTopo(relationInfo[0])
            },
            loadImg (src) {
                return new Promise((resolve, reject) => {
                    let img = new Image()
                    img.onload = resolve(img)
                    img.onerror = reject
                    img.src = src
                })
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
                            model.insts.map(inst => {
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
                                        to: key === 'prev' ? inst['bk_inst_id'] : relationInfo['curr']['bk_inst_id'],
                                        from: key === 'prev' ? relationInfo['curr']['bk_inst_id'] : inst['bk_inst_id']
                                    })
                                }

                                // 处理nodes
                                let isNodeExist = nodes.findIndex(({id}) => {
                                    return id === inst['bk_inst_id']
                                }) > -1
                                if (!isNodeExist) {
                                    insertNode.push({
                                        bk_inst_id: inst['bk_inst_id'],
                                        bk_inst_name: inst['bk_inst_name'],
                                        bk_obj_icon: model['bk_obj_icon']
                                    })
                                }
                            })
                        })
                    } else {
                        let isNodeExist = nodes.findIndex(({id}) => {
                            return id === relationInfo[key]['bk_inst_id']
                        }) > -1
                        if (!isNodeExist) {
                            insertNode.push({
                                bk_inst_id: relationInfo[key]['bk_inst_id'],
                                bk_inst_name: relationInfo[key]['bk_inst_name'],
                                bk_obj_icon: relationInfo[key]['bk_obj_icon']
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

                    nodes.push({
                        id: node['bk_inst_id'],
                        label: node['bk_inst_name'],
                        value: this.instId === node['bk_inst_id'] ? 25 : 15,  // 设置大小
                        image: {
                            selected: selectedUrl,
                            unselected: unselectedUrl
                        }
                    })
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
                let data = {
                    nodes: this.nodes,
                    edges: this.edges
                }
                let network = new vis.Network(this.container, this.graphData, this.options)
                // let network = window.network

                // 绑定事件
                let networkCanvas = this.container.getElementsByTagName('canvas')[0]
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
                        this.initPopBox(id, event)

                        // if (this.isBelongtoCurclassify(id) && id !== 'plat') {
                        //     self.$emit('nodeClick', self.getModelById(id))
                        // }
                    }
                })
                this.isLoading = false
            },
            deleteRelation (instId) {

            },
            showInstDetail (instId) {

            },
            initPopBox (instId, event) {
                this.removePop()

                // 创建popBox
                this.popBox.rand = Math.random().toString(36).substr(2)
                let X = event.clientX
                let Y = event.clientY
                let div = document.createElement('div')
                div.setAttribute('class', 'topo-pop-box')
                div.setAttribute('id', this.popBox.rand)
                div.style.top = `${Y - 50}px`
                div.style.left = `${X}px`
                div.innerHTML = '<span class="detail" id="instDetail">详情</span> | <span class="color-danger" id="deleteRelation">删除关联</span>'
                document.body.appendChild(div)

                // 监听事件
                document.getElementById('instDetail').addEventListener('click', (e) => {
                    e.stopPropagation()
                    this.showInstDetail(instId)
                }, false)
                document.getElementById('deleteRelation').addEventListener('click', (e) => {
                    e.stopPropagation()
                    this.deleteRelation(instId)
                }, false)
                // 确保元素已经加载到dom
                this.popBox.isPopShow = true
                document.body.addEventListener('click', this.removePop, false)
                setTimeout(() => {
                    this.popBox.isPopShow = false
                })
            },
            removePop () {
                if (!this.popBox.isPopShow && this.popBox.rand) {
                    let div = document.getElementById(this.popBox.rand)
                    document.body.removeChild(div)
                    this.popBox.rand = ''
                }
            }
        },
        mounted () {
            this.container = document.getElementById('topo')
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
