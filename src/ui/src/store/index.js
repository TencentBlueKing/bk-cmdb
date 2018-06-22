/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

/*
    Vuex 配置文件
*/

import Vue from 'vue'
import Vuex from 'vuex'
import common from './modules/common'
import main from './modules/main'
import process from './modules/process'
import index from './modules/index'
import hostTransferPop from './modules/hostTransferPop'
import hostSnapshot from './modules/hostSnapshot'
import usercustom from './modules/usercustom'
import navigation from './modules/navigation'
import object from './modules/object'
import association from './modules/association'
Vue.use(Vuex)

export default new Vuex.Store({
    modules: {
        common,
        main,
        process,
        index,
        hostTransferPop,
        hostSnapshot,
        usercustom,
        navigation,
        object,
        association
    }
})
