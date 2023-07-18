<!--
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
-->

<template>
  <div class="node-create-layout">
    <h2 class="node-create-title">{{ $t('新建模块') }}</h2>
    <div class="node-create-path" :title="topoPath">
      {{ $t('添加节点已选择') }}：{{ topoPath }}
    </div>
    <div
      class="node-create-form"
      :style="{
        'max-height': Math.min($APP.height - 400, 400) + 'px',
      }">
      <div class="form-item clearfix mt30">
        <div class="create-type fl">
          <input
            id="formTemplate"
            v-model="withTemplate"
            class="type-radio"
            type="radio"
            name="createType"
            :value="1" />
          <label for="formTemplate">{{ $t('从模板新建') }}</label>
        </div>
        <div class="create-type fl ml50">
          <input
            id="createDirectly"
            v-model="withTemplate"
            class="type-radio"
            type="radio"
            name="createType"
            :value="0" />
          <label for="createDirectly">{{ $t('直接新建') }}</label>
        </div>
      </div>
      <div v-if="withTemplate" class="form-item">
        <label>{{ $t('服务模板') }}</label>
        <bk-select
          key="template"
          v-model="template"
          v-validate="'required'"
          style="width: 100%"
          :clearable="false"
          :searchable="templateList.length > 7"
          :loading="$loading(request.serviceTemplate)"
          data-vv-name="template">
          <bk-option
            v-for="(option, index) in templateList"
            :id="option.id"
            :key="index"
            :name="option.name">
          </bk-option>
          <div
            v-if="!templateList.length"
            slot="extension"
            class="add-template"
            @click="jumpServiceTemplate">
            <i class="bk-icon icon-plus-circle"></i>
            <span>{{ $t('新建服务模板') }}</span>
          </div>
        </bk-select>
        <span v-if="errors.has('template')" class="form-error">{{
          errors.first('template')
        }}</span>
      </div>
      <div class="form-item">
        <label>
          {{ $t('模块名称') }}
          <font color="red">*</font>
          <i
            v-if="withTemplate === 1"
            v-bk-tooltips.top="$t('模块名称提示')"
            class="icon-cc-tips">
          </i>
        </label>
        <cmdb-form-singlechar
          key="moduleName"
          v-model="moduleName"
          v-validate="'required|businessTopoInstNames|length:256'"
          data-vv-name="moduleName"
          data-vv-validate-on="blur"
          :placeholder="$t('请输入xx', { name: $t('模块名称') })"
          :disabled="!!withTemplate">
        </cmdb-form-singlechar>
        <span v-if="errors.has('moduleName')" class="form-error">{{
          errors.first('moduleName')
        }}</span>
      </div>
      <div v-if="!withTemplate" class="form-item clearfix">
        <label>{{ $t('所属服务分类') }}<font color="red">*</font></label>
        <cmdb-selector
          key="firstClass"
          v-model="firstClass"
          v-validate="'required'"
          class="service-class fl"
          data-vv-name="firstClass"
          :auto-select="false"
          :list="firstClassList"
          :loading="$loading(request.serviceCategory)"
          @on-selected="updateCategory">
        </cmdb-selector>
        <cmdb-selector
          key="secondClass"
          v-model="secondClass"
          v-validate="'required'"
          class="service-class fr"
          data-vv-name="secondClass"
          :list="secondClassList"
          :loading="$loading(request.serviceCategory)">
        </cmdb-selector>
        <span v-if="errors.has('firstClass')" class="form-error">{{
          errors.first('firstClass')
        }}</span>
        <span
          v-if="errors.has('secondClass')"
          class="form-error second-class"
          >{{ errors.first('secondClass') }}</span
        >
      </div>
    </div>
    <div class="node-create-options">
      <bk-button
        v-test-id="'createModuleSave'"
        theme="primary"
        :disabled="$loading() || errors.any()"
        @click="handleSave">
        {{ $t('提交') }}
      </bk-button>
      <bk-button theme="default" @click="handleCancel">{{
        $t('取消')
      }}</bk-button>
    </div>
  </div>
</template>

<script>
import has from 'has'

import { MENU_BUSINESS_SERVICE_TEMPLATE } from '@/dictionary/menu-symbol'
import serviceTemplateService from '@/service/service-template/index.js'

