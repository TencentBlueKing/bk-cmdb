/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */

<template lang="html">
    <div class="allField">
        <v-base-info ref="baseInfo"
        :isShow="isShow"
        :objId="objId"
        :classificationId='classificationId'
        :associationId="associationId"
        :type="type"
        :isReadOnly="isModelDetailReadOnly"
        :isMainLine="isMainLine"
        @baseInfoSuccess="baseInfoSuccess"
        @confirm="baseInfoSuccess"
        @cancel="cancel">
        </v-base-info>
        <div class="tab-content model-field-content pb20" v-show="type==='change'">
            <div class="add-field clearfix">
                <bk-button type="primary" :title="$t('ModelManagement[\'新增字段\']')" @click="addField" v-if="!isReadOnly">
                    {{$t('ModelManagement["新增字段"]')}}
                </bk-button>
                <div class="btn-group" :class="{'disabled': isReadOnly}">
                    <bk-button type="default" :title="$t('ModelManagement[\'导入\']')" :disabled="isReadOnly" class="btn mr10">
                        {{$t('ModelManagement["导入"]')}}
                        <input v-if="!isReadOnly" ref="fileInput" type="file" @change.prevent="handleFile">
                    </bk-button>
                    <form :action="exportUrl" method="POST" class="form">
                        <bk-button type="default" btnType="submit" :title="$t('ModelManagement[\'导出\']')" class="btn">
                            {{$t('ModelManagement["导出"]')}}
                        </bk-button>
                    </form>
                </div>
            </div>
            <div class="table-content" v-bkloading="{isLoading: isLoading}">
                <div class="title-content">
                    <ul v-show="fieldList.length>0">
                        <li>{{$t('ModelManagement["唯一"]')}}</li>
                        <li>{{$t('ModelManagement["必填字段"]')}}</li>
                        <li>{{$t('ModelManagement["类型"]')}}</li>
                        <li>{{$t('ModelManagement["字段名"]')}}</li>
                        <li>{{$t('ModelManagement["操作"]')}}</li>
                    </ul>
                </div>
                <div class="list-content-wrapper">
                    <form id="validate-form-change">
                        <div class="list-content" v-for="(item, index) in fieldList" :class="{'editable':item['ispre'] || isReadOnly}">
                            <ul @click="toggleDetailShow(item, index)">
                                <li><i class=" fb bk-icon icon-check-1" v-show="item['isonly']"></i></li>
                                <li><i class=" fb bk-icon icon-check-1" v-show="item['isrequired']"></i></li>
                                <li>{{formatFieldType(item['bk_property_type'])}}</li>
                                <li :title="`${item['bk_property_name']}(${item['bk_property_id']})`">{{item['bk_property_name']}}({{item['bk_property_id']}})</li>
                                <li>
                                    <div class="btn-contain" v-if="item['bk_property_id']==='InstName'" v-bktooltips="{
                                            isShow: tips.innerField.isShow,
                                            content: tips.innerField.content,
                                            direction: tips.innerField.direction
                                        }">
                                        <i class="icon-cc-del f14 vm editable"></i>
                                    </div>
                                    <div class="btn-contain" v-else>
                                        <i class="icon-cc-del f14 vm"  v-if="!item['ispre'] && !isReadOnly"
                                            :class="{'editable':item['ispre'] || isReadOnly}"
                                            @click.stop="showConfirmDialog('delete',item, {id:item['id'], index:index})"
                                        ></i>
                                    </div>
                                </li>
                            </ul>
                            <div class="list-content-hidden" v-show="item.isShow">
                                <form class="form-common clearfix">
                                    <div class="clearfix mb30">
                                        <h3>{{$t('ModelManagement["字段配置"]')}}</h3>
                                        <div class="form-common-item" :class="{'disabled': isReadOnly}">
                                            <label class="form-common-label">{{$t('ModelManagement["中文名"]')}}<span class=""> * </span></label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" :disabled="isReadOnly" class="from-input" name="" :placeholder="$t('ModelManagement[\'请输入字段名称\']')" v-model.trim="curFieldInfo['bk_property_name']"
                                                maxlength="15"
                                                data-parsley-required="true"
                                                :data-parsley-required-message="$t('ModelManagement[\'该字段是必填项\']')"
                                                data-parsley-maxlength="20"
                                                :data-parsley-pattern="reg"
                                                :data-parsley-pattern-message="$t('ModelManagement[\'包含了非法字符\']')"
                                                data-parsley-trigger="input blur"
                                                >
                                            </div>
                                        </div>
                                        <div class="form-common-item disabled tr">
                                            <label class="form-common-label">{{$t('ModelManagement["英文名"]')}}<span class=""> * </span></label>
                                            <div class="form-common-content interior-width-control tl">
                                                <input type="text" disabled class="from-input" name="" value="" :placeholder="$t('ModelManagement[\'下划线/数字/字母\']')" v-model.trim="item['bk_property_id']">
                                            </div>
                                        </div>
                                        <div class="form-common-item tr" :class="{'disabled': isReadOnly}">
                                            <label class="form-common-label">{{$t('ModelManagement["单位"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" class="from-input" name="" :placeholder="$t('ModelManagement[\'请输入单位\']')"
                                                :disabled="isReadOnly"
                                                v-model.trim="curFieldInfo['unit']">
                                            </div>
                                        </div>
                                        <div class="form-common-item block mt20" :class="{'disabled': isReadOnly}">
                                            <label class="form-common-label">{{$t('ModelManagement["提示语"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" :disabled="isReadOnly" class="from-input" name="" :placeholder="$t('ModelManagement[\'请输入提示语\']')" v-model.trim="curFieldInfo['placeholder']">
                                            </div>
                                        </div>
                                    </div>
                                    <!-- 数字 -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'int'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model.trim="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isrequired']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text">{{$t('ModelManagement["是否唯一"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isonly']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                        <div class="form-common-item mt20" :class="{'disabled': isReadOnly}">
                                            <label class="form-common-label">{{$t('ModelManagement["最小值"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" maxlength="11" class="from-input" name="" :placeholder="$t('ModelManagement[\'请输入最小值\']')" v-model.trim="item.option.min" v-if="item.option" :disabled="isReadOnly">
                                                <span class="error-msg" v-show="isIntErrorShow.min">{{$t('Common["内容不合法"]')}}</span>
                                            </div>
                                        </div>
                                        <div class="form-common-item mt20 ml10" :class="{'disabled': isReadOnly}">
                                            <label class="form-common-label">{{$t('ModelManagement["最大值"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" maxlength="11" class="from-input" name="" :placeholder="$t('ModelManagement[\'请输入最大值\']')" v-model.trim="item.option.max" v-if="item.option" :disabled="isReadOnly">
                                                <span class="error-msg" v-show="isIntErrorShow.max">{{$t('Common["内容不合法"]')}}</span>
                                            </div>
                                        </div>
                                    </div>
                                    <!-- 长字符 -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'longchar'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isrequired']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text">{{$t('ModelManagement["是否唯一"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isonly']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                        <div class="form-common-item mt20" :class="{'disabled': isReadOnly}">
                                            <label class="form-common-label">{{$t('Common["正则验证"]')}}</label>
                                            <div class="form-common-content reg-verification ">
                                                <input type="text" class="from-input" name="" placeholder="" v-model.trim="item.option" :disabled="isReadOnly">
                                            </div>
                                        </div>
                                    </div>
                                    <!-- 短字符 -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'singlechar'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isrequired']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text">{{$t('ModelManagement["是否唯一"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isonly']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                        <div class="form-common-item mt20" :class="{'disabled': isReadOnly}">
                                            <label class="form-common-label">{{$t('Common["正则验证"]')}}</label>
                                            <div class="form-common-content reg-verification ">
                                                <input type="text" class="from-input" name="" placeholder="" v-model.trim="item.option" :disabled="isReadOnly">
                                            </div>
                                        </div>
                                    </div>
                                    <!-- 枚举 -->
                                    <div class="mt20 clearfix" v-if="item['bk_property_type'] === 'enum'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                        <div class="enum-table" :class="{'disabled':item['ispre'] || isReadOnly}">
                                            <div v-if="item.isShow">
                                                <div class="form-enum-wrapper" v-for="(field, fieldIndex) in item.option.list">
                                                    <span class="span-enum-radio" @click="item.option.defaultIndex = fieldIndex" :title="$t('ModelManagement[\'设置为默认值\']')" :class="{'active': fieldIndex === item.option.defaultIndex}"></span>
                                                    <div class="enum-id">
                                                        <input type="text" :placeholder="$t('ModelManagement[\'请输入ID\']')"
                                                            v-model.trim="field.id"
                                                            maxlength="15"
                                                            data-parsley-required="true"
                                                            :data-parsley-required-message="$t('ModelManagement[\'该字段是必填项\']')"
                                                            data-parsley-pattern="^[a-zA-Z0-9_]{1,20}$"
                                                            :data-parsley-pattern-message="$t('ModelManagement[\'包含了非法字符\']')"
                                                            data-parsley-trigger="blur"
                                                            data-parsley-no-repeat="changeId"
                                                            @input="forceUpdate('newId')"
                                                        >
                                                    </div>
                                                    <div class="enum-name">
                                                        <input type="text" :placeholder="$t('ModelManagement[\'请输入名称英文数字\']')"
                                                            v-model.trim="field.name"
                                                            maxlength="15"
                                                            data-parsley-required="true"
                                                            :data-parsley-required-message="$t('ModelManagement[\'该字段是必填项\']')"
                                                            data-parsley-maxlength="15"
                                                            :data-parsley-pattern="reg"
                                                            :data-parsley-pattern-message="$t('ModelManagement[\'包含了非法字符\']')"
                                                            data-parsley-trigger="blur"
                                                            data-parsley-no-repeat="change"
                                                            :data-parsley-errors-container="'#changeEnumError'+fieldIndex"
                                                            @input="forceUpdate('change')"
                                                        >
                                                        <!-- 表单验证错误信息容器 -->
                                                        <div class="form-enum-error" :id="'changeEnumError'+fieldIndex"></div>
                                                    </div>
                                                    <button class="bk-icon"
                                                        :disabled="item.option.list.length === 1"
                                                        @click.prevent="deleteEnum('change',fieldIndex,index)"
                                                    ><i class="icon-cc-del"></i></button>
                                                    <button class="bk-icon icon-plus" @click.prevent="addEnum('change',fieldIndex,index)" v-if="fieldIndex === (item.option.list.length -1)"></button>
                                                    <!-- 拖拽标识点，暂未实现，隐藏 -->
                                                    <i class="form-enum-wrapper-dot" hidden></i>
                                                </div>
                                            </div>
                                            <div class="enum-disabled" v-if="isReadOnly"></div>
                                        </div>
                                    </div>

                                    <!-- 日期 -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'date'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isrequired']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                    </div>
                                    <!-- 时间 -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'time'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isrequired']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                    </div>
                                    <!-- 单关联 -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'singleasst'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                        <div class="form-common-item selcet-width-control mt20" :class="{'disabled':item['ispre'] || isReadOnly}">
                                            <label class="form-common-label">{{$t('ModelManagement["关联模型"]')}}</label>
                                            <div class="form-common-content">
                                                <bk-select
                                                    disabled
                                                    :selected="curModelType"
                                                    @on-selected="modelChange">
                                                    <bk-option-group
                                                        v-for="(group, groupIndex) of modelList"
                                                        :label="group['bk_classification_name']"
                                                        :key="groupIndex">
                                                        <bk-select-option
                                                            v-for="(option, optionIndex) of group['bk_objects']"
                                                            :key="optionIndex"
                                                            :value="option['bk_obj_id']"
                                                            :label="option['bk_obj_name']">
                                                        </bk-select-option>
                                                    </bk-option-group>
                                                </bk-select>
                                            </div>
                                        </div>
                                    </div>
                                    <!-- 多关联 -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'multiasst'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30 correlate-single-control">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                        <div class="form-common-item mt20" :class="{'disabled':item['ispre'] || isReadOnly}">
                                            <label class="form-common-label">{{$t('ModelManagement["关联模型"]')}}</label>
                                            <div class="form-common-content selcet-width-control">
                                                <bk-select
                                                    disabled
                                                    :selected="curModelType"
                                                    @on-selected="modelChange">
                                                    <bk-option-group
                                                        v-for="(group, groupIndex) of modelList"
                                                        :label="group['bk_classification_name']"
                                                        :key="groupIndex">
                                                        <bk-select-option
                                                            v-for="(option, optionIndex) of group['bk_objects']"
                                                            :key="optionIndex"
                                                            :value="option['bk_obj_id']"
                                                            :label="option['bk_obj_name']">
                                                        </bk-select-option>
                                                    </bk-option-group>
                                                </bk-select>
                                            </div>
                                        </div>
                                    </div>
                                    <!-- 用户 -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'objuser'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isrequired']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                    </div>
                                    <!-- 时区 -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'timezone'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['isrequired']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                    </div>
                                    <!-- bool -->
                                    <div class="mt20 clearfix" v-show="item['bk_property_type'] === 'bool'">
                                        <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                        <div class="form-common-item disabled">
                                            <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                            <div class="form-common-content interior-width-control">
                                                <input type="text" disabled class="from-input" name="" placeholder="" :value="formatFieldType(item['bk_property_type'])">
                                            </div>
                                        </div>
                                        <div class="form-common-item form-common-item2 pl30">
                                            <div class="from-selcet-wrapper mr30">
                                                <label class="bk-form-checkbox bk-checkbox-small">
                                                    <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                    <input type="checkbox" name="checkbox1" v-model="curFieldInfo['editable']" :disabled="item['ispre'] || isReadOnly">
                                                </label>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="submit-btn" v-if="!isReadOnly">
                                        <bk-button type="primary" :loading="$loading('saveChange')" class="save-btn main-btn mr10" :class="{'loading': $loading('saveChange')}" @click="saveFieldChange(item, index)">
                                            {{$t('Common["保存"]')}}
                                        </bk-button>
                                        <bk-button type="default" class="cancel-btn vice-btn" @click="cancelFieldChange(item, index)">
                                            {{$t('Common["取消"]')}}
                                        </bk-button>
                                    </div>
                                </form>
                            </div>
                        </div>
                    </form>
                </div>
            </div>
            <!-- <div class="add-field-detail" v-show="isAddFieldShow || (!isCreateField && !fieldList.length)"> -->
            <!-- 新增字段 -->
            <form id="validate-form-new">
                <div class="add-field-wrapper" v-show="isAddFieldShow">
                    <div class="add-field-detail">
                        <div class="bg-titel" @click="closeAddFieldBox"><img src="../../../common/images/down_icon.png" alt="" ></div>
                        <div class="border-control">
                            <form class="form-common clearfix">
                                <div class="clearfix mb30">
                                    <h3>{{$t('ModelManagement["字段配置"]')}}</h3>
                                    <div class="form-common-item tl">
                                        <label class="form-common-label">{{$t('ModelManagement["中文名"]')}}<span class=""> * </span></label>
                                        <div class="form-common-content interior-width-control">
                                            <input type="text" class="from-input" name="" placeholder="" v-model.trim="newFieldInfo.propertyName"
                                            maxlength="15"
                                            data-parsley-required="true"
                                            :data-parsley-required-message="$t('ModelManagement[\'该字段是必填项\']')"
                                            data-parsley-maxlength="20"
                                            :data-parsley-pattern="reg"
                                            :data-parsley-pattern-message="$t('ModelManagement[\'包含了非法字符\']')"
                                            data-parsley-trigger="input blur"
                                            >
                                        </div>
                                    </div>
                                    <div class="form-common-item tr pr">
                                        <label class="form-common-label">{{$t('ModelManagement["英文名"]')}}<span class=""> * </span></label>
                                        <div class="form-common-content interior-width-control tl">
                                            <input type="text" class="from-input" name="" value="" v-model.trim="newFieldInfo.propertyId"
                                            maxlength="20"
                                            data-parsley-required="true"
                                            :data-parsley-required-message="$t('ModelManagement[\'该字段是必填项\']')"
                                            data-parsley-maxlength="20"
                                            data-parsley-pattern="^[a-zA-Z0-9_]{1,20}$"
                                            :data-parsley-pattern-message="$t('ModelManagement[\'必须以英文开头，由英文、数字及下划线组成\']')"
                                            data-parsley-trigger="input blur"
                                            >
                                        </div>
                                        <i class="icon-info-png" v-tooltip="$t('ModelManagement[\'下划线/数字/字母\']')"></i>
                                    </div>
                                    <div class="form-common-item tr">
                                        <label class="form-common-label">{{$t('ModelManagement["单位"]')}}</label>
                                        <div class="form-common-content interior-width-control">
                                            <input type="text" class="from-input" name="" value="" :placeholder="$t('ModelManagement[\'请输入单位\']')" v-model.trim="newFieldInfo.unit">
                                        </div>
                                    </div>
                                    <div class="form-common-item mt20 block">
                                        <label class="form-common-label">{{$t('ModelManagement["提示语"]')}}</label>
                                        <div class="form-common-content interior-width-control">
                                            <input type="text" class="from-input" name="" value="" :placeholder="$t('ModelManagement[\'请输入提示语\']')" v-model.trim="newFieldInfo.placeholder">
                                        </div>
                                    </div>
                                </div>
    
                                <h3>{{$t('ModelManagement["选项"]')}}</h3>
                                <div class="clearfix">
                                    <div class="form-common-item mr0">
                                        <label class="form-common-label">{{$t('ModelManagement["类型"]')}}</label>
                                        <div class="form-common-content interior-width-control">
                                            <div class="select-content tc">
                                                <bk-select
                                                    :selected.sync="newFieldInfo.propertyType"
                                                    @on-selected="fieldTypeChange">
                                                    <bk-select-option
                                                        v-for="(option, index) of fieldTypeList"
                                                        :key="option.value"
                                                        :value="option.value"
                                                        :label="option.label">
                                                    </bk-select-option>
                                                </bk-select>
                                            </div>
                                        </div>
                                    </div>
                                    <div class="form-common-item form-common-item2 pl30">
                                        <div class="from-selcet-wrapper mr30">
                                            <label class="bk-form-checkbox bk-checkbox-small">
                                                <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否可编辑"]')}}</i>
                                                <input type="checkbox" name="checkbox1" v-model="newFieldInfo.editable">
                                            </label>
                                        </div>
                                        <div class="from-selcet-wrapper mr30" v-if="isShowRequired(newFieldInfo.propertyType)">
                                            <label class="bk-form-checkbox bk-checkbox-small">
                                                <i class="bk-checkbox-text mr5">{{$t('ModelManagement["是否必填"]')}}</i>
                                                <input type="checkbox" name="checkbox1" v-model="newFieldInfo.isRequired">
                                            </label>
                                        </div>
                                        <div class="from-selcet-wrapper" v-if="isShowOnly(newFieldInfo.propertyType)">
                                            <label class="bk-form-checkbox bk-checkbox-small">
                                                <i class="bk-checkbox-text">{{$t('ModelManagement["是否唯一"]')}}</i>
                                                <input type="checkbox" name="checkbox1" v-model="newFieldInfo['isonly']">
                                            </label>
                                        </div>
                                    </div>
                                </div>

                                <!-- 数字 -->
                                <div class="clearfix" v-show="newFieldInfo.propertyType === 'int'">
                                    <div class="form-common-item mt20">
                                        <label class="form-common-label">{{$t('ModelManagement["最小值"]')}}</label>
                                        <div class="form-common-content interior-width-control">
                                            <input type="text" maxlength="11" class="from-input" name="" :placeholder="$t('ModelManagement[\'请输入最小值\']')" v-model.trim="newFieldInfo.option.min">
                                            <span class="error-msg" v-show="isIntErrorShow.min">{{$t('Common["内容不合法"]')}}</span>
                                        </div>
                                    </div>
                                    <div class="form-common-item  mt20 tr">
                                        <label class="form-common-label">{{$t('ModelManagement["最大值"]')}}</label>
                                        <div class="form-common-content interior-width-control tl">
                                            <input type="text" maxlength="11" class="from-input" name="" :placeholder="$t('ModelManagement[\'请输入最大值\']')" v-model.trim="newFieldInfo.option.max">
                                            <span class="error-msg" v-show="isIntErrorShow.max">{{$t('Common["内容不合法"]')}}</span>
                                        </div>
                                    </div>
                                </div>
                                <!-- 长字符 -->
                                <div class="clearfix" v-show="newFieldInfo.propertyType === 'longchar'">
                                    <div class="form-common-item mt20">
                                        <label class="form-common-label">{{$t('Common["正则验证"]')}}</label>
                                        <div class="form-common-content reg-verification ">
                                            <input type="text" class="from-input" name="" placeholder="" v-model.trim="newFieldInfo.option">
                                        </div>
                                    </div>
                                </div>
                                <!-- 短字符 -->
                                <div class="clearfix" v-show="newFieldInfo.propertyType === 'singlechar'">
                                    <div class="form-common-item mt20">
                                        <label class="form-common-label">{{$t('Common["正则验证"]')}}</label>
                                        <div class="form-common-content reg-verification ">
                                            <input type="text" class="from-input" name="" placeholder="" v-model.trim="newFieldInfo.option">
                                        </div>
                                    </div>
                                </div>
                                <!-- 枚举 -->
                                <div class="clearfix form-option" v-if="newFieldInfo.propertyType === 'enum'">
                                    
                                    <!-- <div v-pre class="clearfix"></div> -->
                                    <div class="form-enum-box clearfix" v-if="newFieldInfo.propertyType === 'enum'">
                                        <div class="form-enum-wrapper" v-for="(field, fieldIndex) in newFieldInfo.option.list">
                                            <span class="span-enum-radio" @click="newFieldInfo.option.defaultIndex = fieldIndex" title="设置为默认值" :class="{'active': fieldIndex === newFieldInfo.option.defaultIndex}"></span>
                                            <div class="enum-id">
                                                <input type="text" :placeholder="$t('ModelManagement[\'请输入ID\']')"
                                                    v-model.trim="field.id"
                                                    maxlength="15"
                                                    data-parsley-required="true"
                                                    :data-parsley-required-message="$t('ModelManagement[\'该字段是必填项\']')"
                                                    data-parsley-pattern="^[a-zA-Z0-9_]{1,20}$"
                                                    :data-parsley-pattern-message="$t('ModelManagement[\'包含了非法字符\']')"
                                                    data-parsley-trigger="blur"
                                                    data-parsley-no-repeat="newId"
                                                    @input="forceUpdate('newId')"
                                                >
                                            </div>
                                            <div class="enum-name">
                                                <input type="text" :placeholder="$t('ModelManagement[\'请输入名称英文数字\']')"
                                                    v-model.trim="field.name"
                                                    maxlength="15"
                                                    data-parsley-required="true"
                                                    :data-parsley-required-message="$t('ModelManagement[\'该字段是必填项\']')"
                                                    data-parsley-maxlength="15"
                                                    :data-parsley-pattern="reg"
                                                    :data-parsley-pattern-message="$t('ModelManagement[\'包含了非法字符\']')"
                                                    data-parsley-trigger="blur"
                                                    data-parsley-no-repeat="new"
                                                    :data-parsley-errors-container="'#newEnumError'+fieldIndex"
                                                    @input="forceUpdate('new')"
                                                >
                                                <div class="form-enum-error" :id="'newEnumError'+fieldIndex"></div>
                                            </div>
                                            <button class="bk-icon"
                                                :disabled="newFieldInfo.option.list.length === 1"
                                                @click.prevent="deleteEnum('new',fieldIndex)"
                                            ><i class="icon-cc-del"></i></button>
                                            <button class="bk-icon icon-plus" @click.prevent="addEnum('new',fieldIndex)" v-if="fieldIndex === (newFieldInfo.option.list.length -1)"></button>
                                            <!-- 表单验证错误信息容器 -->
                                            <!-- 拖拽标识点，暂未实现，隐藏 -->
                                            <i class="form-enum-wrapper-dot" hidden></i>
                                        </div>
                                    </div>
                                    <div class="select-error tc" v-if="isEnumErrorShow&&!newFieldInfo.option.list.length">{{$t('ModelManagement["请先设置枚举内容"]')}}</div>
                                </div>
                                <!-- 单关联 -->
                                <div class="clearfix" v-show="newFieldInfo.propertyType === 'singleasst'">
                                    <div class="form-common-item mt20">
                                        <label class="form-common-label">{{$t('ModelManagement["关联模型"]')}}</label>
                                        <div class="form-common-content selcet-width-control">
                                            <bk-select
                                                ref="singleasstSelect"
                                                :selected="''"
                                                @on-selected="modelSelected">
                                                <bk-option-group
                                                    v-for="(group, groupIndex) of modelList"
                                                    :label="group['bk_classification_name']"
                                                    :key="groupIndex">
                                                    <bk-select-option
                                                        v-for="(option, optionIndex) of group['bk_objects']"
                                                        :key="option['bk_obj_id']"
                                                        :value="option['bk_obj_id']"
                                                        :label="option['bk_obj_name']">
                                                    </bk-select-option>
                                                </bk-option-group>
                                            </bk-select>
                                            <span class="select-error" v-if="isSelectErrorShow">{{$t('ModelManagement["请选择关联模型"]')}}</span>
                                        </div>
                                    </div>
                                </div>
                                <!-- 多关联 -->
                                <div class="clearfix" v-show="newFieldInfo.propertyType === 'multiasst'">
                                    <div class="form-common-item mt20">
                                        <label class="form-common-label">{{$t('ModelManagement["关联模型"]')}}</label>
                                        <div class="form-common-content selcet-width-control tc">
                                            <bk-select
                                                ref="multiasstSelect"
                                                :selected="''"
                                                @on-selected="modelSelected">
                                                <bk-option-group
                                                    v-for="(group, groupIndex) of modelList"
                                                    :label="group['bk_classification_name']"
                                                    :key="groupIndex">
                                                    <bk-select-option
                                                        v-for="(option, oIndex) of group['bk_objects']"
                                                        :key="option['bk_obj_id']"
                                                        :value="option['bk_obj_id']"
                                                        :label="option['bk_obj_name']">
                                                    </bk-select-option>
                                                </bk-option-group>
                                            </bk-select>
                                            <span class="select-error" v-if="isSelectErrorShow">{{$t('ModelManagement["请选择关联模型"]')}}</span>
                                        </div>
                                    </div>
                                </div>
                                <!-- 保存取消按钮 -->
                                <div class="button-wraper">
                                    <bk-button type="primary" class="save-btn main-btn mr10" :loading="$loading('saveNew')" @click="saveNewField">
                                        {{$t('Common["保存"]')}}
                                    </bk-button>
                                    <bk-button type="default" class="cancel-btn vice-btn" @click="closeAddFieldBox">
                                        {{$t('Common["取消"]')}}
                                    </bk-button>
                                </div>
                            </form>
                        </div>
                    </div>
                </div>
            </form>
        </div>
    </div>
</template>

<script type="text/javascript">
    import $ from 'jquery'
    import Parsley from 'parsleyjs'
    import vBaseInfo from './baseInfo'
    import '@/common/js/parsley_locale'
    import {mapGetters} from 'vuex'
    export default {
        props: {
            isShow: {
                default: false
            },
            associationId: {
                default: 0
            },
            type: {
                default: Boolean
            },
            isMainLine: {
                default: false
            },
            isModelDetailReadOnly: {
                default: false
            },
            isReadOnly: {
                default: false,
                type: Boolean
            },
            id: {               // 模型ID
                default: 0
            },
            objId: {
                default: ''
            },
            isCreateField: {
                default: true       // 是否处于新增字段状态
            },
            classificationId: {
                default: 0
            }
        },
        components: {
            vBaseInfo
        },
        watch: {
            'newFieldInfo.isonly' (isonly) {
                if (isonly) {
                    this.newFieldInfo.isRequired = true
                }
            },
            'newFieldInfo.isRequired' (isRequired) {
                if (!isRequired) {
                    this.newFieldInfo.isonly = false
                }
            },
            'curFieldInfo.isonly' (isonly) {
                if (isonly) {
                    this.curFieldInfo.isrequired = true
                }
            },
            'curFieldInfo.isrequired' (isrequired) {
                if (!isrequired) {
                    this.curFieldInfo.isonly = false
                }
            },
            'newFieldInfo.option': {
                handler (newOption, oldOption) {
                    if (this.newFieldInfo.propertyType === 'int') {
                        this.isIntErrorShow.min = false
                        this.isIntErrorShow.max = false
                    }
                },
                deep: true
            },
            objId () {
                if (this.objId === '') {
                    this.$refs.baseInfo.clearData()
                } else {
                    this.$refs.baseInfo.getBaseInfo(this.objId)
                }
                // this.$refs.baseInfo.clearData()
            },
            'newFieldInfo.propertyType' (val) {
                if (val === 'singleasst') {
                    this.$refs.singleasstSelect.curLabel = ''
                    this.$refs.singleasstSelect.curValue = ''
                    this.$refs.singleasstSelect.model = ''
                } else if (val === 'multiasst') {
                    this.$refs.multiasstSelect.curLabel = ''
                    this.$refs.multiasstSelect.curValue = ''
                    this.$refs.multiasstSelect.model = ''
                }
            },
            language (lang) {
                if (lang === 'zh_cn') {
                    this.fieldTypeList = this.fieldTypeListForZh
                } else {
                    this.fieldTypeList = this.fieldTypeListForEn
                }
            }
        },
        data () {
            return {
                reg: '^([a-zA-Z0-9_]|[\u4e00-\u9fa5]|[()+-《》,，；;“”‘’。."\' \\/:]){1,15}$',
                isSelectErrorShow: false,       // 关联模型为空时的提示状态
                isEnumErrorShow: false,         // 枚举内容为空是的提示状态
                isIntErrorShow: {
                    min: false,
                    max: false
                },
                tips: {
                    innerField: {
                        isShow: false,
                        direction: 'top',
                        content: '内置字段不可删除'
                    }
                },
                isLoading: false,           // 是否处于加载列表状态
                fieldList: [],          // 字段配置列表
                defaultModel: '',
                curFieldInfo: {         // 当前改动项
                    isonly: false,
                    isrequired: false
                },
                curFieldInfoCopy: {
                    isonly: false,
                    isrequired: false
                },
                newFieldInfo: {
                    propertyName: '',       // 字段名称
                    propertyId: '',         // API标识
                    propertyType: 'singlechar',      // 字段类型
                    isRequired: false,      // 是否必填
                    isReadOnly: false,
                    editable: true,
                    propertyGroup: 'default',
                    isOnly: false,          // 是否唯一
                    enumList: [],            // 枚举列表
                    fieldType: {
                        int: {
                            min: '',
                            max: ''
                        },
                        longchar: {
                            reg: ''
                        },
                        singlechar: {
                            reg: ''
                        },
                        enum: {

                        },
                        singleasst: {
                            label: '',
                            value: ''
                        }
                    },
                    option: []
                },
                newFieldInfoCopy: {},
                modelList: [],          // 模型分类及附属模型信息列表
                curModelType: '',
                curIndex: 0,            // 当前展开项索引
                isAddFieldShow: false
            }
        },
        computed: {
            ...mapGetters([
                'bkSupplierAccount',
                'language'
            ]),
            fieldTypeList () {
                let list = [
                    {
                        value: 'singlechar',
                        label: this.$t('ModelManagement["短字符"]')
                    },
                    {
                        value: 'int',
                        label: this.$t('ModelManagement["数字"]')
                    },
                    {
                        value: 'enum',
                        label: this.$t('ModelManagement["枚举"]')
                    },
                    {
                        value: 'date',
                        label: this.$t('ModelManagement["日期"]')
                    },
                    {
                        value: 'time',
                        label: this.$t('ModelManagement["时间"]')
                    },
                    {
                        value: 'longchar',
                        label: this.$t('ModelManagement["长字符"]')
                    },
                    {
                        value: 'singleasst',
                        label: this.$t('ModelManagement["单关联"]')
                    },
                    {
                        value: 'multiasst',
                        label: this.$t('ModelManagement["多关联"]')
                    },
                    {
                        value: 'objuser',
                        label: this.$t('ModelManagement["用户"]')
                    },
                    {
                        value: 'timezone',
                        label: this.$t('ModelManagement["时区"]')
                    },
                    {
                        value: 'bool',
                        label: 'bool'
                    }
                ]
                if (this.isMainLine) {
                    list.splice(6, 2)
                }
                return list
            },
            exportUrl () {
                return `${window.siteUrl}object/owner/${this.bkSupplierAccount}/object/${this.objId}/export`
            },
            importUrl () {
                return `${window.siteUrl}object/owner/${this.bkSupplierAccount}/object/${this.objId}/import`
            }
        },
        methods: {
            isShowRequired (type) {
                switch (type) {
                    case 'singlechar':
                    case 'int':
                    case 'date':
                    case 'time':
                    case 'longchar':
                    case 'objuser':
                    case 'timezone':
                        return true
                    default:
                        return false
                }
            },
            isShowOnly (type) {
                switch (type) {
                    case 'singlechar':
                    case 'int':
                    case 'longchar':
                        return true
                    default:
                        return false
                }
            },
            addTableList (index, option) {
                if (this.isAddFieldShow) { // 新增
                    this.newFieldInfo.option.splice(index + 1, 0, {
                        list_header_name: '',
                        list_header_describe: '',
                        isEditDesc: false,
                        isEditName: false,
                        errorMsg: ''
                    })
                } else {
                    option.splice(index + 1, 0, {
                        list_header_name: '',
                        list_header_describe: '',
                        isEditDesc: false,
                        isEditName: false,
                        errorMsg: ''
                    })
                }
            },
            deleteTableList (index, option) {
                if (this.isAddFieldShow) {
                    this.newFieldInfo.option.splice(index, 1)
                } else {
                    option.splice(index, 1)
                }
            },
            listViewEdit (opt, view, isEdit) {
                this.$set(opt, view, isEdit)
                this.$forceUpdate()
            },
            isCloseConfirmShow () {
                // 校验字段
                if (this.isAddFieldShow) {
                    if (JSON.stringify(this.newFieldInfoCopy) !== JSON.stringify(this.newFieldInfo)) {
                        return true
                    }
                } else {
                    if (JSON.stringify(this.curFieldInfo) !== JSON.stringify(this.curFieldInfoCopy)) {
                        return true
                    }
                }
                // 校验模型名
                return this.$refs.baseInfo.isCloseConfirmShow()
            },
            handleFile (e) {
                this.isLoading = true
                let files = e.target.files
                let formData = new FormData()
                formData.append('file', files[0])
                this.$axios.post(this.importUrl, formData).then(res => {
                    if (res.result) {
                        let data = res.data[this.objId]
                        if (data.hasOwnProperty('insert_failed')) {
                            this.$alertMsg(data['insert_failed'][0])
                        } else if (data.hasOwnProperty('update_failed')) {
                            this.$alertMsg(data['update_failed'][0])
                        } else {
                            this.$alertMsg(this.$t('ModelManagement["导入成功"]'), 'success')
                            this.getModelField()
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.$refs.fileInput.value = ''
                    this.isLoading = false
                }).catch(reject => {
                    this.$alertMsg(this.$t('ModelManagement["导入失败"]'))
                    this.$refs.fileInput.value = ''
                    this.isLoading = false
                })
            },
            /*
                保存基本信息成功
            */
            baseInfoSuccess (obj) {
                this.$emit('baseInfoSuccess', obj)
            },
            /*
                取消按钮
            */
            cancel () {
                this.$emit('cancel')
            },
            inputOptionMin (item) {
                if (item.option.max !== '') {
                    if (parseInt(item.option.max) < parseInt(event.target.value)) {
                        item.option.max = event.target.value
                    }
                } else {
                    item.option.max = event.target.value
                    this.$set(item.option, 'max', event.target.value)
                }
            },
            inputOptionMax (item) {
                if (item.option.min !== '') {
                    if (parseInt(item.option.min) > parseInt(event.target.value)) {
                        item.option.min = event.target.value
                    }
                }
            },
            /*
                立即创建按钮
            */
            createField () {
                this.resetNewField()
                this.$emit('update:isCreateField', false)
                this.addField()
            },
            /*
                格式化显示字段类型
            */
            formatFieldType (type) {
                for (var i = 0; i < this.fieldTypeList.length; i++) {
                    if (this.fieldTypeList[i].value === type) {
                        return this.fieldTypeList[i].label
                    }
                }
            },
            /*
                新增时格式化选项内容
                type: 类型
                option: option的内容
            */
            formatFieldOption (type, option) {
                let opt = null
                switch (type) {
                    case 'int':
                        opt = {
                            min: option.min,
                            max: option.max
                        }
                        break
                    case 'enum':
                        option.list.map((item, index) => {
                            item['is_default'] = index === option.defaultIndex
                        })
                        opt = option.list
                        break
                    case 'longchar':
                    case 'singlechar':
                    case 'singleasst':
                    case 'multiasst':
                        opt = option
                        break
                }
                return opt
            },
            /*
                格式化获取到的字段选项
                item: 当前项
                index: 索引
            */
            parseFieldOption (item, index) {
                let option = null
                switch (item['bk_property_type']) {
                    case 'int':
                    case 'singleasst':
                    case 'multiasst':
                        if (item['Option'] !== 'undefined') {
                            option = item['Option']
                        }
                        break
                    case 'enum':
                        if (item['Option'] !== 'undefined') {
                            let opt = item['Option']
                            let defaultIndex = ''
                            for (let i = 0; i < opt.length; i++) {
                                if (opt[i].hasOwnProperty('is_default') && opt[i]['is_default']) {
                                    defaultIndex = i
                                    break
                                }
                            }
                            option = {
                                list: opt,
                                defaultIndex: defaultIndex
                            }
                        }
                        break
                    case 'longchar':
                    case 'singlechar':
                        option = item['Option']
                        break
                }
                this.fieldList[index].option = this.$deepClone(option)
            },
            formatAttrOption (data) {
                data.map(item => {
                    switch (item['bk_property_type']) {
                        case 'int':
                            if (item.option === null) {
                                item.option = {
                                    min: '',
                                    max: ''
                                }
                            }
                            break
                        case 'enum':
                            if (!Array.isArray(item.option) || (Array.isArray(item.option) && !item.option.length)) {
                                item.option = [{
                                    id: '',
                                    name: '',
                                    is_default: true
                                }]
                            }
                            break
                    }
                })
            },
            /*
                获取字段配置
            */
            getModelField () {
                let params = {
                    bk_obj_id: this.objId,
                    bk_supplier_account: this.bkSupplierAccount
                }
                this.isAddFieldShow = false
                this.isLoading = true
                this.$store.dispatch('object/getAttribute', {objId: this.objId, force: true}).then(res => {
                    if (res.result) {
                        this.formatAttrOption(res.data)
                        for (var i = 0; i < this.fieldList.length; i++) {
                            this.fieldList[i]['isShow'] = false
                        }
                        if (res.data.length) {
                            this.$emit('update:isCreateField', false)
                        } else {
                            this.$emit('update:isCreateField', true)
                        }
                        let arr = []
                        let empty = []         // 没有勾选唯一值和或者必须项时
                        let haveValue = []    // 有勾选唯一值和或者必须项时
                        this.fieldList = []
                        for (let item of res.data) {
                            // 解决后端变量与前端重名问题
                            item.Option = item.option
                            if (item['isonly'] && item['isrequired']) {
                                haveValue.unshift(item)
                            } else if (item['isonly']) {
                                haveValue.push(item)
                            } else if (item['isrequired']) {
                                haveValue.push(item)
                            } else if (item['isonly'] === false && item['isrequired'] === false) {
                                empty.push(item)
                            }
                            this.fieldList = haveValue.concat(empty)
                        }
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                    this.isLoading = false
                }).catch(reject => {
                    this.isLoading = false
                })
            },
            /*
                展开或隐藏字段详情
                item: 展开的当前项
                index: 索引
            */
            toggleDetailShow (item, index) {
                $('#validate-form-change').parsley().reset()
                this.isIntErrorShow.min = false
                this.isIntErrorShow.max = false
                if (!this.fieldList[index].isShow) {
                    this.parseFieldOption(item, index)
                    this.curFieldInfo['bk_property_name'] = item['bk_property_name']
                    this.curFieldInfo['isrequired'] = item['isrequired']
                    this.curFieldInfo['isonly'] = item['isonly']
                    this.curFieldInfo['editable'] = item['editable']
                    this.curFieldInfo['placeholder'] = item['placeholder']
                    this.curFieldInfo['unit'] = item['unit']
                    this.curFieldInfo['bk_asst_forward'] = ''
                } else {
                    this.curFieldInfo = {}
                }
                this.curFieldInfoCopy = this.$deepClone(this.curFieldInfo)
                for (var i = 0; i < this.fieldList.length; i++) {
                    if (index === i) {
                        this.fieldList[i].isShow = !this.fieldList[i].isShow
                        this.curIndex = index
                        // 处理单关联和多关联两种特殊情况
                        if (this.fieldList[i]['bk_property_type'] === 'singleasst' || this.fieldList[i]['bk_property_type'] === 'multiasst') {
                            this.curModelType = this.fieldList[i]['bk_asst_obj_id']
                        }
                    } else {
                        this.fieldList[i].isShow = false
                    }
                }
                this.fieldList.splice()
            },
            /*
                添加字段
            */
            addField () {
                this.resetNewField()
                this.isAddFieldShow = true
            },
            /*
                清空新增字段内容
            */
            resetNewField () {
                $('#validate-form-new').parsley().reset()
                this.isSelectErrorShow = false
                this.isEnumErrorShow = false
                this.newFieldInfo.propertyType = 'bool'
                this.$nextTick(() => {
                    this.newFieldInfo = {
                        propertyName: '',       // 字段名称
                        propertyId: '',         // API标识
                        propertyType: 'singlechar',      // 字段类型
                        isRequired: false,      // 是否必填
                        isReadOnly: false,
                        editable: true,
                        propertyGroup: 'default',
                        isOnly: false,          // 是否唯一
                        enumList: [],            // 枚举列表
                        fieldType: {
                            int: {
                                min: '',
                                max: ''
                            },
                            longchar: {
                                reg: ''
                            },
                            singlechar: {
                                reg: ''
                            },
                            enum: {
    
                            },
                            singleasst: {
                                type: ''
                            }
                        },
                        option: ''
                    }
                })
            },
            /*
                删除字段
                id: 需要删除项的ID
                index: 索引
            */
            deleteField (id, index) {
                this.$axios.delete('object/attr/' + id, {}).then(res => {
                    if (res.result) {
                        this.fieldList.splice(index, 1)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
                二次确认弹窗
                type: 类型
                params: 参数
            */
            showConfirmDialog (type, item, params) {
                if (item['ispre']) {
                    return
                }
                let self = this
                switch (type) {
                    case 'delete':
                        this.$bkInfo({
                            title: this.$t('ModelManagement["确定删除字段？"]'),
                            confirmFn () {
                                self.deleteField(params.id, params.index)
                            }
                        })
                        break
                }
            },
            /*
                验证新增配置字段是否为空
            */
            checkParams () {
                if (this.newFieldInfo.propertyType === 'singleasst' || this.newFieldInfo.propertyType === 'multiasst') {
                    if (this.newFieldInfo['bk_asst_obj_id'] === '') {
                        this.isSelectErrorShow = true
                        return false
                    }
                }
                if (this.newFieldInfo.propertyType === 'enum') {
                    if (this.newFieldInfo.option.length === 0) {
                        this.isEnumErrorShow = true
                        return false
                    }
                }
                if (this.newFieldInfo.propertyType === 'int') {
                    this.isIntErrorShow.min = !/^(-)?[0-9]*$/.test(this.newFieldInfo.option.min)
                    this.isIntErrorShow.max = !/^(-)?[0-9]*$/.test(this.newFieldInfo.option.max)
                    if (this.isIntErrorShow.min || this.isIntErrorShow.max) {
                        return false
                    }
                    if (parseInt(this.newFieldInfo.option.min) > parseInt(this.newFieldInfo.option.max)) {
                        this.isIntErrorShow.min = true
                        return false
                    }
                }
                this.isIntErrorShow.min = false
                this.isIntErrorShow.max = false
                this.isSelectErrorShow = false
                this.isEnumErrorShow = false
                return true
            },
            /*
                验证修改字段是否为空
            */
            checkChangeParams (item, index) {
                if (item['bk_property_type'] === 'singleasst' || item['bk_property_type'] === 'multiasst') {
                    if (!item['bk_asst_obj_id']) {
                        this.isSelectErrorShow = true
                        return false
                    }
                }
                if (item['bk_property_type'] === 'enum') {
                    if (item.option.list.length === 0) {
                        this.isEnumErrorShow = true
                        return false
                    }
                }
                if (item['bk_property_type'] === 'int') {
                    this.isIntErrorShow.min = !/^(-)?[0-9]*$/.test(item.option.min)
                    this.isIntErrorShow.max = !/^(-)?[0-9]*$/.test(item.option.max)
                    if (this.isIntErrorShow.min || this.isIntErrorShow.max) {
                        return false
                    }
                    if (parseInt(item.option.min) > parseInt(item.option.max)) {
                        this.isIntErrorShow.min = true
                        return false
                    }
                }
                this.isIntErrorShow.min = false
                this.isIntErrorShow.max = false
                this.isSelectErrorShow = false
                this.isEnumErrorShow = false
                return true
            },
            /*
                新增字段确认按钮
            */
            saveNewField () {
                $('#validate-form-new').parsley().validate()
                if (!$('#validate-form-new').parsley().isValid()) return
                if (!this.checkParams()) {
                    return
                }
                let params = {
                    creator: window.userName,
                    isonly: this.newFieldInfo['isonly'],
                    isreadonly: false,
                    isrequired: this.newFieldInfo.isRequired,
                    bk_property_group: 'default',
                    bk_obj_id: this.objId,
                    option: this.formatFieldOption(this.newFieldInfo.propertyType, this.newFieldInfo.option),
                    bk_supplier_account: this.bkSupplierAccount,
                    bk_property_id: this.newFieldInfo.propertyId,
                    bk_property_name: this.newFieldInfo.propertyName,
                    bk_property_type: this.newFieldInfo.propertyType,
                    editable: this.newFieldInfo.editable,
                    placeholder: this.newFieldInfo.placeholder,
                    unit: this.newFieldInfo.unit,
                    bk_asst_obj_id: this.newFieldInfo['bk_asst_obj_id']
                }
                this.$axios.post('object/attr', params, {id: 'saveNew'}).then(res => {
                    if (res.result) {
                        this.getModelField()
                        this.closeAddFieldBox()
                        this.$emit('newField') // 更新字段分栏列表
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
                关闭添加字段弹窗
            */
            closeAddFieldBox () {
                this.isAddFieldShow = false
            },
            /*
                新增字段时类型改变回调 主要是清空内容
                opt: 当前项数据
                index: 当前项索引
            */
            fieldTypeChange (opt, index) {
                switch (opt.value) {
                    case 'int':
                        this.newFieldInfo.option = {
                            min: '',
                            max: ''
                        }
                        break
                    case 'longchar':
                        this.newFieldInfo.option = ''
                        break
                    case 'singlechar':
                        this.newFieldInfo.option = ''
                        break
                    case 'enum':
                        this.newFieldInfo.option = {
                            list: [{id: '', name: ''}],
                            defaultIndex: 0
                        }
                        break
                    case 'singleasst':
                    case 'multiasst':
                        this.newFieldInfo.option = ''
                        this.newFieldInfo['bk_asst_obj_id'] = ''
                }
            },
            /*
                保存变更
                item: 当前项
                index: 索引
            */
            saveFieldChange (item, index) {
                $('#validate-form-change').parsley().validate()
                if (!$('#validate-form-change').parsley().isValid()) return
                if (!this.checkChangeParams(item, index)) {
                    return
                }
                let option = this.formatFieldOption(item['bk_property_type'], item.option)
                let params = {
                    description: item['description'],
                    editable: this.curFieldInfo['editable'],
                    placeholder: this.curFieldInfo['placeholder'],
                    unit: this.curFieldInfo['unit'],
                    isonly: this.curFieldInfo['isonly'],
                    isreadonly: false,
                    isrequired: this.curFieldInfo['isrequired'],
                    bk_property_group: 'default',
                    option: this.formatFieldOption(item['bk_property_type'], item.option),
                    bk_property_name: this.curFieldInfo['bk_property_name'],
                    bk_property_type: item['bk_property_type'],
                    bk_asst_obj_id: '',
                    bk_asst_forward: ''
                }
                // 只有关联类型才添加以下三个参数
                if (item['bk_property_type'] === 'singleasst' || item['bk_property_type'] === 'multiasst') {
                    if (option !== 'undefined') {
                        params['bk_asst_obj_id'] = option.value
                    }
                    params['bk_property_id'] = item['bk_property_id']
                    params['bk_obj_id'] = this.objId
                    params['bk_supplier_account'] = this.bkSupplierAccount
                }
                this.$axios.put(`object/attr/${item['id']}`, params, {id: 'saveChange'}).then(res => {
                    if (res.result) {
                        this.getModelField()
                        this.curFieldInfoCopy = this.$deepClone(this.curFieldInfo)
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            /*
                取消变更
                item: 当前项
                index: 索引
            */
            cancelFieldChange (item, index) {
                this.toggleDetailShow(item, index)
            },
            /*
                添加枚举字段
                type: new 新增字段时新增 change 查看详情时新增
            */
            addEnum (type, index, idx) {
                if (type === 'new') {
                    this.newFieldInfo.option.list.push({name: ''})
                } else {
                    let option = this.fieldList[idx].option
                    this.fieldList[idx].option.list.push({name: ''})
                }
                this.$forceUpdate()
                this.forceValidate(type)
            },
            edit (type, index, idx) {
                if (type === 'new') {
                    this.newFieldInfo.option[index].type = 'input'
                } else {
                    this.fieldList[idx].option[index].type = 'input'
                    this.fieldList.splice()
                }
            },
            /*
                向上调整枚举字段位置
                type: 新增还是修改 new 新增  change 修改
                index: 索引
                idx: type为new时才有值  值为当前字段在fieldList列表中的索引
            */
            enumUp (type, index, idx) {
                let enumList = []
                if (type === 'new') {
                    enumList = this.newFieldInfo.option.list
                    if (index === this.newFieldInfo.option.defaultIndex) {
                        this.newFieldInfo.option.defaultIndex--
                    } else if (index === this.newFieldInfo.option.defaultIndex + 1) {
                        this.newFieldInfo.option.defaultIndex++
                    }
                } else {
                    enumList = this.fieldList[idx].option.list
                    if (index === this.fieldList[idx].option.defaultIndex) {
                        this.fieldList[idx].option.defaultIndex--
                    } else if (index === this.fieldList[idx].option.defaultIndex + 1) {
                        this.fieldList[idx].option.defaultIndex++
                    }
                }
                let temp = enumList[index - 1]
                enumList.splice(index - 1, 1, enumList[index])
                enumList.splice(index, 1, temp)
                this.$forceUpdate()
                this.forceValidate(type)
            },
            /*
                向下调整枚举字段位置
                type: 新增还是修改 new 新增  change 修改
                index: 索引
                idx: type为new时才有值  值为当前字段在fieldList列表中的索引
            */
            enumDown (type, index, idx) {
                let enumList = []
                if (type === 'new') {
                    enumList = this.newFieldInfo.option.list
                    if (index === this.newFieldInfo.option.defaultIndex) {
                        this.newFieldInfo.option.defaultIndex++
                    } else if (index === this.newFieldInfo.option.defaultIndex - 1) {
                        this.newFieldInfo.option.defaultIndex--
                    }
                } else {
                    enumList = this.fieldList[idx].option
                    if (index === this.fieldList[idx].option.defaultIndex) {
                        this.fieldList[idx].option.defaultIndex++
                    } else if (index === this.fieldList[idx].option.defaultIndex - 1) {
                        this.fieldList[idx].option.defaultIndex--
                    }
                }
                let temp = enumList[index + 1]
                enumList.splice(index + 1, 1, enumList[index])
                enumList.splice(index, 1, temp)
                this.$forceUpdate()
                this.forceValidate(type)
            },
            /*
                删除枚举字段
                type: 新增还是修改 new 新增  change 修改
                index: 索引
                idx: type为new时才有值  值为当前字段在fieldList列表中的索引
            */
            deleteEnum (type, index, idx) {
                if (type === 'new') {
                    if (index === this.newFieldInfo.option.defaultIndex) {
                        this.newFieldInfo.option.defaultIndex = 0
                    }
                    this.newFieldInfo.option.list.splice(index, 1)
                } else {
                    if (index === this.fieldList[idx].option.defaultIndex) {
                        this.fieldList[idx].option.defaultIndex = 0
                    }
                    this.fieldList[idx].option.list.splice(index, 1)
                }
                this.$forceUpdate()
                this.forceValidate(type)
            },
            /*
                获取模型分类以及模型信息
            */
            getClassification () {
                this.$axios.post(`object/classification/${this.bkSupplierAccount}/objects`, {}).then(res => {
                    if (res.result) {
                        this.modelList = []
                        res.data.map(val => {
                            if (val.hasOwnProperty('bk_objects') && val['bk_objects'].length) {
                                // 不显示自己 obj[i].ObjId === this.objId
                                // 暂时不显示主机 进程 业务 集群 模块
                                for (let obj = val['bk_objects'], i = obj.length - 1; i >= 0; i--) {
                                    if (obj[i]['bk_obj_id'] === 'plat' || obj[i]['bk_obj_id'] === 'process' || obj[i]['bk_obj_id'] === 'set' || obj[i]['bk_obj_id'] === 'module' || obj[i]['bk_obj_id'] === this.objId) {
                                        obj.splice(i, 1)
                                    }
                                }
                                // 去掉临时不显示的内容后长度不为0时才添加到列表中
                                if (val['bk_objects'].length && val['bk_classification_id'] !== 'bk_biz_topo') {
                                    this.modelList.push(val)
                                }
                            }
                        })
                    } else {
                        this.$alertMsg(res['bk_error_msg'])
                    }
                })
            },
            modelChange (item) {
                this.fieldList[this.curIndex].option = item
            },
            modelSelected (item) {
                this.isSelectErrorShow = false
                this.newFieldInfo['bk_asst_obj_id'] = item.value
                // this.newFieldInfo.option = item
            },
            /*
                页面初始化
            */
            init () {
                this.getModelField()
                this.getClassification()
            },
            addNoRepeatValidator () {
                if (!window.Parsley.hasValidator('noRepeat')) {
                    window.Parsley.addValidator('noRepeat', {
                        requirementType: 'string',
                        validateString: function (inputValue, inputType) {
                            let allFieldsName = []
                            let enumFields = document.querySelectorAll('[data-parsley-no-repeat="' + inputType + '"]')
                            enumFields.forEach((enumField, index) => {
                                allFieldsName.push(enumField.value)
                            })
                            let firstIndex = allFieldsName.indexOf(inputValue)
                            let lastIndex = allFieldsName.lastIndexOf(inputValue, allFieldsName.length - 1)
                            if (inputValue && firstIndex !== -1 && firstIndex !== lastIndex) {
                                return false
                            } else {
                                return true
                            }
                        },
                        messages: {
                            'en': 'This value should not be repeated',
                            'zh-cn': this.$t('Common["重复的值"]')
                        }
                    })
                }
            },
            forceValidate (type) {
                this.$nextTick(() => {
                    this.$el.querySelectorAll('[data-parsley-no-repeat="' + type + '"]').forEach((enumField) => {
                        $(enumField).parsley().validate()
                    })
                })
            },
            forceUpdate (type) {
                this.$nextTick(() => {
                    if (type === 'change') {
                        this.$forceUpdate()
                    }
                    this.forceValidate(type)
                })
            }
        },
        directives: {
            focus: {
                inserted: function (el) {
                    el.focus()
                }
            }
        },
        mounted () {
            this.newFieldInfo.propertyGroup = 'default'
            this.addNoRepeatValidator()
            if (this.objId === '') {
                this.$refs.baseInfo.clearData()
            } else {
                this.$refs.baseInfo.getBaseInfo(this.objId)
            }
        }
    }
</script>


<style media="screen" lang="scss" scoped>
    $borderColor: #e7e9ef; //边框色
    $defaultColor: #ffffff; //默认
    $primaryColor: #f9f9f9; //主要
    $fnMainColor: #bec6de; //文案主要颜色
    $primaryHoverColor: #6b7baa; // 主要颜色
    .select-error,
    .error-msg{
        font-size: 12px;
        color: #ff3737;
    }
    .add-field-btn{
        width: 90px;
        height: 32px;
        line-height: 32px;
        border: none;
        background-color: #30d878;
        color:#fff;
    }
    .icon-btn{  //单纯图标的按钮
        background: #ffffff;
        color: $primaryHoverColor;
        cursor: pointer;
        &:hover{
            background: $primaryHoverColor;
            color: $defaultColor;
        }
    }
    .icon-cc-del{
        color: #c3cdd7;
    }
    .no-border-btn{    //无边框按钮
        background: #fff;
        color: $primaryHoverColor;
        cursor: pointer;
        &:hover{
            background: $primaryHoverColor;
            color: #fff;
            i{
                background: $primaryHoverColor;
                color: #fff;
            }
            span{
                background: $primaryHoverColor;
                color: #fff;
            }
        }
    }
    .allField{
        height: 100%;
        padding: 0 30px;
    }
    .tab-content{
        /* border-bottom: solid 1px $borderColor; */
        height: auto !important;
        &.model-field-content{
            height: calc(100% - 187px) !important;
        }
        .add-field{
            text-align: left;
            border-top: solid 1px $borderColor;
            padding-top: 20px;
            .btn-group{
                float: right;
                font-size: 0;
                .btn{
                    position: relative;
                    overflow: hidden;
                    input[type="file"]{
                        position: absolute;
                        left: -70px;
                        top: 0;
                        opacity: 0;
                        width: calc(100% + 70px);
                        height: 100%;
                        cursor: pointer;
                    }
                }
                .form{
                    display: inline-block;
                }
            }
        }
        .table-content{
            height: calc(100% - 33px) !important;
            margin-top: 10px;
            overflow-y: auto;
            width: 720px;
            /*无数据提示样式*/
            .no-data-tip{
                p{
                    font-size: 14px;
                    color: $primaryHoverColor;
                }
                .create-field-btn{
                    height: 38px;
                    line-height: 38px;
                    border-radius: 2px;
                    padding: 0 74px;
                    color: #fff;
                    font-size: 14px;
                    border: none;
                    margin-top: 15px;
                }
            }
            /*字段配置列表头部*/
            .title-content{
                width: 100%;
                height:40px;
                >ul{
                    >li{
                      float:left;
                      height:40px;
                      line-height:40px;
                      background:#f9f9f9;
                      text-align:center;
                      border:1px solid $borderColor;
                      border-right:none;
                      &:nth-child(1){
                        width:80px;
                        }
                        &:nth-child(2){
                            width:80px;
                        }
                        &:nth-child(3){
                            width:148px;
                        }
                        &:nth-child(4){
                            width:287px;
                        }
                        &:nth-child(5){
                            width:105px;
                            border-right:1px solid $borderColor;
                        }
                    }
                }
            }
            .list-content-wrapper{
                overflow-y: auto;
                height: calc(100% - 40px);
                /* margin-right: 10px; */
                &::-webkit-scrollbar{
                    width: 6px;
                    height: 5px;
                }
                &::-webkit-scrollbar-thumb{
                    border-radius: 20px;
                    background: #a5a5a5;
                }
            }
            /*列表内容样式*/
            .list-content{
                width:100%;
                cursor: pointer;
                &.editable{
                    ul{
                        li{
                            color: #ccc;
                        }
                    }
                }
                >ul{
                    width:100%;
                    height:40px;
                    padding:0;
                    margin:0;
                    >li{
                        float:left;
                        height:40px;
                        line-height:40px;
                        background:#ffffff;
                        text-align:center;
                        border-bottom:1px solid $borderColor;
                        text-overflow: ellipsis;
                        overflow: hidden;
                        &:nth-child(1){
                            width:80px;
                            border-left:1px solid $borderColor;
                        }
                        &:nth-child(2){
                            width:80px;
                        }
                        &:nth-child(3){
                            width:148px;
                        }
                        &:nth-child(4){
                            width:287px;
                        }
                        &:nth-child(5){
                            width:105px;
                            border-right:1px solid $borderColor;
                        }
                    }
                }
                .list-content-hidden{
                    width:100%;
                    .enum-table{
                        position: relative;
                        &.disabled{
                            cursor: not-allowed;
                        }
                        .enum-disabled{
                            position: absolute;
                            left: 0;
                            right: 0;
                            top: 0;
                            bottom: 0;
                        }
                    }
                    .form-common{
                        width: 700px;
                        /*background: #f9f9f9;*/
                        padding: 30px 19px 30px 17px;
                        border: 1px solid #e7e9ef;
                        border-top: 0;
                        &.dn{
                            display:none;
                        }
                        h3{
                            margin:0;
                            margin-bottom:10px;
                        }
                        .form-common-item{
                            width:213px;
                            margin-right:0;
                            &.form-common-item2{
                                width: 66.7%;
                            }
                            &.block{
                                width: 100%;
                                .form-common-content{
                                    width: calc(100% - 92px);
                                }
                            }
                            .form-common-label{
                                display: inline-block;
                                width:75px;
                                vertical-align: top;
                                line-height: 30px;
                                text-align:right;
                            }
                            .form-common-content{
                                // margin-left:5px;
                                width:128px;
                                input{
                                    width:100%;
                                }
                            }
                            .selcet-width-control{
                                text-align:left;
                            }
                        }
                        .correlate-more-control{
                            width:290px!important;
                            label{
                                width: auto !important;
                            }
                        }
                    }
                }
                .btn-contain{
                    .editable{
                        cursor: not-allowed;
                    }
                    i{
                        cursor:pointer;
                    }
                }
            }
        }
    }
    /*新增字段*/
    .add-field-wrapper{
        position: absolute;
        bottom: 0;
        left: 0;
        right: 0;
        top: 127px;
        height: calc(100% - 139px);
        .add-field-detail{
            padding-left:20px;
            position:absolute;
            bottom:60px;
            left:0;
            right:0;
            z-index: 400;
            top: 0;
            background: #fff;
            padding-bottom: 50px;
            padding:0 20px 20px 20px;
            bottom: 0;
            .bg-titel{
                text-align: center;
                width: 100%;
                height: 40px;
                background:url('../../../common/images/bg_title.png') no-repeat;
                display: block;
                cursor: pointer;
                >img{
                    cursor: pointer;
                    margin-top: 14px;
                }
            }
            .content-hidden{
                display:none;
            }
            .border-control{
                padding:0 40px 0 40px ;
                min-height: 400px;
                height: calc(100% - 30px);
            }
            .content-replace{
                display:none;
            }
            .title{
                h3{
                    border-left: 4px solid $primaryHoverColor;
                    padding-left:4px;
                    font-size:14px;
                    line-height:1;
                }
            }
            .form-common{
                width: 661px;
                margin-top: 20px;
                height: calc(100% - 20px);
                .form-option {
                    margin-top: 10px;
                    overflow: auto;
                    max-height: calc(100% - 254px);
                    @include scrollbar;
                }
                .form-common-item{
                    .form-common-label{
                        width: 75px;
                        display: inline-block;
                        text-align: right;
                        vertical-align: top;
                        line-height: 30px;
                        color: $primaryHoverColor;
                    }
                }
            }
            .button-wraper{
                margin-left: 82px;
                .bk-button{
                    height: 30px;
                    line-height: 28px;
                    border-radius: 2px;
                    text-align: center;
                    display:inline-block;
                    margin-top:30px;
                    cursor:pointer;
                    font-size: 12px;
                    padding: 0 20px;
                    min-width: 90px;
                }
            }
        }
    }
    .form-common{
        color: $primaryHoverColor;
        h3{
            font-size:14px;
            text-align:left;
            border-left:4px solid $primaryHoverColor;
            line-height:1;
            padding-left:5px;
            font-weight: normal;
            margin: 0;
            margin-bottom: 10px;
        }
        .bk-form-radio{
            margin-right: 10px;
        }
        .form-common-item{
            width: 33.3%;
            float:left;
            .icon-info-png{
                position: absolute;
                display: inline-block;
                top: 7px;
                right: -20px;
                width: 16px;
                height: 16px;
                background: url(../../../common/images/icon/icon-info.png);
            }
            &.form-common-item2{
                width: 66.7%;
            }
            &.block{
                width: 100%;
                .form-common-content{
                    width: calc(100% - 83px);
                    input{
                        width: 100% !important;
                    }
                }
            }
            &.disabled{
                input{
                    background:#f9f9f9;
                }
            }
            .from-selcet-wrapper{
                display:inline-block;
                .bk-form-checkbox{
                    margin-right:0;
                }
                label{
                    color: $primaryHoverColor;
                    font-style: normal;
                }
            }
            .form-common-label{
                span{
                    color:#f05d5d;
                }
            }
            .form-common-content{
                display:inline-block;
                margin-left:2px;
                width: 130px;
                &.reg-verification{
                    input{
                        width:130px;
                    }
                }
                &.interior-width-control{
                    input{
                        width: 130px;
                    }
                }
                &.selcet-width-control{
                   width:130px;
                }
                .select-content{
                   width: 130px;
                }
                input{
                   width: 130px;
                   height: 30px;
                   line-height: 30px;
                   border: 1px solid #bec6de;
                   padding: 0 10px;
                   outline: none;
                   border-radius: 2px;
                }
                input[type=number] {
                    padding: 0 0 0 10px;
                }
            }
        }
        .submit-btn{
            margin-left: 82px;
            .bk-button{
                height: 30px;
                line-height: 28px;
                border-radius: 2px;
                display:inline-block;
                margin-top:30px;
                cursor:pointer;
                font-size: 12px;
                text-align: center;
                padding: 0 20px;
                min-width: 90px;
            }
        }
        .from-table-content{
            width:100%;
            margin-top:60px;
            .add-enum{
                display:inline-block;
                cursor:pointer;
            }
            >thead{
                tr{
                    th{
                        text-align:center;
                        background:#e7e9ef;
                        padding:5px 20px;
                    }
                }
            }
            >tbody{
                tr{
                    td{
                        text-align:center;
                        border: none;
                        border-bottom:1px solid #e7e9ef;
                        padding:5px 20px;
                        color:#bec6de;
                        input{
                            outline: none;
                            border: 1px solid #bec6de;
                            padding: 2px 10px;
                        }
                        i{
                           cursor:pointer;
                        }
                    }
                }
            }
        }
    }
    .list-wrapper{
        margin: 20px 0 0 69px;
        border-top: 1px solid #dde4eb;
        border-left: 1px solid #dde4eb;
        color: #737987;
        .list-item{
            position: relative;
            font-size: 0;
            line-height: 40px;
            display: flex;
            &:first-child{
                background: #fafbfd;
                color: #333948;
            }
            .icon{
                position: absolute;
                font-size: 16px;
                left: -25px;
                top: 13px;
                color: #ff5656;
            }
            div{
                text-align: center;
                font-size: 14px;
                border-right: 1px solid #dde4eb;
                border-bottom: 1px solid #dde4eb;
                text-overflow: ellipsis;
                overflow: hidden;
                &:not(:last-child) {
                    flex: 2;
                }
                &:last-child{
                    flex: 1;
                }
                .add{
                    cursor: pointer;
                    &:hover{
                        color: #3c96ff;
                    }
                }
                .delete{
                    margin-left: 8px;
                    .icon-cc-del{
                        color: #737987;
                        &:hover{
                            color: #ff3737;
                        }
                    }
                    cursor: pointer;
                }
                .list-view{
                    cursor: pointer;
                    display: inline-block;
                    padding: 0 10px;
                    width: 100%;
                    height: 40px;
                    vertical-align: bottom;
                }
            }
        }
    }
</style>

<style>
    .bk-tooltips{
        z-index: 9999999
    }
</style>
<style lang="scss" scoped>
    .form-enum-box {
        width: 100%;
        height: 100%;
        overflow: auto;
    }
    .form-enum-wrapper{
        margin-top: 10px;
        margin-left: 82px;
        font-size: 0;
        position: relative;
        float: left;
        // width: calc(100% - 82px);
        .enum-id{
            float: left;
            width: 90px;
            margin-right: 10px;
            input{
                width: 90px;
            }
        }
        .enum-name{
            float: left;
            width: 250px;
            input{
                width: 250px;
            }
        }
        input {
            font-size: 12px;
            vertical-align: middle;
            // width: 350px;
            height: 30px;
            border-radius: 2px;
            border: 1px solid #bec6de;
            padding: 0 10px;
        }
        .enum-radio{
            display: inline-block;
            width: 30px;
            position: absolute;
            left: -30px;
            top: 0;
            cursor: pointer;
        }
        .span-enum-radio{
            display: inline-block;
            position: absolute;
            left: -30px;
            line-height: 30px;
            top: 0;
            cursor: pointer;
            margin-top: 7px;
            width: 16px;
            height: 16px;
            border: 1px solid #aeaeae;
            border-radius: 50%;
            &.active{
                border: 5px solid #3c96ff;
            }
        }
        button.bk-icon{
            display: inline-block;
            width: 30px;
            height: 30px;
            margin-left: 5px;
            vertical-align: middle;
            text-align: center;
            font-size: 14px;
            line-height: 1;
            border: 1px solid #bec6de;
            background-color: #fff;
            outline: 0;
            &:disabled{
                cursor: not-allowed;
                background-color: #eee;
                border-color: #eee;
            }
        }
        .form-enum-wrapper-dot,
        .form-enum-wrapper-dot:before,
        .form-enum-wrapper-dot:after{
            position: absolute;
            left: 10px;
            width: 3px;
            height: 3px;
            background-color: #bec6de;
        }
        .form-enum-wrapper-dot{
            top: 14px;
            &:before{
                content: '';
                left: 0;
                bottom: 5px;
            }
            &:after{
                content: '';
                left: 0;
                top: 5px;
            }
        }
        input:focus{
            & ~ .form-enum-wrapper-dot,
            & ~ .form-enum-wrapper-dot:before,
            & ~ .form-enum-wrapper-dot:after{
                background-color: #498fe0;
            }
        }
    }
</style>
