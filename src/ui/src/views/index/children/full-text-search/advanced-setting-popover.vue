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
  <bk-popover class="search-dropdown" v-bind="popoverProps" ref="popover">
    <bk-link theme="default" class="anchor"
      icon="bk-icon icon-angle-down"
      icon-placement="right"
      @click="handleShow">{{$t('高级搜索')}}</bk-link>
    <template #content>
      <div class="advanced-search-form">
        <div class="setting-group">
          <div class="title">{{$t('检索对象')}}</div>
          <div class="content">
            <div class="bk-button-group">
              <bk-button
                @click="handleTargetClick('model')"
                :class="{ 'is-selected': targets.includes('model') }">{{$t('模型')}}</bk-button>
              <bk-button
                @click="handleTargetClick('instance')"
                :class="{ 'is-selected': targets.includes('instance') }">{{$t('实例')}}</bk-button>
            </div>
          </div>
        </div>
        <div class="setting-group scope">
          <div class="title">{{$t('模型范围')}}</div>
          <div class="content">
            <div class="setting-item" v-show="targets.includes('model')">
              <label class="label">{{$t('模型')}}</label>
              <model-selector multiple searchable class="form-el" :placeholder="$t('默认全部')" v-model="models" />
            </div>
            <div class="setting-item" v-show="targets.includes('instance')">
              <label class="label">{{$t('实例')}}</label>
              <model-selector multiple searchable class="form-el" :placeholder="$t('默认全部')" v-model="instances" />
            </div>
          </div>
        </div>
        <div class="buttons">
          <bk-button theme="primary" @click="handleConfirm">{{$t('确定')}}</bk-button>
          <bk-button theme="default" @click="handleCancel">{{$t('取消')}}</bk-button>
        </div>
      </div>
    </template>
  </bk-popover>
</template>

<script>
  import { defineComponent, computed, toRefs, ref } from 'vue'
  import routerActions from '@/router/actions'
  import RouterQuery from '@/router/query'
  import ModelSelector  from '@/components/ui/selector/model.vue'
  import useAdvancedSetting from './use-advanced-setting.js'
  import { pickQuery } from './use-route.js'

  export default defineComponent({
    components: {
      ModelSelector
    },
    setup() {
      const popover = ref(null)
      const route = computed(() => RouterQuery.route)

      const popoverProps = {
        width: 500,
        trigger: 'click',
        distance: 12,
        // sticky: true,
        theme: 'light',
        placement: 'bottom',
        trigger: 'manual',
        tippyOptions: {
          hideOnClick: false
        }
      }

      const {
        activeSetting,
        handleShow,
        handleConfirm,
        handleCancel,
        handleTargetClick
      } = useAdvancedSetting({
        onShow() {
          popover.value.showHandler()
        },
        onConfirm() {
          const query = pickQuery(route.value.query, ['tab', 'keyword'])
          routerActions.redirect({
            query: {
              ...query,
              t: Date.now()
            }
          })
          popover.value.hideHandler()
        },
        onCancel() {
          popover.value.hideHandler()
        }
      })

      return {
        ...toRefs(activeSetting),
        popover,
        popoverProps,
        handleShow,
        handleTargetClick,
        handleConfirm,
        handleCancel
      }
    }
  })
</script>

<style lang="scss" scoped>
  .advanced-search-form {
    margin: 8px;

    .buttons {
      .bk-button {
        & + .bk-button {
          margin-left: 4px;
        }
      }
    }
  }

  .search-dropdown {
    .anchor {
      /deep/ .bk-link-text {
        font-size: 12px;
      }
    }
  }

  .bk-button-group {
    .bk-button {
      min-width: 70px;
      border-radius: 2px;
      & + .bk-button {
        margin-left: 6px;
      }
    }
  }

  .setting-group {
    .title {
      font-size: 14px;
      font-weight: 700;
      color: #63656E;
    }
    .content {
      margin: 14px 0 24px 0;
    }
    &.scope {
      .setting-item {
        display: flex;
        align-items: center;
        margin-bottom: 10px;

        .label {
          flex: none;
          height: 32px;
          width: 80px;
          line-height: 32px;
          text-align: center;
          color: #63656E;
          background: #f0f1f5;
          border: 1px solid #c4c6cc;
          border-radius: 2px 0px 0px 2px;

          & + .form-el {
            margin-left: -1px;
          }
        }
        .form-el {
          flex: 1;
          width: calc(100% - 80px);
        }
      }
    }
  }
</style>
