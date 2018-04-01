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
    <div class="discover-content">
        <!-- 左侧tab切换导航 start-->
        <div class="left-tap-contain">
            <div class="list-tap">
                <ul>
                    <li v-for="(item,index) in nav" @click="changeItem(item,index)" :class="{'active': curIndex===index}">
                        <i :class="item.icon"></i>
                        <span class="text" :title="item.name">{{item.name}}</span>
                    </li>
                </ul>
            </div>
        </div>
        <!-- 右侧内容展示 -->
        <div class="right-content">
            <!-- 采集配置 -->
            <div class="acquisition-configuration width-control" v-if="category==='采集配置'">
               <div class="title clearfix">
                   <div class="left fl">
                       <span class="version-munber">版本 v 105.beta</span>
                       <span class="spilt"></span>
                       <span class="version-history">版本历史<img src="../../common/images/left_edition_icon.png" alt=""></span>
                   </div>
                   <div class="right fr">
                       <button class="main-btn mr10">保存</button>
                       <button class="vice-btn mr10">调试</button>
                       <button class="vice-btn" @click="issued">下发</button>
                   </div>
               </div>
               <div class="script-editor pr">
                   <v-ace
                   :value="editorConfig.value"
                   :width="editorConfig.width"
                   :height="editorConfig.height"
                   :lang="editorConfig.lang"
                   >
                   </v-ace>
                    <!-- 下发弹窗 -->
                    <div class="issued-pop" v-show="isIssued">
                        <span class="del-text" @click="close"><i class="bk-icon icon-close f12"></i></span>
                        <div class="issued-detail">
                            <h3>正在下发新版本采集脚本：20/330971  <span class="failure-record">( 失败记录 <span class="num">2</span> )</span></h3>    
                            <div class="progress-content">
                                <div class="progress-bar">
                                </div>
                            </div>
                            <div class="help">
                                <a>查看详情</a>
                                <span>|</span>
                                <a>帮助</a>
                            </div>
                        </div>
                    </div>
               </div>
            </div>
            <!-- 属性绑定 -->
            <div class="attributes-bind width-control" v-else-if="category==='属性绑定'">
                <div class="table-content">
                    <v-table ref="table"
                    :hasCheckbox="false"
                    :tableHeader="tableHeaderModule"
                    :tableList="tableList"
                    :paging="paging"
                    :isLoading="isLoading">
                    </v-table>
                </div>
            </div>
            <!-- 脚本下发记录 -->
            <div class="script-record width-control" v-else-if="category==='脚本下发记录'">
                待定
            </div>
        </div>
    </div>

</template>
<script type="text/javascript">
    import vTable from '@/components/table/table2'
    import vAce from '@/components/ace-editor/ace-editor'
    export default {
        components: {
            vTable,
            vAce
        },
        data () {
            return {
                editorConfig: {               // 脚本编辑器数据配置
                    value: '',
                    width: '100%',
                    height: '533',
                    lang: 'javascript'
                },
                nav: [                         // 导航列表
                    {
                        name: '采集配置',
                        icon: 'icon-cc-collect'
                    },
                    {
                        name: '属性绑定',
                        icon: 'icon-cc-setting'
                    },
                    {
                        name: '脚本下发记录',
                        icon: 'icon-cc-setting'
                    }
                ],
                curIndex: 0,                  // 导航当前索引
                category: '采集配置',          // 导航分类
                paging: {
                    current: 1,
                    totalPage: 1,
                    showItem: 10,
                    sort: ''
                },
                isLoading: false,               // 表格加载
                tableHeaderModule: [
                    {
                        id: 'id',
                        name: '模块名'
                    },
                    {
                        id: 'name',
                        name: '所属集群数'
                    },
                    {
                        id: 'age',
                        name: '主机数量'
                    },
                    {
                        id: 'email',
                        name: '状态'
                    }
                ],
                tableList: [                     // 表格列表
                    {
                        'id': 1463,
                        'name': 'Edward Brown',
                        'age': 41,
                        'email': 'z.qgeyitbrzg@sdhra.by',
                        'url': 'nntp://zlqtkmfd.bf/latxycl'
                    },
                    {
                        'id': 1464,
                        'name': 'Mary Moore',
                        'age': 41,
                        'email': 'z.qgeyitbrzg@sdhra.by',
                        'url': 'nntp://zlqtkmfd.bf/latxycl'
                    },
                    {
                        'id': 1465,
                        'name': 'Jessica Young',
                        'age': 41,
                        'email': 'z.qgeyitbrzg@sdhra.by',
                        'url': 'nntp://zlqtkmfd.bf/latxycl'
                    },
                    {
                        'id': 1466,
                        'name': 'Elizabeth Jones',
                        'age': 41,
                        'email': 'z.qgeyitbrzg@sdhra.by',
                        'url': 'nntp://zlqtkmfd.bf/latxycl'
                    },
                    {
                        'id': 1467,
                        'name': 'Timothy White',
                        'age': 41,
                        'email': 'z.qgeyitbrzg@sdhra.by',
                        'url': 'nntp://zlqtkmfd.bf/latxycl'
                    }
                ],
                isIssued: false                  // 下发弹窗显示
            }
        },
        methods: {
            changeItem (itme, index) {
                this.curIndex = index
                this.category = itme.name
            },
            issued () {
                this.isIssued = true
            },
            close () {
                this.isIssued = false
            }
        }
    }
