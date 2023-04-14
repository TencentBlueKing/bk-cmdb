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

import businessSelector from './selector/business.vue'
import businessMixSelector from './selector/business-mix.vue'
import clipboardSelector from './selector/clipboard.vue'
import selector from './selector/selector.vue'
import cloudSelector from './selector/cloud.vue'
import serviceCategorySelector from './selector/service-category.vue'
import details from './details/details.vue'
import form from './form/form.vue'
import formMultiple from './form/form-multiple.vue'
import bool from './form/bool.vue'
import boolInput from './form/bool-input.vue'
import date from './form/date.vue'
import dateRange from './form/date-range.vue'
import time from './form/time.vue'
import int from './form/int.vue'
import float from './form/float.vue'
import longchar from './form/longchar.vue'
import singlechar from './form/singlechar.vue'
import timezone from './form/timezone.vue'
import enumeration from './form/enum.vue'
import enummulti from './form/enummulti.vue'
import enumquote from './form/enumquote.vue'
import objuser from './form/user.vue'
import resize from './other/resize.vue'
import collapseTransition from './transition/collapse.js'
import collapse from './collapse/collapse'
import dotMenu from './dot-menu/dot-menu.vue'
import input from './form/input.vue'
import searchInput from './form/search-input.vue'
import inputSelect from './selector/input-select.vue'
import iconButton from './button/icon-button.vue'
import tips from './other/tips.vue'
import dialog from './dialog/dialog.vue'
import auth from './auth/auth.vue'
import authOption from './auth/auth-option.vue'
import tableEmpty from './table-empty/table-empty.vue'
import dataEmpty from './data-empty/data-empty.vue'
import list from './form/list.vue'
import table from './form/table.vue'
import leaveConfirm from './dialog/leave-confirm.vue'
import textButton from './button/link-button.vue'
import stickyLayout from './other/sticky-layout.vue'
import permission from './permission/embed-permission.vue'
import routerSubview from './other/router-subview.vue'
import organization from './form/organization.vue'
import propertyValue from './other/property-value.vue'
import tagInput from './tag-input/tag-input.vue'
import innertable from './form/inner-table/index.vue'
const install = (Vue) => {
  const components = [
    businessSelector,
    businessMixSelector,
    clipboardSelector,
    selector,
    details,
    form,
    formMultiple,
    bool,
    boolInput,
    date,
    dateRange,
    time,
    int,
    float,
    longchar,
    singlechar,
    timezone,
    enumeration,
    enummulti,
    enumquote,
    objuser,
    resize,
    collapseTransition,
    collapse,
    dotMenu,
    input,
    searchInput,
    inputSelect,
    iconButton,
    tips,
    dialog,
    cloudSelector,
    serviceCategorySelector,
    auth,
    authOption,
    tableEmpty,
    dataEmpty,
    list,
    table,
    leaveConfirm,
    textButton,
    stickyLayout,
    permission,
    routerSubview,
    organization,
    propertyValue,
    tagInput,
    innertable
  ]
  components.forEach((component) => {
    Vue.component(component.name, component)
  })
}

export default {
  install,
  businessSelector,
  businessMixSelector,
  clipboardSelector,
  selector,
  details,
  form,
  formMultiple,
  bool,
  boolInput,
  date,
  dateRange,
  time,
  int,
  float,
  longchar,
  singlechar,
  timezone,
  enumeration,
  enummulti,
  enumquote,
  objuser,
  resize,
  collapseTransition,
  dotMenu,
  input,
  searchInput,
  inputSelect,
  iconButton,
  tips,
  dialog,
  cloudSelector,
  serviceCategorySelector,
  auth,
  authOption,
  tableEmpty,
  dataEmpty,
  list,
  table,
  leaveConfirm,
  textButton,
  stickyLayout,
  permission,
  routerSubview,
  organization,
  propertyValue,
  tagInput
}