export default {
  props: {
    parentNode: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      withTemplate: 1,
      createTypeList: [
        {
          id: 1,
          name: this.$t('从模板创建'),
        },
        {
          id: 0,
          name: this.$t('直接创建'),
        },
      ],
      template: '',
      templateList: [],
      moduleName: '',
      firstClass: '',
      firstClassList: [],
      secondClass: '',
      values: {},
      request: {
        serviceTemplate: Symbol('serviceTemplate'),
        serviceCategory: Symbol('serviceCategory'),
      },
    }
  },
  computed: {
    topoPath() {
      const nodePath = [...this.parentNode.parents, this.parentNode]
      return nodePath.map(node => node.data.bk_inst_name).join('/')
    },
    business() {
      return this.$store.getters['objectBiz/bizId']
    },
    serviceTemplateMap() {
      return this.$store.state.businessHost.serviceTemplateMap
    },
    currentTemplate() {
      return this.templateList.find(item => item.id === this.template) || {}
    },
    categoryMap() {
      return this.$store.state.businessHost.categoryMap
    },
    currentCategory() {
      return (
        this.firstClassList.find(category => category.id === this.firstClass) ||
        {}
      )
    },
    secondClassList() {
      return this.currentCategory.secondCategory || []
    },
  },
  watch: {
    withTemplate(withTemplate) {
      if (withTemplate) {
        this.updateCategory()
        this.template = this.templateList.length ? this.templateList[0].id : ''
      } else {
        this.template = ''
        this.updateCategory(1)
        this.getServiceCategories()
      }
    },
    template(template) {
      if (template) {
        this.moduleName = this.currentTemplate.name
      } else {
        this.moduleName = ''
      }
    },
  },
  created() {
    this.getServiceTemplates()
  },
  methods: {
    async getServiceTemplates() {
      if (has(this.serviceTemplateMap, this.business)) {
        this.templateList = this.serviceTemplateMap[this.business]
      } else {
        try {
          const templates = await serviceTemplateService.findAll(
            { bk_biz_id: this.business, page: { sort: '-last_time' } },
            {
              requestId: this.request.serviceTemplate,
            }
          )

          this.templateList = templates
          this.$store.commit('businessHost/setServiceTemplate', {
            id: this.business,
            templates,
          })
        } catch (e) {
          console.error(e)
          this.templateList = []
        }
      }
      this.template = this.templateList[0] ? this.templateList[0].id : ''
    },
    async getServiceCategories() {
      if (has(this.categoryMap, this.business)) {
        this.firstClassList = this.categoryMap[this.business]
      } else {
        try {
          const data = await this.$store.dispatch(
            'serviceClassification/searchServiceCategory',
            {
              params: { bk_biz_id: this.business },
              config: {
                requestId: this.request.serviceCategory,
              },
            }
          )
          const categories = this.collectServiceCategories(data.info)
          this.firstClassList = categories
          this.$store.commit('businessHost/setCategories', {
            id: this.business,
            categories,
          })
        } catch (e) {
          console.error(e)
          this.firstClassList = []
        }
      }
    },
    collectServiceCategories(data) {
      const categories = []
      data.forEach(item => {
        if (!item.category.bk_parent_id) {
          categories.push(item.category)
        }
      })
      categories.forEach(category => {
        // eslint-disable-next-line max-len
        category.secondCategory = data
          .filter(item => item.category.bk_parent_id === category.id)
          .map(item => item.category)
      })
      return categories
    },
    updateCategory(firstClass) {
      if (firstClass) {
        this.firstClass = firstClass
        this.secondClass = this.secondClassList.length
          ? this.secondClassList[0].id
          : ''
      } else {
        this.firstClass = ''
        this.secondClass = ''
      }
    },
    handleSave() {
      this.$validator.validateAll().then(isValid => {
        if (isValid) {
          this.$emit('submit', {
            bk_module_name: this.moduleName,
            service_category_id: this.withTemplate
              ? this.currentTemplate.service_category_id
              : this.secondClass,
            service_template_id: this.withTemplate ? this.template : 0,
          })
        }
      })
    },
    handleCancel() {
      this.$emit('cancel')
    },
    jumpServiceTemplate() {
      this.$routerActions.redirect({
        name: MENU_BUSINESS_SERVICE_TEMPLATE,
        params: {
          bizId: this.business,
        },
      })
    },
  },
}
</script>

<style lang="scss" scoped>
.node-create-layout {
  position: relative;
}

.node-create-title {
  margin-top: -14px;
  padding: 0 26px;
  line-height: 30px;
  font-size: 22px;
  color: #333948;
}

.node-create-path {
  padding: 23px 26px 0;
  margin: 0 0 -5px;
  font-size: 12px;

  @include ellipsis;
}

.node-create-form {
  padding: 0 26px 27px;
  overflow: visible;
}

.form-item {
  margin: 15px 0 0;
  position: relative;

  label {
    display: block;
    padding: 7px 0;
    line-height: 19px;
    font-size: 14px;
    color: #63656e;
  }

  .service-class {
    width: 260px;

    @include inlineBlock;
  }

  .form-error {
    position: absolute;
    top: 100%;
    left: 0;
    font-size: 12px;
    color: $cmdbDangerColor;

    &.second-class {
      left: 270px;
    }
  }

  .create-type {
    display: flex;
    align-items: center;

    .type-radio {
      appearance: none;
      width: 16px;
      height: 16px;
      padding: 3px;
      border: 1px solid #979ba5;
      border-radius: 50%;
      background-clip: content-box;
      outline: none;
      cursor: pointer;

      &:checked {
        border-color: #3a84ff;
        background-color: #3a84ff;
      }
    }

    label {
      padding: 0 0 0 6px;
      font-size: 14px;
      cursor: pointer;
    }
  }
}

.node-create-options {
  padding: 9px 20px;
  border-top: 1px solid $cmdbBorderColor;
  text-align: right;
  background-color: #fafbfd;
}

font {
  padding: 0 2px;
}

.add-template {
  width: 20%;
  cursor: pointer;

  .icon-plus-circle {
    @include inlineBlock;

    font-size: 14px;
  }

  span {
    @include inlineBlock;
  }
}
</style>