</script>
<style media="screen" lang="scss" scoped>
$borderColor: #e7e9ef; //边框色
$defaultColor: #ffffff; //默认
$primaryColor: #f9f9f9; //主要
$fnMainColor: #bec6de; //文案主要颜色
$primaryHoverColor: #6b7baa; // 主要颜色
@mixin scrollbar {                 // 页面滚动条
    &::-webkit-scrollbar {
        width: 6px;
        height: 5px;
        &-thumb {
            border-radius: 20px;
            background: #a5a5a5;
            box-shadow: inset 0 0 6px hsla(0,0%,80%,.3);
        }
    }
}
.discover-content{
    width: 100%;
    height: 100%;
    .left-tap-contain{
        float: left;
        width:188px;
        border-right:1px solid $borderColor;
        float:left;
        border-left: none;
        border-top: none;
        height: 100%;
        .list-tap{
            height: 100%;
            overflow-y: auto;
            @include scrollbar;
            ul{
                padding-top:30px;
                >li{
                    width: 100px;
                    height: 40px;
                    line-height: 40px;
                    padding: 0 30px 0 44px;
                    width: 100%;
                    cursor: pointer;
                    font-size: 14px;
                    color: #6b7baa;
                    font-size: 14px;
                    position: relative;
                    white-space:nowrap;
                    text-overflow:ellipsis;
                    -o-text-overflow:ellipsis;
                    overflow: hidden;
                    margin-bottom: 8px;
                    .icon-left{
                        margin-left: -12px;
                    }
                    &:hover{
                        color: #498fe0;
                        background: #f9f9f9;
                        border-right:4px solid #498fe0;
                    }
                    .text{
                        padding:0 3px 0 5px;
                        min-width:64px;
                        vertical-align: top;
                    }
                    &.active{
                        color: #498fe0;
                        background: #f9f9f9;
                        border-right:4px solid #498fe0;
                    }

                }
            }
        }
    }
    .right-content{
        padding: 20px;
        background: #f9f9f9;
        width: calc(100% - 188px);
        height: 100%;
        float: left;
        .width-control{
            width: 100%;
            height: 100%;
        }
        .acquisition-configuration{
            .title{
                width: 100%;
                height: 36px;
                line-height: 36px;
                .left{
                    font-size: 0;
                    .version-munber{
                        font-size: 14px;
                        color: $primaryHoverColor; 
                        font-weight: bold;
                        line-height: 0;
                    }
                    .spilt{
                        width: 2px;
                        height: 12px;
                        color: #e7e9ef;
                        background: #e7e9ef;
                        margin: 0 10px;
                        vertical-align: sub;
                        display: inline-block;
                    }
                    .version-history{
                        font-size: 14px;
                        color: #498fe0;
                        cursor: pointer;
                        >img{
                            vertical-align: middle;
                            margin-left: 3px;
                        }
                    }
                }
                .right{
                    font-size: 0;
                    button{
                        padding: 0 20px;
                        min-width: 120px;
                        font-size: 14px;
                        height: 36px;
                        line-height: 36px;
                    }
                }
            }
        }
        .attributes-bind{
            .table-content{
                table{
                    tbody{
                        background: #fff;
                    }
                }
            }
        }
    }
    .issued-pop{
        position: absolute;
        top: 50%;
        margin-top: -119px;
        left: 50%;
        margin-left: -305px;
        width: 610px;
        height: 238px;
        border-radius: 2px;
        background: #fff;
        padding: 50px 60px;
        .issued-detail{
            >h3{
                margin: 0;
                padding: 0;
                color: $primaryHoverColor;
                font-size: 16px;
                text-align: center;
                .failure-record{
                    font-size: 14px;
                    .num{
                        color: red;
                    }
                }
            }
            .progress-content{
                width: 490px;
                height: 36px;
                border: 1px solid $borderColor;
                margin-top: 30px;
                .progress-bar{
                    width: 50%;
                    float: left;
                    /*width: 0;*/
                    height: 100%;
                    font-size: 12px;
                    line-height: 20px;
                    color: #fff;
                    text-align: center;
                    background: -webkit-linear-gradient(left, #6b7caa , #1fd4fc);
                    -webkit-box-shadow: inset 0 -1px 0 rgba(0,0,0,.15);
                    box-shadow: inset 0 -1px 0 rgba(0,0,0,.15);
                    -webkit-transition: width .6s ease;
                    -o-transition: width .6s ease;
                    transition: width .6s ease;
                }
            }
            .help{
                color: #498fe0;
                line-height: 1;
                text-align: center;
                margin-top: 40px;
                font-size: 14px;
                a{
                    cursor: pointer;
                    &:hover{
                        color: #276dbe;
                    }
                }
                >span{
                    color: $borderColor;
                    margin: 0 5px;
                    display: inline-block;
                }
            }
        }
        .del-text{
            position: absolute;
            top: 0;
            right: 0;
            width: 27px;
            height: 27px;
            line-height: 26px;
            border-radius: 50%;
            text-align: center;
            margin: 4px 4px 0 0;
            border-radius: 50%;
            background-repeat: no-repeat;
            background-size: 11px 11px;
            background-position: 50% 50%;
            cursor: pointer;
            display: inline-block;
            &:hover{
                background-color: #f3f3f3;
            }
        }

    }
}
</style>