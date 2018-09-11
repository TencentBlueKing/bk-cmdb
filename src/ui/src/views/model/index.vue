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
        <v-model-nav @createClassify="createClassify"></v-model-nav>
        <!-- <div class="topo-wrapper">
            <v-topo-list 
                v-if="bkClassificationId === 'bk_biz_topo'"
                @createModel="createModel"
                @editModel="editModel"
                ref="topo"
            ></v-topo-list>
            <v-topo
                v-else-if="bkClassificationId !== void 0"
                @createModel="createModel"
                @editModel="editModel"
                @editClassify="editClassify"
                ref="topo"
            ></v-topo>
            <v-global-models v-else></v-global-models>
        </div> -->
        <router-view
            @createModel="createModel"
            @editModel="editModel"
            @editClassify="editClassify"
        ></router-view>
        <!-- <cmdb-slider
            :hasCloseConfirm="true"
            :isCloseConfirmShow="slider.isCloseConfirmShow"
            :isShow.sync="slider.isShow" :title="slider.title"
            @closeSlider="closeSlider">
            <v-details slot="content"
                ref="details"
                :isEdit.sync="slider.isEdit"
                @createModel="updateTopo(true)"
                @updateModel="updateTopo(false)"
                @cancel="slider.isShow = false"
            ></v-details>
        </cmdb-slider> -->
        <v-pop
            ref="pop"
            v-if="pop.isShow"
            :isEdit="pop.isEdit"
            @closePop="closePop"
        ></v-pop>
    </div>
</template>

<script>
    import vPop from './pop'
    import vModelNav from './model-nav'
    import vDetails from './details'
    import { mapGetters, mapMutations, mapActions } from 'vuex'
    export default {
        components: {
            vModelNav,
            vDetails,
            vPop
        },
        data () {
            return {
                slider: {
                    isShow: false,
                    title: '',
                    isEdit: false,
                    isCloseConfirmShow: false
                },
                pop: {
                    isShow: false,
                    isEdit: false
                }
            }
        },
        computed: {
            bkClassificationId () {
                return this.$route.params.classifyId
            }
        },
        methods: {
            ...mapActions('objectModelClassify', [
                'searchClassificationsObjects'
            ]),
            ...mapMutations('objectModel', [
                'setActiveModel'
            ]),
            createModel (prevModelId) {
                this.slider.title = this.$t('ModelManagement["新增模型"]')
                this.slider.isShow = true
                this.slider.isEdit = false
                this.setActiveModel({
                    bk_classification_id: this.bkClassificationId,
                    bk_asst_obj_id: prevModelId
                })
            },
            editModel (model) {
                this.slider.title = model['bk_obj_name']
                this.setActiveModel(model)
                this.slider.isEdit = true
                this.slider.isShow = true
            },
            editClassify () {
                this.pop.isEdit = true
                this.pop.isShow = true
            },
            createClassify () {
                this.pop.isEdit = false
                this.pop.isShow = true
            },
            closeSlider () {
                this.slider.isCloseConfirmShow = this.$refs.details.isCloseConfirmShow()
                if (!this.slider.isCloseConfirmShow) {
                    this.slider.isShow = false
                }
            },
            closePop () {
                this.pop.isShow = false
            },
            async updateTopo (isShow) {
                await this.searchClassificationsObjects({})
                this.slider.isShow = isShow
                this.$refs.topo.initTopo()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .model-wrapper {
        height: 100%;
        padding: 0;
        .topo-wrapper{
            position: relative;
            float: left;
            width: calc(100% - 188px);
            height: 100%;
            background-color: #f4f5f8;
            background-image: linear-gradient(#eef1f5 1px, transparent 0), linear-gradient(90deg, #eef1f5 1px, transparent 0);
            background-size: 10px 10px;
            @include ellipsis;
            overflow: auto;
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
                            color: $cmdbMainBtnColor;
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
    }
</style>
