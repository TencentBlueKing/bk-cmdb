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
  <bk-sideslider
    :transfer="true"
    :width="760"
    :title="$t('操作详情')"
    :is-show.sync="isShow"
    @hidden="handleHidden">
    <div class="details-content" slot="content"
      v-bkloading="{ isLoading: pending }">
      <component
        v-if="details"
        :is="detailsType"
        :details="details">
      </component>
    </div>
  </bk-sideslider>
</template>

<script>
  import DetailsJson from './details-json'
  import DetailsTable from './details-table'
  import {
    BUILTIN_MODELS,
    BUILTIN_MODEL_RESOURCE_TYPES
  } from '@/dictionary/model-constants.js'

  export default {
    components: {
      [DetailsJson.name]: DetailsJson,
      [DetailsTable.name]: DetailsTable
    },
    props: {
      id: Number,
      bizId: {
        type: Number
      },
      objId: {
        type: String,
      },
      resourceType: {
        type: String,
        default: ''
      },
      /**
       * 审计目标，不同类型的对象的权限可能不一样，所以需要区分审计目标，可选值 instance（资源实例）、common（通用），如果不需要特殊鉴权，默认为 common
       */
      aduitTarget: {
        type: String,
        default: 'common',
        validator: val => ['instance', 'common'].includes(val)
      }
    },
    data() {
      return {
        details: null,
        isShow: false,
        pending: true
      }
    },
    computed: {
      detailsType() {
        if (!this.details) {
          return null
        }

        // 用表格来展示审计详情的审计目标，除开这里的其他都用 JSON 来展示。
        const tableViewTargets = [
          BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.HOST],
          BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS],
          BUILTIN_MODEL_RESOURCE_TYPES[BUILTIN_MODELS.BUSINESS_SET],
          'module',
          'set',
          'mainline_instance',
          'model_instance',
          'cloud_area'
        ]

        let isTableViewTarget = tableViewTargets.includes(this.details.resource_type)

        // 如果拥有模型实例，则需要判断模型是否还存在
        if (this.details?.operation_detail?.bk_obj_id) {
          isTableViewTarget = this.isModelExisted()
        }

        return isTableViewTarget ? DetailsTable.name : DetailsJson.name
      }
    },
    async created() {
      try {
        this.pending = true
        await this.getDetails()
      } catch (error) {
        console.log(error)
      } finally {
        this.pending = false
      }
    },
    methods: {
      isModelExisted() {
        const modelId = this.details?.operation_detail?.bk_obj_id
        return Boolean(this.$store.getters['objectModelClassify/getModelById'](modelId))
      },
      show() {
        this.isShow = true
      },
      handleHidden() {
        this.$emit('close')
      },
      async getDetails() {
        try {
          if (this.aduitTarget === 'instance') {
            this.details = await this.$store.dispatch('audit/getInstDetails', {
              params: {
                condition: {
                  bk_biz_id: this.bizId,
                  bk_obj_id: this.objId,
                  resource_type: this.resourceType,
                  id: [this.id]
                },
                with_detail: true
              }
            })
          }

          if (this.aduitTarget === 'common') {
            this.details = await this.$store.dispatch('audit/getDetails', {
              id: this.id
            })
          }
        } catch (error) {
          console.error(error)
          this.details = null
        }
      }
    }
  }
</script>

<style lang="scss" scoped>
    .details-content {
        height: calc(100vh - 60px);
        padding: 20px;
    }
</style>
