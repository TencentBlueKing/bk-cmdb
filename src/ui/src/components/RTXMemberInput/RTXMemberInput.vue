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
    <div class="member-panel">
        <input ref="memberRef"></input>
    </div>
</template>

<script>

import $ from 'jquery'

export default {
    props: {
        maxLength: {
            type: Number,
            default: null,
            required: false
        },
        IsRequired: {
            type: Boolean,
            default: false,
            required: false
        }
    },
    data () {
        return {
            validater: null,
            developers: []
        }
    },
    mounted () {
        this.init()
    },
    methods: {
        /*
         *  人员选择器初始化
        */
        init () {
            const lenSet = {}
            const _this = this
            if (this.maxLength) {
                lenSet['maxData'] = this.maxLength
            }

            this.validater = $(this.$el)

            $(this.$refs.memberRef).bkMemberSelector({
                ...lenSet,
                type: 'rtx',
                resultScroll: true,
                onSelect (member, e) {
                    e.stopPropagation()
                    e.preventDefault()
                    if (_this.validater.find('.bkMember-required').length === 1) {
                        _this.validater.removeClass('bkMember-data-error').find('.bkMember-required').remove()
                    }
                }
            })

            $(this.$refs.memberRef).data('bkMember').clear()

            /*
             *  人员选择器focus / blur 状态切换时标识类添加（样式定制）
            */
            const memberInputSelector = '.bk-data-editor input[name="bk-data-input"]'
            const selectorObj = this.validater.find(memberInputSelector)
            selectorObj.on('focus blur', (event) => {
                const parentDiv = $(event.target).parents('div.member-panel')
                if ($(event.target).is(':focus')) {
                    parentDiv.addClass('bkMember-data-focus')
                } else {
                    this.developers = $(this.$refs.memberRef).data('bkMember') ? $(this.$refs.memberRef).data('bkMember').getValue() : []
                    parentDiv.removeClass('bkMember-data-focus')
                    this.isValid()
                }
            })
        },

        /*
         *  人员选择器设置值
        */
        setValue (members) {
            $(this.$refs.memberRef).data('bkMember').setValue(members)
            if (members.length !== 0) {
                this.validation()
            } else {
                $(this.$refs.memberRef).data('bkMember').clear()
            }
        },

        /*
         *  人员选择器取值
        */
        getValue () {
            return $(this.$refs.memberRef).data('bkMember').getValue()
        },

        /*
         *  人员选择器验证
        */
        isValid () {
            if (this.IsRequired && this.developers.length === 0) {
                if (this.validater.find('.bkMember-required').length === 0) {
                    this.validater.addClass('bkMember-data-error').append(`<p class="bkMember-required">该字段是必填项</p>`)
                }
            } else {
                this.validater.removeClass('bkMember-data-error').find('.bkMember-required').remove()
            }
        },
        validation () {
            if (this.IsRequired && $(this.$refs.memberRef).data('bkMember').getValue().length === 0) {
                if (this.validater.find('.bkMember-required').length === 0) {
                    this.validater.addClass('bkMember-data-error').append(`<p class="bkMember-required">该字段是必填项</p>`)
                }
                return false
            } else {
                this.validater.removeClass('bkMember-data-error').find('.bkMember-required').remove()
                return true
            }
        },
        clear () {
            // 清空已选人员并取消验证状态
            $(this.$refs.memberRef).data('bkMember').setValue([])
            this.validater.removeClass('bkMember-data-error').find('.bkMember-required').remove()
        }
    },
    watch: {
        developers () {
            this.$emit('getValue', this.developers)
        },
        '$route' () {
            $(this.$refs.memberRef).data('bkMember').setValue([])
        }
    }
}
</script>

<style lang="scss">
    @import './../../assets/bkDataSelector-1.0/css/bkDataSelector.css';
    .bkMember-data-error {
        .bk-data-wrapper {
            border: 1px solid #ff3737!important;
        }
        .bkMember-required {
            margin: 0;
            color: #ff3737;
            font-size: 12px;
        }
    }
</style>
