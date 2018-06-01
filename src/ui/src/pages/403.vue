/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template lang="html">
    <div class="img-403">
        <div class="img-box">
            <img src="../common/images/403.png" alt="">
            <div class="tips-403">
                {{$t('Common["Sorry，您没有业务管理的权限，请联系配置平台管理员"]')}}
                <p :class="{'zh': language === 'zh-cn'}" class="name"></p>
            </div>
        </div>
    </div>
</template>

<script>
    import { mapGetters } from 'vuex'
    export default {
        computed: {
            ...mapGetters([
                'language'
            ])
        }
    }
</script>


<style media="screen" lang="scss">
    .img-403{
        position: relative;
        height: 100%;
        .img-box{
            position: absolute;
            left: 50%;
            top: 50%;
            transform: translate(-50%, -50%);
            width: 650px;
            text-align: center;
            font-size: 16px;
        }
        img{
            display: block;
            margin: 0 auto;
        }
        .tips-403{
            margin-top: 10px;
            font-size: 20px;
            color: #979797;
            .name{
                margin: 10px auto 0;
                width: 500px;
                overflow: hidden;
                text-overflow: ellipsis;
                white-space: break-word;
                white-space: nowrap;
                &.zh {
                    width: 400px;
                }
            }
        }
    }
</style>
