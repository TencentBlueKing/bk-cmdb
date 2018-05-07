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
    <div class="business-wrapper">
        <div class="selector-wrapper clearfix">
            <h2 class="selector-title fl">{{$t('Permission["业务角色"]')}}</h2>
            <div class="selector-container clearfix">
                <bk-select class="role-selector fl" :selected.sync="selectedBusinessRole">
                    <bk-select-option v-for="(role, index) in businessRoles"
                        :key="index"
                        :value="role['bk_property_id']"
                        :label="role['bk_property_name']"
                    ></bk-select-option>
                </bk-select>
                <!-- <router-link class="fl" :to="{path:'model', query: {'bk_classification_id': 'bk_organization'}}" v-show="false"><i class="icon-cc-derivation"></i>{{$t('Permission["角色配置"]')}}</router-link> -->
            </div>
        </div>
        <div class="authority-wrapper clearfix">
            <h2 class="authority-title fl">{{$t('Permission["功能选择"]')}}</h2>
            <div class="checkbox-container clearfix">
                <span v-for="authority in authorities.list" class="checkbox-span fl">
                    <label class="bk-form-checkbox bk-checkbox-small authority-checkbox"
                        :class="{'disabled': isMaintainers}"
                        :for="'business-authority-' + authority.id"
                        :title="$t(authority.name)"
                        @click="updateAuthorities">
                        <input type="checkbox" 
                            :value="authority.id" 
                            :id="'business-authority-' + authority.id" 
                            :disabled="isMaintainers"
                            v-model="authorities.selected">
                        {{$t(authority.name)}}
                    </label>
                </span>
            </div>
        </div>
    </div>
</template>
<script>
    import { mapGetters } from 'vuex'
    import Throttle from 'lodash.throttle'
    export default {
        computed: {
            ...mapGetters([
                'bkSupplierAccount'
            ]),
            isMaintainers () {
                return this.selectedBusinessRole === 'bk_biz_maintainer'
            }
        },
        data () {
            return {
                selectedBusinessRole: '',
                businessRoles: [],
                authorities: {
                    list: [{
                        id: 'hostupdate',
                        name: 'Permission["主机编辑"]'
                    }, {
                        id: 'hosttrans',
                        name: 'Permission["主机转移"]'
                    }, {
                        id: 'topoupdate',
                        name: 'Permission["拓扑编辑"]'
                    }, {
                        id: 'customapi',
                        name: 'Permission["自定义查询"]'
                    }, {
                        id: 'proconfig',
                        name: 'Permission["进程管理"]'
                    }],
                    selected: []
                }
            }
        },
        watch: {
            selectedBusinessRole () {
                this.getBusinessRoleAuthorities()
            },
            isMaintainers (isMaintainers) {
                if (isMaintainers) {
                    this.authorities.selected = this.authorities.list.map(authority => {
                        return authority.id
                    })
                }
            }
        },
        methods: {
            getBusinessRoles () {
                this.$axios.post('object/attr/search', {bk_obj_id: 'biz', bk_supplier_account: this.bkSupplierAccount}).then((res) => {
                    if (res.result) {
                        let roles = []
                        res.data.map((role) => {
                            if (role['bk_property_type'] === 'objuser') {
                                roles.push(Object.assign(role, {selectedAuthorities: []}))
                            }
                        })
                        this.businessRoles = roles
                        this.selectedBusinessRole = roles[0]['bk_property_id']
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            getBusinessRoleAuthorities () {
                this.$axios.get(`topo/privilege/${this.bkSupplierAccount}/biz/${this.selectedBusinessRole}`).then((res) => {
                    if (res.result) {
                        if (!this.isMaintainers) {
                            this.authorities.selected = res.data.length ? res.data : []
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            updateAuthorities: Throttle(function () {
                this.$nextTick(() => {
                    this.$axios.post(`topo/privilege/${this.bkSupplierAccount}/biz/${this.selectedBusinessRole}`, this.authorities.selected)
                    .then((res) => {
                        if (!res.result) {
                            this.$alertMsgThrottle(res['bk_error_msg'])
                        }
                    })
                })
            }, 300, {leading: false, trailing: true})
        },
        mounted () {
            this.getBusinessRoles()
        }
    }
</script>

<style lang="scss" scoped>
    .business-wrapper{
        padding: 50px 0 0 0;
        color: #737987;
        font-size: 14px;
    }
    .selector-wrapper{
        .selector-title{
            text-align: right;
            width: 137px;
            height: 36px;
            line-height: 36px;
            font-size: 14px;
            margin: 0;
            padding: 0 30px 0 0;
        }
        .selector-container{
            padding-left: 96px;
            overflow: visible;
            .role-selector{
                width: 286px;
            }
            a{
                height: 14px;
                line-height: 1;
                display: inline-block;
                vertical-align: unset;
                margin: 11px 0 0 26px;
                color: #498fe0;
                text-decoration: none;
                .icon-cc-derivation{
                    vertical-align: -1px;
                    margin: 0 4px 0 0;
                }
            }
        }
    }
    .authority-wrapper{
        padding: 40px 0 0 0;
        .authority-title{
            width: 137px;
            height: 14px;
            text-align: right;
            line-height: 14px;
            font-size: 14px;
            padding: 0 30px 0 0;
            margin: 0;
        }
        .checkbox-container{
            width: 741px;
            padding-left: 137px;
            overflow: visible;
            line-height: 14px;
        }
    }
    .checkbox-span{
        width: 180px;
        height: 14px;
        margin: 0 16px 28px 0;
        .authority-checkbox{
            max-width: 180px;
            line-height: 14px;
            padding: 0;
            margin: 0;
            white-space: nowrap;
            overflow: hidden;
            text-overflow: ellipsis;
            cursor: pointer;
            &.disabled{
                cursor: not-allowed;
            }
            input{
                width: 14px;
                height: 14px;
                vertical-align: top;
                &:checked,
                &:focus{
                    border-color: transparent !important;
                }
            }
        }
    }
</style>