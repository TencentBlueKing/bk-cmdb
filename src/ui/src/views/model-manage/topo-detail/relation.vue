<template>
    <div>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["源模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstId')}">
                <input type="text" class="cmdb-form-input"
                name="asstId"
                v-validate="'required|singlechar'"
                :placeholder="$t('ModelManagement[\'请输入英文标识\']')">
                <p class="form-error">{{errors.first('asstId')}}</p>
            </div>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["目标模型"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstId')}">
                <input type="text" class="cmdb-form-input"
                name="asstId"
                v-validate="'required|singlechar'"
                :placeholder="$t('ModelManagement[\'请输入英文标识\']')">
                <p class="form-error">{{errors.first('asstId')}}</p>
            </div>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关系类型"]')}}
                <span class="color-danger">*</span>
            </span>
            <ul class="relation-label cmdb-form-item clearfix" :class="{'is-error': errors.has('asstId')}">
                <li class="active">asdf</li>
                <li>啊但是发生地方</li>
                <li>啊但是发生地方</li>
                <li>啊但是发生地方</li>
                <li>啊但是发生地方</li>
                <li>啊但是发生地方</li>
                <li>啊但是发生地方</li>
                <li>自定义关系</li>
            </ul>
        </label>
        <label class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关系描述"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('asstId')}">
                <input type="text" class="cmdb-form-input"
                name="asstId"
                v-validate="'required|singlechar'"
                :placeholder="$t('ModelManagement[\'请输入关系描述\']')">
                <p class="form-error">{{errors.first('asstId')}}</p>
            </div>
        </label>
        <div class="form-label">
            <span class="label-text">
                {{$t('ModelManagement["关系约束"]')}}
                <span class="color-danger">*</span>
            </span>
            <div class="cmdb-form-item" :class="{'is-error': errors.has('mapping')}">
                <cmdb-selector
                    :list="mappingList"
                    v-validate="'required'"
                    name="mapping"
                    v-model="relationInfo.mapping"
                ></cmdb-selector>
                <p class="form-error">{{errors.first('mapping')}}</p>
            </div>
            <i class="bk-icon icon-info-circle"></i>
        </div>
    </div>
</template>

<script>
    export default {
        data () {
            return {
                mappingList: [{
                    id: '1:1',
                    name: '1-1'
                }, {
                    id: '1:n',
                    name: '1-N'
                }, {
                    id: 'n:n',
                    name: 'N-N'
                }],
                relationInfo: {
                    mapping: ''
                }
            }
        }
    }
</script>

<style lang="scss" scoped>
    .relation-label {
        li {
            float: left;
            margin: 5px 7px 5px 0;
            padding: 0 8px;
            height: 26px;
            line-height: 24px;
            font-size: 12px;
            border: 1px solid $cmdbTableBorderColor;
            background: #f5f7f9;
            cursor: pointer;
            &:hover {
                background: #fafafa;
            }
            &.active {
                color: #fff;
                background: $cmdbBorderFocusColor;
                border-color: $cmdbBorderFocusColor;
            }
            &:last-child {
                margin-right: 0;
                padding: 0;
                border: none;
                color: $cmdbBorderFocusColor;
                background: #fff;
            }
        }
    }
</style>
