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
    <transition name="fade">
        <div class="alertbox" v-show="show">
            <div class="alert-left">
                <i v-if="options.msgIsSuccess" class="bk-icon icon-check-circle-shape"></i>
                <i v-else class="bk-icon icon-close-circle-shape"></i>
            </div>
            <div class="alert-right">
                {{options.msgText}}
            </div>
        </div>
    </transition>
</template>

<script type="text/javascript">
    export default {
        data () {
            return {
                timers: [],
                show: false
            }
        },
        props: {
            options: {
                type: Object,
                default: {}
            }
        },
        watch: {
            options () {
                this.timers.forEach((timer) => {
                    window.clearTimeout(timer)
                })
                this.timers = []
                this.countdown()
            }
        },
        methods: {
            close () {
                this.show = false
            },
            countdown () {
                this.show = true
                const t = setTimeout(() => {
                    this.close()
                }, 1000)
                this.timers.push(t)
            }
        }
    }
</script>

<style media="screen" lang="scss" scoped>
    .alertbox{
        box-shadow: 0 0 10px #cccccc;
        border-radius: 4px;
        background-color: #fff;
        opacity: 0.98;
        position: fixed;
        top: 50%;
        left: 50%;
        transform: translate(-50%, -50%);
        -o-transform: translate(-50%, -50%);
        -ms-transform: translate(-50%, -50%);
        -moz-transform: translate(-50%, -50%);
        -webkit-transform: translate(-50%, -50%);
        padding: 17px 44px 17px 24px;
        z-index: 99999999;
        .alert-left{
            margin-right: 20px;
            width: 40px;
            height: 40px;
            float: left;
            .bk-icon{
                font-size: 40px;
                display: inline-block;
                background-size: cover;
                border-radius: 50%;
                &.icon-check-circle-shape{
                    background: #ffffff;
                    border-color: #30d878;
                    color: #30d878;
                }
                &.icon-close-circle-shape{
                    color: #ff5656;
                    background: #ffffff;
                    border-color: #ff5656;
                }
            }
        }
        .alert-right{
            float: left;
            color: #4f515e;
            line-height: 40px;
            vertical-align: middle;
        }
    }
    .fade-enter,.fade-leave-active{
        transition: opacity .5s;
        opacity: 0;
    }
    .fade-enter-active,.fade-leave{
        transition: opacity .5s;
        opacity: .98;
    }
</style>
