/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

/**
 * 内置到 Vue.proottype 的全局变量，命名规则为 $ 开头
 */
import Vue from 'vue'

/**
 * @global Site 是后台编译时内置在 builder/config/index.js 中变量，随着前端编译后会渲染到 index.html 并保存在全局变量 window.Site 中。
 */
Vue.prototype.$Site = window.Site
