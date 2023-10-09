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
  <bk-select
    v-bind="$attrs"
    v-model="localValue"
    :scroll-height="420">
    <bk-option-group
      class="model-option-group"
      v-for="(group, index) in displayModelList"
      :name="group.bk_classification_name"
      :key="index">
      <bk-option v-for="option in group.bk_objects"
        :key="option.bk_obj_id"
        :id="option.bk_obj_id"
        :name="option.bk_obj_name">
        <div
          class="option-item-content"
          :title="option.name">
          <div class="text">
            <span class="item-name">{{option.bk_obj_name}}</span>
          </div>
          <template v-if="isShowLink">
            <i class="icon-cc-share link-icon" @click.prevent.stop="handleClickLink(option)"></i>
          </template>
        </div>
      </bk-option>
    </bk-option-group>
  </bk-select>
</template>

<script>
  import { MENU_MODEL_DETAILS } from '@/dictionary/menu-symbol'

  export default {
    props: {
      value: {
        type: [Array, String],
        default: ''
      },
      multiple: {
        type: Boolean,
        default: false
      },
      exclude: {
        type: Array,
        default: () => ([])
      },
      isShowLink: {
        type: Boolean,
        default: true
      }
    },
    data() {
      return {
        classifications: []
      }
    },
    computed: {
      localValue: {
        get() {
          return this.multiple ? (this.value || []) : (this.value || '')
        },
        set(values) {
          this.$emit('input', values)
          this.$emit('change', values)
        }
      },
      displayModelList() {
        const displayModelList = []
        this.classifications.forEach((classification) => {
          displayModelList.push({
            ...classification,
            bk_objects: classification.bk_objects
              .filter(model => !model.bk_ispaused && !model.bk_ishidden && !this.exclude.includes(model.bk_obj_id))
          })
        })
        return displayModelList.filter(item => item.bk_objects.length > 0)
      }
    },
    created() {
      this.getModelList()
    },
    methods: {
      async getModelList() {
        try {
          this.classifications = await this.$store.dispatch('objectModelClassify/searchClassificationsObjects', {
            fromCache: true
          })
        } catch (error) {
          this.classifications = []
        }
      },
      handleClickLink(model) {
        this.$routerActions.open({
          name: MENU_MODEL_DETAILS,
          params: {
            modelId: model.bk_obj_id
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
  .model-option-group {
    /deep/.bk-option-group-name {
      @include ellipsis;
    }

  }
  .option-item-content {
    color: #63656E;
    font-size: 14px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    .text {
      flex: 1;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .link-icon {
      font-size: 12px;
      color: #3A84FF;
      margin-left: 8px;
      display: none;
    }

    &:hover {
      .link-icon {
        display: block;
      }
    }
  }
</style>
