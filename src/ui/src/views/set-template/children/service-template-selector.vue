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
  <section>
    <div class="top" v-if="allTemplates.length">
      <bk-select
        class="mr10"
        style="width: 210px"
        :popover-width="260"
        :placeholder="$t('所有一级分类')"
        :allow-clear="true"
        :searchable="true"
        v-model="filter.primaryCategory">
        <bk-option v-for="category in primaryCategoryList"
          :key="category.id"
          :id="category.id"
          :name="category.name">
        </bk-option>
      </bk-select>
      <bk-select
        class="mr10"
        style="width: 210px"
        :popover-width="260"
        :placeholder="$t('所有二级分类')"
        :allow-clear="true"
        :searchable="true"
        :empty-text="secCategoryEmptyText"
        v-model="filter.secCategory">
        <bk-option v-for="category in secCategoryList"
          :key="category.id"
          :id="category.id"
          :name="category.name">
        </bk-option>
      </bk-select>
      <bk-input
        class="search"
        type="text"
        :placeholder="$t('请输入模板名称搜索')"
        clearable
        right-icon="bk-icon icon-search"
        v-model.trim="filter.templateName">
      </bk-input>
      <span class="select-all">
        <bk-checkbox :value="isSelectAll" :indeterminate="isHalfSelected" @change="handleSelectAll">全选</bk-checkbox>
      </span>
    </div>
    <ul class="template-list clearfix"
      v-bkloading="{ isLoading: $loading('getServiceTemplate') }"
      :style="{ height: !!templates.length ? '264px' : '306px' }"
      :class="{ 'is-loading': $loading('getServiceTemplate') }">
      <template v-if="templates.length">
        <li v-for="(template, index) in templates"
          class="template-item fl clearfix"
          :class="{
            'is-selected': localSelected.includes(template.id),
            'is-middle': index % 3 === 1,
            'disabled': $parent.$parent.serviceExistHost(template.id)
          }"
          :key="template.id"
          @click="handleClick(template, $parent.$parent.serviceExistHost(template.id))"
          @mouseenter="handleShowDetails(template, $event, $parent.$parent.serviceExistHost(template.id))"
          @mouseleave="handleHideTips">
          <i class="select-icon bk-icon icon-check-circle-shape fr"></i>
          <span class="template-name">{{template.name}}</span>
        </li>
      </template>
      <li class="template-empty" v-else>
        <div class="empty-content">
          <img class="empty-image" src="../../../assets/images/empty-content.png">
          <i18n class="empty-tips" path="无服务模板提示">
            <template #link>
              <a class="empty-link" href="javascript:void(0)" @click="handleLinkClick">{{$t('去添加服务模板')}}</a>
            </template>
          </i18n>
        </div>
      </li>
    </ul>
    <div ref="templateDetails"
      class="template-details"
      v-bkloading="{
        isLoading: $loading(processRequestId), mode: 'spin', theme: 'primary', size: 'mini'
      }"
      v-show="tips.show">
      <template v-if="!$loading(processRequestId)">
        <div class="disabled-tips" v-show="processInfo.disabled">{{$t('该模块下有主机不可取消')}}</div>
        <div class="info-item">
          <span class="label">{{$t('模板名称')}} ：</span>
          <div class="details">{{curTemplate.name}}</div>
        </div>
        <div class="info-item">
          <span class="label">{{$t('服务分类')}} ：</span>
          <div class="details">{{curTemplate.category}}</div>
        </div>
        <div class="info-item">
          <span class="label">{{$t('服务进程')}} ：</span>
          <div class="details">
            <p v-for="(item, index) in processInfo.processes" :key="index">{{item}}</p>
            <template v-if="!processInfo.processes.length">
              <p>{{$t('模板没配置进程')}}</p>
            </template>
          </div>
        </div>
      </template>
    </div>
  </section>
</template>

