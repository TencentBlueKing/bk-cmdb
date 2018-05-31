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
    <bk-selector class="userapi-compare-selector fl"
        :list="operatorList"
        :selected.sync="currentSelected"
        @item-selected="operatorSelectedHandler">
    </bk-selector>
</template>

<script>
    export default {
        props: {
            type: {
                type: String,
                required: true
            },
            selected: {
                default: ''
            },
            property: {
                type: Object
            }
        },
        data () {
            return {
                typeMap: {
                    'default': '$eq,$ne',
                    'singlechar': '$regex,$eq,$ne',
                    'longchar': '$regex,$eq,$ne',
                    'objuser': '$regex,$eq,$ne',
                    'singleasst': '$regex,$eq,$ne',
                    'multiasst': '$regex,$eq,$ne',
                    'name': '$in,$eq,$ne'
                },
                operatorMap: {
                    '$nin': this.$t("Common['不包含']"),
                    '$in': 'IN',
                    '$regex': this.$t("Common['包含']"),
                    '$eq': this.$t("Common['等于']"),
                    '$ne': this.$t("Common['不等于']")
                }
            }
        },
        computed: {
            operatorList () {
                let list = []
                let operators = null
                if (this.property['bkPropertyId'] === 'bk_module_name' || this.property['bkPropertyId'] === 'bk_set_name') {
                    operators = this.typeMap['name'].split(',')
                } else if (this.typeMap.hasOwnProperty(this.type)) {
                    operators = this.typeMap[this.type].split(',')
                } else {
                    operators = this.typeMap['default'].split(',')
                }
                operators.forEach((operator, index) => {
                    if (this.operatorMap.hasOwnProperty(operator)) {
                        list.push({
                            id: operator,
                            name: this.operatorMap[operator]
                        })
                    }
                })
                return list
            },
            currentSelected: {
                get () {
                    if (!this.selected) {   // 初始化默认选中，通知父组件
                        this.$emit('update:selected', this.operatorList[0]['id'])
                    }
                    return this.selected || this.operatorList[0]['id']
                },
                set (val) {
                    this.$emit('update:selected', val)
                }
            }
        },
        methods: {
            operatorSelectedHandler (id, data, index) {
                this.$emit('item-selected', id, data, index)
            }
        }
    }
</script>