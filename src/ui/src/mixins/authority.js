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
        ...mapGetters(['navigation'])
    },
    data () {
        return {
            unauthorized: {
                search: false,
                create: true,
                update: true,
                delete: true
            }
        }
    },
    watch: {
        navigation () {
            this.setAuthority()
        },
        '$route.fullPath' () {
            this.setAuthority()
        }
    },
    methods: {
        setAuthority () {
            let authority = []
            let navigation = this.navigation
            let path = this.$route.fullPath
            for (let navType in navigation) {
                if (navigation[navType]['authorized']) {
                    let subNav = navigation[navType]['children']
                    let isFound = false
                    for (let i = 0; i < subNav.length; i++) {
                        if (subNav[i]['path'] === path) {
                            authority = subNav[i]['authority']
                            isFound = true
                            break
                        }
                    }
                    if (isFound) break
                }
            }
            this.unauthorized.search = authority.indexOf('search') === -1
            this.unauthorized.create = authority.indexOf('create') === -1
            this.unauthorized.update = authority.indexOf('update') === -1
            this.unauthorized.delete = authority.indexOf('delete') === -1
        }
    },
    beforeMount () {
        this.setAuthority()
    }
}