<script>
  import { MENU_BUSINESS_SERVICE_TEMPLATE } from '@/dictionary/menu-symbol'
  import { mapGetters } from 'vuex'
  import debounce from 'lodash.debounce'
  import serviceTemplateService from '@/service/service-template/index.js'

  export default {
    name: 'serviceTemplateSelector',
    props: {
      selected: {
        type: Array,
        default: () => []
      },
      servicesHost: {
        type: Array,
        default: () => []
      }
    },
    data() {
      return {
        allTemplates: [],
        templates: [],
        localSelected: [...this.selected],
        templateDetailsData: {},
        processRequestId: Symbol('processDetails'),
        tips: {
          show: false,
          instance: null
        },
        curTemplate: {},
        categoryList: [],
        primaryCategoryList: [],
        secCategoryList: [],
        processInfo: {
          disabled: false,
          processes: []
        },
        filter: {
          primaryCategory: '',
          secCategory: '',
          templateName: ''
        }
      }
    },
    computed: {
      ...mapGetters('objectBiz', ['bizId']),
      isSelectAll() {
        return this.localSelected.length === this.templates.length
      },
      isHalfSelected() {
        return !this.isSelectAll && this.localSelected.length > 0
      },
      isEditMode() {
        return this.$parent.$parent.mode === 'edit'
      },
      secCategoryEmptyText() {
        return this.filter.primaryCategory ? this.$t('没有二级分类') : this.$t('请选择一级分类')
      }
    },
    watch: {
      localSelected: {
        handler(value) {
          this.$emit('select-change', value)
        },
        immediate: true
      },
      allTemplates(value) {
        this.$emit('template-loaded', value)
      },
      'filter.primaryCategory'(id) {
        // 当前分类下的二级分类
        this.secCategoryList = this.categoryList.filter(item => item.bk_parent_id === id)

        if (!id) {
          this.filter.secCategory = ''
        }

        this.filterTemplate()
      },
      'filter.secCategory'() {
        this.filterTemplate()
      },
      'filter.templateName'() {
        this.filterTemplate()
      }
    },
    async created() {
      await this.getServiceCategory()
      this.getTemplates()
      this.handleShowDetails = debounce(this.showDetails, 300)
      this.handleHideTips = debounce(this.hideDetailsTips, 300)
    },
    methods: {
      async getServiceCategory() {
        try {
          const data = await this.$store.dispatch('serviceClassification/searchServiceCategoryWithoutAmout', {
            params: { bk_biz_id: this.bizId },
            config: {
              requestId: 'getServiceCategoryWithoutAmount',
              fromCache: true
            }
          })
          this.categoryList = data?.info || []
          this.primaryCategoryList = this.categoryList.filter(category => !category.bk_parent_id)
        } catch (e) {
          console.error(e)
          this.categoryList = []
        }
      },
      async getTemplates() {
        try {
          const templates = await serviceTemplateService.findAll({ bk_biz_id: this.bizId, page: { sort: 'name' } }, {
            requestId: 'getServiceTemplate'
          })

          // 将分类数据写入模板中
          templates.forEach((template) => {
            const secCategory = this.categoryList.find(item => item.id === template.service_category_id)
            const primaryCategory = this.categoryList.find(item => item.id === secCategory.bk_parent_id)
            template.parent_service_category_id = primaryCategory?.id
            template.category = `${primaryCategory?.name || '--'} / ${secCategory?.name || '--'}`
          })

          // 选中的显示在前
          this.templates = templates.sort((a, b) => this.selected.includes(b.id) - this.selected.includes(a.id))

          // 备份全量模板，用于搜索
          this.allTemplates = this.templates
        } catch (e) {
          console.error(e)
          this.templates = []
        }
      },
      handleClick(template, disabled) {
        if (disabled) return
        const index = this.localSelected.indexOf(template.id)
        if (index > -1) {
          this.localSelected.splice(index, 1)
        } else {
          this.localSelected.push(template.id)
        }
      },
      getSelectedServices() {
        return this.localSelected.map(id => this.allTemplates.find(template => template.id === id))
      },
      filterTemplate() {
        const { primaryCategory, secCategory, templateName } = this.filter
        let results = []

        if (!primaryCategory && !secCategory && !templateName) {
          results = this.allTemplates.slice()
        }

        if (primaryCategory) {
          results = this.allTemplates.filter(template => template.parent_service_category_id === primaryCategory)
        }
        if (secCategory) {
          const data = primaryCategory ? results : this.allTemplates
          results = data.filter(template => template.service_category_id === secCategory)
        }
        if (templateName) {
          const data = (primaryCategory || secCategory) ? results : this.allTemplates
          results = data.filter(template => template.name.indexOf(templateName) > -1)
        }

        this.templates = results
      },
      async showDetails(template = {}, event, disabled) {
        this.curTemplate = template
        this.processInfo.disabled = disabled

        this.tips.instance && this.tips.instance.destroy()
        this.tips.instance = this.$bkPopover(event.target, {
          content: this.$refs.templateDetails,
          zIndex: 9999,
          width: 'auto',
          trigger: 'manual',
          boundary: 'window',
          arrow: true
        })
        this.tips.show = true
        this.$nextTick(() => {
          this.tips.instance && this.tips.instance.show()
        })

        const curInfo = this.templateDetailsData[template.id]
        if (curInfo) {
          this.setProcessInfo(curInfo)
          return
        }

        try {
          const data = await this.$store.dispatch('processTemplate/getBatchProcessTemplate', {
            params: {
              bk_biz_id: this.bizId,
              service_template_id: template.id
            },
            config: {
              requestId: this.processRequestId,
              cancelPrevious: true
            }
          })
          this.setProcessInfo(data.info)
          this.templateDetailsData[template.id] = data.info
        } catch (e) {
          console.error(e)
        }
      },
      setProcessInfo(data = []) {
        this.processInfo.processes = data.map((process) => {
          const port = this.$tools.getValue(process, 'property.port.value') || ''
          return `${process.bk_process_name}${port ? `:${port}` : ''}`
        })
      },
      hideDetailsTips() {
        this.tips.instance && this.tips.instance.destroy()
        this.tips.instance = null
      },
      handleLinkClick() {
        this.$routerActions.redirect({
          name: MENU_BUSINESS_SERVICE_TEMPLATE
        })
      },
      handleSelectAll(checked) {
        const { serviceExistHost } = this.$parent.$parent
        if (checked) {
          this.localSelected = this.templates.map(template => template.id)
        } else {
          if (this.isEditMode) {
            const selectedTemplate = this.templates.filter(template => serviceExistHost(template.id))
            this.localSelected = selectedTemplate.map(template => template.id)
          } else {
            this.localSelected = []
          }
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .top {
        display: flex;
        margin-bottom: 24px;
        .search {
            width: 210px;
        }
        .select-all {
            line-height: 32px;
            margin-left: auto;
        }
    }
    .template-list {
        height: 264px;
        @include scrollbar-y;
        &.is-loading {
            min-height: 144px;
        }
        .template-item {
            width: calc((100% - 30px) / 3);
            height: 32px;
            margin: 0 0 16px 0;
            padding: 0 6px 0 10px;
            line-height: 30px;
            border-radius: 2px;
            border: 1px solid #DCDEE5;
            color: #63656E;
            cursor: pointer;
            &.is-middle {
                margin: 0 10px 16px;
            }
            &.is-selected {
                background-color: #E1ECFF;
                .select-icon {
                    font-size: 18px;
                    border: none;
                    border-radius: initial;
                    background-color: initial;
                    color: #3A84FF;
                }
            }
            &.disabled {
                cursor: not-allowed;
                .select-icon {
                    color: #C4C6CC;
                    cursor: not-allowed;
                }
            }
            .select-icon {
                width: 18px;
                height: 18px;
                font-size: 0px;
                margin: 6px 0;
                color: #fff;
                background-color: #fff;
                border-radius: 50%;
                border: 1px solid #979BA5;
            }
            .template-name {
                display: block;
                max-width: calc(100% - 18px);
                @include ellipsis;
            }
            .bk-tooltip {
                width: 100%;
                /deep/ .bk-tooltip-ref {
                    width: 100%;
                }
            }
        }
    }
    .template-empty {
        height: 280px;
        &:before {
            content: "";
            height: 100%;
            width: 0;
            @include inlineBlock;
        }
        .empty-content {
            width: 100%;
            @include inlineBlock;
            .empty-image {
                display: block;
                margin: 0 auto;
            }
            .empty-tips {
                display: block;
                text-align: center;
            }
            .empty-link {
                color: #3A84FF;
            }
        }
    }
    .template-details {
        line-height: 20px;
        min-width: 160px;
        min-height: 60px;
        .disabled-tips {
            border-bottom: 1px solid #FFFFFF;
            padding: 0 0 6px;
            margin-bottom: 6px;
            font-size: 12px;
        }
        .info-item {
            display: flex;

            .label {
                font-size: 12px;
                font-weight: 700;
            }
        }
        .details {
            font-size: 12px;
        }

        ::v-deep .bk-loading {
            // 在全局被修改了，这里再特殊修改
            background: rgba(0, 0, 0, 0) !important;
        }
    }
</style>
