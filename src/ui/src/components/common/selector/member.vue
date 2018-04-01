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
    <div class="bk-select" v-click-outside="close">
        <div class="bk-select-wrapper" @click="toggleSlide">
            <div class="member-selected" :class="{'disabled': disabled}">
                <template v-if="localSelected.length">
                    <span class="member-selected-item"
                        v-for="(englishName, index) in localSelected" 
                        :class="{'disabled': disabled}" 
                        :key="index"
                        :title="englishName"
                        @click.stop="deleteMember(englishName, index)">
                        {{englishName}}
                    </span>
                </template>
                <template v-else>
                    <span class="member-empty">{{placeholder}}</span>
                </template>
            </div>
        </div>
        <transition name="toggle-slide" v-if="!disabled">
            <div class="bk-select-list" v-show="open">
                <div class="bk-select-list-filter">
                    <input ref="filter" type="text" autofocus="autofocus" class="bk-select-filter-input" v-model="filter">
                    <i class="bk-icon icon-search"></i>
                </div>
                <ul>
                    <li ref="memberItem" :class="['bk-select-list-item',{'selected': checkSelected(member)}]" 
                        v-for="(member, index) in members">
                        <label class="bk-form-checkbox bk-checkbox-small bk-select-list-label" 
                            :title="getLabel(member, index)" 
                            @click="setSelected(member, index)">
                            <input type="checkbox" 
                                v-if="multiple"
                                :checked="checkSelected(member, index)">
                            {{getLabel(member, index)}}
                        </label>
                    </li>
                </ul>
            </div>
        </transition>
    </div>
</template>

<script>
    import { mapGetters, mapActions } from 'vuex'
    export default {
        props: {
            placeholder: {
                type: String,
                required: false,
                default: '请选择'
            },
            multiple: {
                type: Boolean,
                required: false,
                default: false
            },
            selected: {
                type: String,
                required: true,
                default: ''
            },
            disabled: {
                type: Boolean,
                default: false
            }
        },
        computed: {
            ...mapGetters({
                'members': 'memberList',
                'memberLoading': 'memberLoading'
            }),
            localSelected () {
                if (this.multiple && this.selected) {
                    return this.selected.split(',')
                } else if (this.selected) {
                    return [this.selected]
                } else {
                    return []
                }
            }
        },
        data () {
            return {
                open: false,
                filter: ''
            }
        },
        watch: {
            open (isOpen) {
                if (isOpen) {
                    this.$nextTick(() => {
                        this.$refs.filter.focus()
                    })
                } else {
                    this.filter = ''
                }
                this.$emit('member-toggle', isOpen)
            },
            filter (val) {
                val = val.toLowerCase()
                this.members.map((member, index) => {
                    if (member.english_name.toLowerCase().indexOf(val) === -1) {
                        this.$refs.memberItem[index].style.display = 'none'
                    } else {
                        this.$refs.memberItem[index].style.display = 'block'
                    }
                })
            }
        },
        methods: {
            ...mapActions(['getMemberList']),
            getLabel (member, index) {
                if (member.chinese_name && member.chinese_name !== member.english_name) {
                    return `${member.english_name}(${member.chinese_name})`
                } else {
                    return member.english_name
                }
            },
            checkSelected (member, index) {
                if (this.multiple) {
                    return this.localSelected.indexOf(member.english_name) !== -1
                } else {
                    return this.selected === member.english_name
                }
            },
            setSelected (member, index) {
                let selected
                if (this.multiple) {
                    selected = this.selected.length ? this.selected.split(',') : []
                    if (selected.indexOf(member.english_name) !== -1) {
                        selected.splice(selected.indexOf(member.english_name), 1)
                    } else {
                        selected.push(member.english_name)
                    }
                    selected = [...new Set(selected)].join(',')
                } else {
                    selected = member.english_name
                    this.close()
                }
                this.$emit('update:selected', selected)
                this.$emit('on-select', member)
            },
            deleteMember (englishName, index) {
                if (!this.disabled) {
                    let selected = this.localSelected.slice()
                    selected.splice(index, 1)
                    this.$emit('update:selected', selected.join(','))
                    this.$emit('on-delete', englishName)
                }
            },
            toggleSlide () {
                if (!this.disabled) {
                    this.open = !this.open
                }
            },
            close () {
                this.open = false
            }
        },
        mounted () {
            if (!this.members.length && !this.memberLoading) {
                this.getMemberList()
            }
        }
    }
</script>

<style lang="scss" scoped>
    .member-empty{
        display: block;
        height: 18px;
        margin-top: 6px;
        padding: 0 2px;
        color: #c3cdd7;
    }
    .member-selected{
        border: solid 1px #c3cdd7;
        border-radius: 2px;
        padding: 2px 6px 8px;
        line-height: 18px;
        max-height: 108px;
        overflow: auto;
        &.disabled{
            background-color: #fafafa;
            cursor: not-allowed;
            .member-selected-item.disabled{
                cursor: not-allowed;
            }
        }
        &::-webkit-scrollbar {
            width: 6px;
            height: 5px;
        }
        &::-webkit-scrollbar-thumb {
            border-radius: 20px;
            background: #a5a5a5;
            -webkit-box-shadow: inset 0 0 6px hsla(0,0%,80%,.3);
        }
    }
    .member-selected-item{
        display: inline-block;
        height: 18px;
        line-height: 16px;
        padding: 0 4px;
        margin: 6px 2px 0;
        background-color: #fafafa;
        border: 1px solid #d9d9d9;
        border-radius: 2px;
        cursor: pointer;
        &.disabled{
            cursor: default;
        }
    }
    .bk-select-list{
        top: 100%;
        margin-top: 4px;
    }
    .bk-select-list-item{
        padding: 0;
    }
    .bk-select-list-label{
        display: block;
        padding: 0 12px;
        height: 42px;
        line-height: 42px;
        cursor: pointer;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }
</style>