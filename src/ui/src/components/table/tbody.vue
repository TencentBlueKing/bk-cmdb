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
    <tbody class="pr">
        <template v-if="tableList.length">
            <slot name="tableBodyRow"
                v-for="(item,index) in tableList" 
                :item="item">
                <tr class="cp" @click="handleRowClick(item)">
                    <slot v-for="header in tableHeader"
                        :name="header.id"
                        :item="item">
                        <td>{{item[header.id]}}</td>
                    </slot>
                </tr>
            </slot>
        </template>
        <template v-else>
            <slot name="tableEmptyRow">
                <tr>
                    <td class="table-empty-col" :colspan="tableHeader.length">
                        <p>暂时没有数据</p>
                    </td>
                </tr>
            </slot>
        </template>
    </tbody>
</template>

<script>
    export default {
        props: {
            tableHeader: {
                type: Array,
                required: true
            },
            tableList: {
                type: Array,
                required: true
            },
            hasCheckbox: {
                type: Boolean,
                default: false
            }
        },
        methods: {
            handleRowClick () {
                this.$emit('handleRowClick', ...arguments)
            }
        }
    }
</script>