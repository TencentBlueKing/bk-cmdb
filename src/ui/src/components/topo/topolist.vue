<template>
    <ul class="topolist-wrapper clearfix">
        <li class="line" 
            :class="{'default': model['bk_obj_id'] === 'biz'}"
            v-for="(model, index) in topoList"
            :key="index"
            @click="editModel(model)"
            >
            <div class="content">
                <i @click="createModel(model)" class="icon-add prev icon-cc-round-plus" v-if="model['bk_obj_id'] !== 'biz' && model['bk_obj_id'] !== 'module'"></i>
                <div>
                    <i class="topo-icon" :class="model['bk_obj_icon']"></i>
                </div>
                <div class="content-name">
                    {{model['bk_obj_name']}}
                </div>
            </div>
        </li>
    </ul>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                topoList: [],       // 拓扑列表
                topoStructure: []   // 主线拓扑关系
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ])
        },
        methods: {
            createModel (model) {
                let prevObjId = this.findPrevModelId(model)
                this.$emit('createModel', prevObjId)
            },
            editModel (model) {
                this.$emit('editModel', model)
            },
            
            /*
                找到上一级的模型
            */
            findPrevModelId (model) {
                for (let i = 0, topoStructure = this.topoStructure; i < topoStructure.length; i++) {
                    if (model['bk_obj_id'] === topoStructure['bk_obj_id']) {
                        return topoStructure['bk_pre_obj_id']
                    }
                }
            },
            init () {
                this.$Axios.all([
                    this.getTopoStructure(),
                    this.getTopoModel(),
                    this.getApp()
                ]).then(this.$Axios.spread((structureRes, modelRes, appRes) => {
                    this.topoStructure = structureRes
                    let topoList = modelRes[0]['bk_objects']
                    topoList.push(appRes[0])
                    this.sortByStructure(structureRes, topoList)
                }))
            },
            /*
                根据模型拓扑排序
            */
            sortByStructure (structure, topoList) {
                let list = []
                structure.map(({bk_obj_id: objId}) => {
                    for (let i = 0; i < topoList.length; i++) {
                        if (topoList[i]['bk_obj_id'] === objId) {
                            list.push(topoList[i])
                            break
                        }
                    }
                })
                this.topoList = list
            },
            /*
                查询模型拓扑
            */
            getTopoStructure () {
                return this.$axios.get(`topo/model/${this.bkSupplierAccount}`).then(res => {
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
                    bk_classification_id: 'bk_biz_topo'
                }
                return this.$axios.post(`object/classification/${this.bkSupplierAccount}/objects`, params).then(res => {
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    return res.data || []
                })
            },
            /*
                获取业务
            */
            getApp () {
                let params = {
                    bk_obj_id: 'biz'
                }
                return this.$axios.post('objects', params).then(res => {
                    if (!res.result) {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    return res.data || []
                })
            }
        },
        created () {
            this.init()
        }
    }
</script>

<style lang="scss" scoped>
    .topolist-wrapper{
        padding-top: 80px;
        margin: 0 auto;
        li{
            position: relative;
            float: left;
            width: 91px;
            height: 91px;
            border: 1px solid #d6d8df;
            box-shadow: 0 0 10px transparent;
            text-align: center;
            margin-right: 20px;
            margin-top: 30px;
            border-radius: 50%;
            padding: 0 5px;
            cursor: pointer;
            background: #6b7baa;
            color: #fff;
            font-size: 12px;
            font-weight: bold;
            &.line{
                float: none;
                margin: 60px auto 0;
                &:first-child{
                    margin-top: 0;
                }
                &::after{
                    content: "";
                    height: 60px;
                    position: absolute;
                    left: 50%;
                    top: 89px;
                    border: 1px dashed #d6d8df;
                }
                &:last-child{
                    &::after{
                        content: "";
                        height: 100%;
                        position: absolute;
                        left: 56px;
                        top: 50px;
                        border:none;
                    }
                }
            }
            &.default{
                border-style: solid;
                border-width: 1px;
                background: transparent;
                background: #fff !important;
                color: #d6d8df !important;
                cursor: default;
            }
            &:not(.default):hover{
                border: 1px solid #d6d8df;
                box-shadow: 0 2.8px 0 rgba(12, 34, 59, 0.05)
            }
            .content{
                white-space: nowrap;
                text-overflow: ellipsis;
                overflow: hidden;
                padding-top: 20px;
                .topo-icon{
                    font-size: 25px;
                }
                .content-name{
                    margin-top: 10px;
                    line-height: 1;
                    white-space: nowrap;
                    text-overflow: ellipsis;
                    overflow: hidden;
                }
            }
            .icon-add{
                position: absolute;
                display: inline-block;
                border-radius: 50%;
                top: -41px;
                padding: 1px 0;
                font-size: 18px;
                color: #498fe0;
                background: #eee;
                &:hover{
                    color: #50abff;
                }
                &.prev,&.next{
                    left: 36px;
                }
            }
        }
    }
</style>
