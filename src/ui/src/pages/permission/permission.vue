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
    <div class="permission-wrapper">
        <bk-tab class="permission-tab" :active-name="activeTabName" @tab-changed="tabChanged">
            <bk-tabpanel name="role" :title="$t('Permission[\'角色\']')">
                <v-role @on-search-success="setRoles"
                    @skipToUser="skipToUser"
                ></v-role>
            </bk-tabpanel>
            <bk-tabpanel name="authority" :title="$t('Permission[\'权限\']')">
                <v-authority :roles="roles"
                    :activeGroup="activeGroup"
                    :activeTabName.sync="activeTabName"
                ></v-authority>
            </bk-tabpanel>
            <bk-tabpanel name="business" :title="$t('Permission[\'业务权限\']')">
                <v-business></v-business>
            </bk-tabpanel>
        </bk-tab>
    </div>
</template>

<script>
    import vRole from './children/role'
    import vAuthority from './children/authority'
    import vBusiness from './children/business'
    export default {
        components: {
            vRole,
            vAuthority,
            vBusiness
        },
        data () {
            return {
                activeTabName: 'role',
                roles: [],
                activeGroup: {}
            }
        },
        methods: {
            skipToUser (group) {
                this.activeTabName = 'authority'
                this.activeGroup = group
            },
            tabChanged (name) {
                this.activeTabName = name
                if (name === 'role') {
                    this.activeGroup = {}
                }
            },
            setRoles (roles, hasFilter) {
                if (!hasFilter) {
                    this.roles = roles
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .permission-wrapper{
        height: 100%;
        padding: 0 20px;
        overflow: auto;
    }
</style>
<style lang="scss">
    .permission-tab.bk-tab2{
        height: 100%;
        border: none;
        .bk-tab2-head{
            height: 80px;
            .tab2-nav-item{
                height: 79px;
                line-height: 79px;
                min-width: 93px;
                text-align: center;
            }
        }
        .bk-tab2-content{
            height: calc(100% - 80px);
            >section{
                height: 100%;
            }
        }
    }
</style>