/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

/*  vue mixins
    用于检测当前页面的权限
 */

import { mapGetters } from 'vuex'
export default {
    computed: {
        ...mapGetters(['isAdmin']),
        ...mapGetters('navigation', ['authority', 'authorizedNavigation']),
        unauthorized () {
            let fullAuthority = ['search', 'update', 'delete']
            let fullAuthorityClassification = ['bk_host_manage', 'bk_back_config', 'bk_index']
            let authorized = []
            if (this.isAdmin) {
                authorized = fullAuthority
            } else {
                let modelAuthority = this.authority['model_config']
                let model = null
                for (let i = 0; i < this.authorizedNavigation.length; i++) {
                    model = this.authorizedNavigation[i]['children'].find(({path}) => path === this.$route.path)
                    if (model) {
                        break
                    }
                }
                if (model) {
                    if (fullAuthorityClassification.includes(model['classificationId'])) {
                        authorized = fullAuthority
                    } else {
                        authorized = modelAuthority.hasOwnProperty(model['classificationId']) ? modelAuthority[model['classificationId']][model['id']] : []
                    }
                }
            }
            return {
                search: authorized.indexOf('search') === -1,
                update: authorized.indexOf('update') === -1,
                delete: authorized.indexOf('delete') === -1
            }
        }
    }
}
