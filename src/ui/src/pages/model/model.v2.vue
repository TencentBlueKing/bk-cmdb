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
    <div class="model-wrapper">
        <v-model-nav
            @changeClassify="getActiveClassify"
        ></v-model-nav>
        <div class="topo-wrapper">
            <v-topo-list 
                v-if="activeClassify['bk_classification_id'] === 'bk_biz_topo'"
                @createModel="createModel"
                @editModel="editModel"
            ></v-topo-list>
            <v-topo
                v-else
                :activeClassify="activeClassify"
                @createModel="createModel"
                @editModel="editModel"
            ></v-topo>
            <bk-button type="primary" class="topo-btn edit" @click="editModelClass">
                <i class="icon icon-cc-edit"></i>
            </bk-button>
            <bk-button type="danger" class="topo-btn del" v-if="activeClassify['bk_classification_type']!=='inner'" @click="deleteModelClass">
                <i class="icon icon-cc-del"></i>
            </bk-button>
        </div>
        <v-sideslider
            :isShow.sync="isSlideShow"
            :title="sliderTitle"
        >
            <template slot="content">
                <div class="slide-content">
                    <bk-tab :active-name="activeTabName" @tab-changed="tabChanged">
                        <bk-tabpanel name="info" title="模型配置">
                            <v-field></v-field>
                        </bk-tabpanel>
                        <bk-tabpanel name="layout" title="字段分组">
                            <v-layout
                                :isShow="activeTabName === 'layout'"
                                :activeModel="activeModel"
                            ></v-layout>
                        </bk-tabpanel>
                        <bk-tabpanel name="other" title="其他操作">
                            <v-other
                                :activeClassify="activeClassify"
                                :activeModel="activeModel"
                                @closeSlider="closeSlider"
                            ></v-other>
                        </bk-tabpanel>
                    </bk-tab>
                </div>
            </template>
        </v-sideslider>
    </div>
</template>

<script type="text/javascript">
    import bus from '@/eventbus/bus'
    import vModelNav from './modelNav'
    import vTopoList from '@/components/topo/topolist'
    import vTopo from '@/components/topo/topo2'
    import vSideslider from '@/components/slider/sideslider'
    import vModelInfo from './modelInfo'
    import vField from './children/field2.vue'
    import vLayout from './children/layout2.vue'
    import vOther from './children/other2.vue'
    export default {
        data () {
            return {
                isSlideShow: false,
                activeTabName: '',
                sliderTitle: {
                    text: ''
                },
                activeClassify: {
                    bk_classification_id: '',
                    bk_objects: []
                },
                activeModel: {}
            }
        },
        methods: {
            closeSlider () {
                this.isSlideShow = false
            },
            editModelClass () {
                bus.$emit('editModelClass', true)
            },
            deleteModelClass () {
                bus.$emit('deleteModelClass')
            },
            createModel () {
                this.sliderTitle.text = '新增模型'
                this.isSlideShow = true
            },
            editModel (model) {
                this.sliderTitle.text = model['bk_obj_name']
                this.activeModel = model
                this.isSlideShow = true
            },
            getActiveClassify (activeClassify) {
                this.activeClassify = activeClassify
            },
            tabChanged (name) {
                this.activeTabName = name
            }
        },
        components: {
            vModelNav,
            vTopoList,
            vTopo,
            vSideslider,
            vModelInfo,
            vField,
            vLayout,
            vOther
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    .model-wrapper{
        height: 100%;
        overflow: hidden;
        .topo-wrapper{
            position: relative;
            float: left;
            width: calc(100% - 188px);
            height: 100%;
            background-color: #f4f5f8;
            background-image: linear-gradient(#eef1f5 1px, transparent 0), linear-gradient(90deg, #eef1f5 1px, transparent 0);
            background-size: 10px 10px;
            .topo-btn{
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
                background: #fff;
                .icon{
                    font-size: 14px;
                    color: #737987;
                }
                &.edit{
                    left: 15px;
                    &:hover{
                        .icon{
                            color: #498fe0;
                        }
                    }
                }
                &.del{
                    right: 9px;
                    &:hover{
                        .icon{
                            color: #ef4c4c;
                        }
                    }
                }
            }
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
    }
    .slide-content{
        padding: 8px 20px 20px;
        height: calc(100% - 60px);
    }
</style>

<style lang="scss">
    .model-wrapper{
        .bk-tab2{
            border: none;
        }
    }
</style>
