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
  import { t } from '@/i18n'

  const props = defineProps({
    defaultFilter: {
      type: Array,
      default: () => []
    }
  })

  const emit = defineEmits(['search'])

  const searchSelectComp = ref(null)
  const isTyeing = ref(false)
  const isFocus = ref(false)
  const filterMenus = [
    {
      id: 'id',
      name: t('id')
    },
    {
      id: 'name',
      name: t('名称')
    },
    {
      id: 'bk_obj_id',
      name: t('查询对象'),
      multiable: false,
      children: [
        {
          name: '主机',
          id: 'host'
        },
        {
          name: '集群',
          id: 'set'
        }
      ]
    }, {
      id: 'modify_user',
      name: t('更新人')
    }
  ]
  const filter = ref([])

  watch(() => props.defaultFilter, (defaultFilter) => {
    filter.value = defaultFilter.map(item => ({
      id: item.id,
      name: filterMenus.find(menu => menu.id === item.id)?.name,
      values: (item.value?.split(',') || []).map(val => ({ name: val }))
    }))
  }, { immediate: true })

  const displayFilterMenus = computed(() => {
    if (filter.value?.length) {
      return filterMenus.filter(menu => !filter.value.some(item => item.id === menu.id))
    }
    return filterMenus.slice()
  })

  watch(displayFilterMenus, () => {
    if (!isFocus.value) {
      return
    }
    // fix菜单项丢失
    searchSelectComp.value?.showMenu()
    setTimeout(() => {
      searchSelectComp.value?.getInputInstance()?.click()
    }, 300)
  })
  const handleSearchSelectChange = (list) => {
    const ids = filterMenus.map(item => item.id)
    list.forEach((item) => {
      const nameItem = list.find(searchItem => searchItem.id === 'name')
      if (!ids.includes(item.id)) {
        // 存在多个就合并
        if (nameItem) {
          nameItem.values.push({ name: item.name })
        } else {
          item.id = 'name'
          item.values = [{ name: item.name }]
        }
      }
    })
    // 合并完之后，被合并的项还会遗留，在这里去掉
    const newList = list.filter(item => ids.includes(item.id))
    emit('search', newList)
  }
  const handleInputChange = () => {
    isTyeing.value = true
  }
  const handleInputFocus = () => {
    isFocus.value = true
  }
  const handleInputClickOutside = () => {
    isFocus.value = false
  }
</script>

<template>
  <bk-search-select
    ref="searchSelectComp"
    class="search-select"
    :clearable="false"
    :placeholder="$t('动态分组查询')"
    :filter="true"
    :show-condition="false"
    :show-popover-tag-change="false"
    :data="displayFilterMenus"
    v-model="filter"
    @input-change="handleInputChange"
    @input-focus="handleInputFocus"
    @input-click-outside="handleInputClickOutside"
    @change="handleSearchSelectChange">
  </bk-search-select>
</template>

<style lang="scss" scoped>
.search-select {
  width: 480px;
  background: #fff;
}
</style>
<style lang="scss">
.bk-search-list {
  max-width: 360px;
  .bk-search-list-menu-item {
    width: 100%;
    .item-name {
      width: 100%;
      @include ellipsis;
    }
  }
}
</style>
