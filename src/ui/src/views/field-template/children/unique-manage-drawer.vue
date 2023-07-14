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

<script setup>
  import { computed, ref, watch } from 'vue'
  import cloneDeep from 'lodash/cloneDeep'
  import { t } from '@/i18n'
  import { $bkInfo } from '@/magicbox'
  import UniqueManage from './unique-manage.vue'

  const props = defineProps({
    open: {
      type: Boolean,
      default: false
    },
    uniqueList: {
      type: Array,
      default: () => ([])
    },
    fieldList: {
      type: Array,
      default: () => ([])
    },
    beforeUniqueList: {
      type: Array,
      default: () => ([])
    }
  })

  const emit = defineEmits(['close', 'change-unique'])

  const sidesliderComp = ref(null)
  const uniqueManageComp = ref(null)

  const uniqueListLocal = ref(cloneDeep(props.uniqueList))
  watch(() => props.uniqueList, (uniqueList) => {
    uniqueListLocal.value = cloneDeep(uniqueList)
  }, { deep: true })

  const isShow = computed({
    get() {
      return props.open
    },
    set() {
      emit('close')
    }
  })

  const isAllValid = ref(true)
  const editedUniqueList = ref(null)
  const hasEmptyUnique = computed(() => {
    if (editedUniqueList.value === null || !editedUniqueList.value?.length) {
      return false
    }
    return editedUniqueList.value.some(item => !item.keys?.length)
  })
  const handleUpdateUnique = (uniqueList, isValid) => {
    isAllValid.value = isValid
    editedUniqueList.value = uniqueList

    if (isValid) {
      emit('change-unique', uniqueList)
    }
  }

  const handleSliderHidden = () => {
    // 必要的重置数据
    isAllValid.value = true
    editedUniqueList.value = null
    emit('close')
  }

  const handleSliderBeforeClose = async () => {
    if (hasEmptyUnique.value || !isAllValid.value) {
      const subTitle = t(hasEmptyUnique.value ? '未编辑完成的数据将被忽略，原数据不会被更改' : '校验未通过，数据不会被更新')
      return new Promise((resolve) => {
        $bkInfo({
          title: t('确认退出'),
          subTitle,
          extCls: 'bk-dialog-sub-header-center',
          confirmFn: () => {
            emit('close')
            resolve(true)
          },
          cancelFn: () => {
            resolve(false)
          }
        })
      })
    }
    return true
  }
</script>

<template>
  <bk-sideslider
    ref="sidesliderComp"
    v-transfer-dom
    :width="460"
    :title="$t('唯一校验')"
    :is-show.sync="isShow"
    :quick-close="true"
    :before-close="handleSliderBeforeClose"
    @hidden="handleSliderHidden">
    <div slot="content" class="content">
      <unique-manage
        :field-list="fieldList"
        :unique-list="uniqueListLocal"
        :before-unique-list="beforeUniqueList"
        @update-unique="handleUpdateUnique"
        ref="uniqueManageComp">
      </unique-manage>
    </div>
  </bk-sideslider>
</template>

<style lang="scss" scoped>
  .content {
    padding: 24px;
  }
</style>
