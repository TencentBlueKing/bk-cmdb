/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

// import VeeValidate, {Validator} from 'vee-validate'
import VeeValidate, {Validator} from 'vee-validate'
/*
    名称
*/
const isName = {
    getMessage (field, args) { // 錯誤提示
        return '请输入合法的' + field
    },
    validate: value => { // 驗證規則
        return /^([a-zA-Z0-9_ ]|[\u4e00-\u9fa5]|[\uac00-\ud7ff]|[\u0800-\u4e00]){1,20}$/.test(value)
    }
}
Validator.extend('name', isName)
/*
    id
*/
const isID = {
    getMessage (field, args) { // 错误提示
        // zh_CN: field => '请输入' + field
        return '请输入' + field
    },
    validate: value => { // 验证规则
        var reg = /^[a-z0-9_]{1,20}$/
        return reg.test(value)
    }
}
Validator.extend('id', isID)

const isHttp = {
    getMessage (field, args) { // 错误提示
        // zh_CN: field => '请输入' + field
        return '请输入' + field
    },
    validate: value => { // 验证规则
        var reg = /^http:\/\/[^\s]+/
        return reg.test(value)
    }
}
Validator.extend('http', isHttp)
/*
    长字符
*/
const longchar = {
    getMessage (field, args) { // 错误提示
        return '请输入' + field
    },
    validate: value => { // 验证规则
        return /^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[\(\)\+\-《》_,，；;“”‘’。\."'\\\/: ]){0,2000}$/.test(value)
    }
}
Validator.extend('longchar', longchar)

/*
    短字符
*/
const singlechar = {
    getMessage (field, args) { // 错误提示
        return '请输入' + field
    },
    validate: value => { // 验证规则
        return /^([a-zA-Z0-9]|[\u4e00-\u9fa5]|[\(\)\+\-《》_,，；;“”‘’。\."'\\\/: ]){0,256}$/.test(value)
    }
}
Validator.extend('singlechar', singlechar)

/*
    更新Validator
*/
const dictionary = {
    zh_CN: {
        messages: {
            name: () => '请输入正确的内容',
            longchar: () => '请输入正确的内容',
            singlechar: () => '请输入正确的内容',
            id: () => '格式不正确，只能包含下划线，数字，英文小写',
            http: () => '请输入以http://开头的URL',
            required: (field) => '请输入' + field,
            numeric: (field) => '请输入数字',
            regex: (field) => field + '不合法'
        },
        attributes: {
            // name: '中文名',
            // name: '英文名'
            // lifeCycle: '生命周期',
            // maintainers: '运维人员',
            // strictCharacter: '严格字符',
            // looseCharacter: '宽松字符',
            // number: '数字',
            // ipAddress: 'ip地址',
            // commonName: '名称'
        }
    },
    en: {
        messages: {
            name: () => 'Please enter the correct content',
            longchar: () => 'Please enter the correct content',
            singlechar: () => 'Please enter the correct content',
            id: () => 'The format is incorrect and can only contain underscores, numbers, and lowercase English',
            http: () => 'Please enter a URL beginning with http://',
            // lifeCycle: () => '请输入格式为 yyyy-MM-dd 的生命周期',
            // maintainers: () => '请选择运维人员',
            // strictCharacter: () => '请输入正确的验证字符',
            // looseCharacter: () => '请输入正确的宽松字符',
            // number: () => '请输入正确的数字',
            // ipAddress: () => '请输入正确的IP地址',
            // commonName: () => '请输入正确的名称',
            required: (field) => 'Please enter ' + field,
            numeric: (field) => 'Please enter the number',
            regex: (field) => field + ' is illegal'
        },
        attributes: {
            Description: 'Description',
            Name: 'Name'
            // lifeCycle: '生命周期',
            // maintainers: '运维人员',
            // strictCharacter: '严格字符',
            // looseCharacter: '宽松字符',
            // number: '数字',
            // ipAddress: 'ip地址',
            // commonName: '名称'
        }
    }
}
Validator.localize(dictionary)

export default dictionary