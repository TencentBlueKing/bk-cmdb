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

import moment from 'moment'
import has from 'has'

const defaultFormatter = (value) => {
  if (!value) {
    return '--'
  }
  return value
}

const timeFormatter = (value, format = 'YYYY-MM-DD HH:mm:ss') => {
  if (!value) {
    return '--'
  }
  const formated = moment(value).format(format)
  if (formated === 'Invalid date') {
    return value
  }
  return formated
}

const numericFormatter = (value) => {
  if (isNaN(value) || value === null || value === undefined || value === '') {
    return '--'
  }
  return value
}

export function singlechar(value) {
  return defaultFormatter(value)
}

export function longchar(value) {
  return defaultFormatter(value)
}

export function int(value) {
  return numericFormatter(value)
}

export function float(value) {
  return numericFormatter(value)
}

export function date(value) {
  return timeFormatter(value, 'YYYY-MM-DD')
}

export function time(value) {
  return timeFormatter(value, 'YYYY-MM-DD HH:mm:ss')
}

export function objuser(value) {
  if (!value) {
    return '--'
  }
  const userList = window.CMDB_USER_LIST || []
  const user = userList.find(user => user.english_name === value)
  if (user) {
    return `${user.english_name}(${user.chinese_name})`
  }
  return value
}

export function timezone(value) {
  return defaultFormatter(value)
}

export function bool(value) {
  if (['true', 'false'].includes(String(value))) {
    return String(value)
  }
  return '--'
}

export function enumeration(value, options, showId = false) {
  const option = (options || []).find(option => option.id === value)
  if (!option) {
    return '--'
  }
  if (showId) {
    return `${option.name}(${option.id})`
  }
  return option.name
}

export function foreignkey(value) {
  if (Array.isArray(value)) {
    return value.map(inst => `${inst.bk_inst_name}[${inst.bk_inst_id}]`).join(',')
  }
  if (String(value).length) {
    return value
  }
  return '--'
}

export function list(value) {
  return defaultFormatter(value)
}

export function implode(value, separator = ',') {
  if (Array.isArray(value)) {
    return value.join(separator)
  }
  return value.toString()
}

export function array(value) {
  if (!value || (Array.isArray(value) && value.length === 0)) {
    return '--'
  }

  if (typeof value === 'string') {
    return value
  }

  // 字符串数组
  if (value.every(val => typeof val === 'string')) {
    return value.toString()
  }

  return object(value)
}

export function object(value) {
  if (!value) {
    return '--'
  }

  let result = '--'
  try {
    result = JSON.stringify(value)
  } catch (e) {
    result = '--'
  }

  return result
}

const formatterMap = {
  singlechar,
  longchar,
  int,
  float,
  date,
  time,
  objuser,
  timezone,
  bool,
  foreignkey,
  list,
  enum: enumeration,
  array,
  object
}

export default function formatter(value, property, options) {
  const isPropertyObject = typeof property === 'object'
  const type = isPropertyObject ? property.bk_property_type : property
  const propertyOptions = isPropertyObject ? property.option : options
  if (has(formatterMap, type)) {
    return formatterMap[type](value, propertyOptions)
  }
  return value
}
