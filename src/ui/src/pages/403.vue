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
        <img src="../common/images/403.png" alt="">
        <p class="tips-403">Sorry，您的权限不足</p>
    </div>
</template>

<style media="screen" lang="scss">
    .img-403{
        position: relative;
        min-height: calc(100% - 50px);
        img{
            position: absolute;
            left: 50%;
            top: 50%;
            transform: translate(-50%, -50%);
        }
        .tips-403{
            font-size: 20px;
            color: #979797;
            position: absolute;
            left: 50%;
            top: 50%;
            margin-top: 108px;
            transform: translate(-50%, -50%);
        }
    }
</style>
