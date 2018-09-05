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
    <div class="nav-wrapper">
        <div class="list-wrapper">
            <ul class="list-box">
                <router-link tag="li" 
                :to="{path: `/model/${classify['bk_classification_id']}`}"
                v-for="(classify, index) in classifications"
                :key="index">
                    <i :class="classify['bk_classification_icon']"></i>
                    <span class="text">{{classify['bk_classification_name']}}</span>
                </router-link>
            </ul>
            <div class="add-btn-wrapper" @click="createClassify()">
                <a class="add-btn" href="javascript:;">
                    <span>{{$t('Common["新增"]')}}</span>
                </a>
            </div>
            <router-link exact to="/model" :class="['global-models']">
                <button class="btn-global" @click="topoView = 'GLOBAL_MODELS'">
                    <i class="icon-cc-fullscreen"></i>
                    <span class="text">{{$t("ModelManagement['全局视图']")}}</span>
                </button>
            </router-link>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        data () {
            return {
                topoView: 'models',
                localClassify: []
            }
        },
        computed: {
            ...mapGetters(['supplierAccount']),
            ...mapGetters('objectModelClassify', [
                'classifications'
            ]),
            bkClassificationId () {
                return this.$route.params.classifyId
            }
        },
        watch: {
            '$route.params.classifyId' (classifyId) {
                this.init()
            }
        },
        methods: {
            createClassify () {
                this.$emit('createClassify')
            },
            init () {
                if (this.bkClassificationId) {
                    let classify = this.classifications.find(({bk_classification_id: bkClassificationId}) => bkClassificationId === this.$route.params.classifyId)
                    if (!classify) {
                        this.$router.push('/404')
                    }
                    this.topoView = 'models'
                } else {
                    this.topoView = 'GLOBAL'
                }
            }
        },
        created () {
            this.init()
        }
    }
</script>

<style lang="scss" scoped>
    $primaryColor: #737987;
    $primaryHoverColor: #3c96ff; 
    $primaryHoverBgColor: #f1f7ff;
    $primaryActiveBgColor: #e2efff;
    $white: #fff;
    $borderColor: #dde4eb;
    $btnColor: #c3cdd7;
    .nav-wrapper {
        position: relative;
        float:left;
        height: 100%;
        width:188px;
        border-right: 1px solid $borderColor;
    }
    .list-wrapper{
        width: 100%;
        border-left: none;
        border-top: none;
        height: calc(100% - 50px);
        overflow-y: auto;
        @include scrollbar;
        .list-box{
            >li{
                height: 48px;
                line-height: 48px;
                padding: 0 30px 0 44px;
                width: 100%;
                cursor: pointer;
                font-size: 14px;
                color: $primaryColor;
                font-size: 14px;
                position: relative;
                white-space: nowrap;
                text-overflow: ellipsis;
                overflow: hidden;
                i{
                    font-size: 16px;
                }
                .icon-left{
                    margin-left: -12px;
                }
                &:hover{
                    color: $primaryHoverColor;
                    background: $primaryHoverBgColor;
                }
                .text{
                    padding: 0 3px 0 5px;
                    min-width: 64px;
                    vertical-align: top;
                }
                &.active{
                    color: $primaryHoverColor;
                    background: $primaryActiveBgColor;
                }
            }
        }
        .add-btn-wrapper{
            width: 148px;
            height: 32px;
            background: $white;
            cursor: pointer;
            font-size: 0;
            margin: 10px auto;
            .add-btn{
                display: block;
                height: 32px;
                line-height: 30px;
                border-radius: 2px;
                color: $btnColor;
                border: dashed 1px $btnColor;
                text-align: center;
                font-size: 14px;
                &:hover{
                    border-color: $primaryHoverColor;
                    color: $primaryHoverColor;
                }
            }
        }
    }
    .global-models{
        position: absolute;
        bottom: 0;
        left: 0;
        width: 100%;
        height: 50px;
        line-height: 50px;
        text-align: center;
        background-color: #f7fafe;
        &.active{
            background-color: #e2efff;
            .btn-global{
                color: #3578da;
                border-color: currentColor;
            }
        }
        .btn-global{
            width: 120px;
            height: 32px;
            padding: 0;
            line-height: 30px;
            background-color: #ffffff;
            border-radius: 2px;
            border: solid 1px #d6d8df;
            font-size: 14px;
            color: $cmdbTextColor;
            outline: 0;
            padding: 0;
            .icon-cc-fullscreen,
            .text{
                display: inline-block;
                vertical-align: middle;
            }
        }
    }
</style>
