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
  <Teleport to=".breadcrumbs-layout">
    <cmdb-auth tag="div" :auth="{ type: $OPERATION.R_MODEL, relation: [modelId] }">
      <template #default="{ disabled }">
        <i class="icon-cc-share share "
          :class="{ disabled }"
          v-bk-tooltips="$t('前往模型管理')"
          @click="handleModelDetail">
        </i>
      </template>
    </cmdb-auth>
  </Teleport>
</template>

<script>
  import { mapGetters } from 'vuex'
  import {
    MENU_MODEL_DETAILS,
  } from '@/dictionary/menu-symbol'
  import Teleport from 'vue2-teleport'

  export default {
    name: 'cmdb-model-fast-link',
    components: {
      Teleport,
    },
    props: {
      objId: {
        type: String,
        default: ''
      }
    },
    computed: {
      ...mapGetters('objectModelClassify', ['getModelById']),
      modelId() {
        return this.model.id
      },
      model() {
        return this.getModelById(this.objId) || {}
      },
    },
    methods: {
      handleModelDetail() {
        this.$routerActions.open({
          name: MENU_MODEL_DETAILS,
          params: {
            modelId: this.objId,
          }
        })
      }
    }
  }
</script>

<style lang="scss" scoped>
  .auth-box {
    @include space-between;
  }
  .share {
    cursor: pointer;
    font-size: 14px;
    margin-left: 5px;
    color: $primaryColor;
  }
  .disabled {
    color: #c4c6cc;
  }
</style>
