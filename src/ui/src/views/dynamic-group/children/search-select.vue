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
  import { useStore } from '@/store'
  import { escapeRegexChar } from '@/utils/util'
  import useSearchUser from '@/hooks/use-search-user'

  const props = defineProps({
    defaultFilter: {
      type: Array,
      default: () => []
    }
  })

  const store = useStore()

  const { search: userSearch, lookup: userLookup } = useSearchUser()

  const emit = defineEmits(['search'])

  const searchSelectComp = ref(null)
  const isTyeing = ref(false)
  const isFocus = ref(false)
  const filterMenus = [
    {
      id: 'id',
      name: t('ID')
    },
    {
      id: 'name',
      name: t('分组名称'),
      multiable: true,
      remote: true
    },
    {
      id: 'bk_obj_id',
      name: t('查询对象'),
      multiable: false,
      children: [
        {
          name: t('主机'),
          id: 'host'
        },
        {
          name: t('集群'),
          id: 'set'
        }
      ]
    }, {
      id: 'modify_user',
      name: t('修改人'),
      multiable: true,
      remote: true
    }
  ]
  const filter = ref([])
  const bizId = computed(() => store.getters['objectBiz/bizId'])

  const displayFilterMenus = computed(() => {
    if (filter.value?.length) {
      return filterMenus.filter(menu => !filter.value.some(item => item.id === menu.id))
    }
    return filterMenus.slice()
  })

  const getFormatVal = (value, id) => {
    if (id === 'bk_obj_id') {
      if (value === 'set') return t('集群')
      if (value === 'host') return t('主机')
    }
    return value
  }
  const fetchOptions = async (val, menu) => {
    const fetchs = {
      name: fetchDynamicGroup,
      modify_user: fetchMember
    }
    return fetchs[menu.id](val, menu)
  }
  const fetchDynamicGroup = async (val, menu) => {
    const params = {
      condition: {
        [menu.id]: escapeRegexChar(val)
      },
      page: {
        start: 0,
        limit: 100,
        sort: 'id'
      }
    }
    if (!isTyeing.value || !val?.length || val === `${menu.name}：`) {
      Reflect.deleteProperty(params, 'condition')
    }
    const { info } = await store.dispatch('dynamicGroup/search', {
      bizId: bizId.value,
      params,
      config: {
        cancelPrevious: true,
        globalPermission: false
      }
    })
    return info
  }
  const fetchMember = async (val, menu) => {
    if (!isTyeing.value || val?.length < 2 || val === `${menu.name}：`) {
      return []
    }
    const result = await userSearch(val)
    return result
  }
  const handleSearchSelectChange = (list) => {
    const ids = filterMenus.map(item => item.id)
    const nameItem = list.find(searchItem => searchItem.id === 'name')
    list.forEach((item) => {
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

  watch(() => props.defaultFilter, async (defaultFilter) => {
    filter.value = defaultFilter.map(item => ({
      id: item.id,
      name: filterMenus.find(menu => menu.id === item.id)?.name,
      values: (getFormatVal(item.value, item.id)?.split(',') || []).map(val => ({ name: val, id: item.value }))
    }))

    // 把用户id回填为名称
    const modifierFilter = filter.value.find(item => item.id === 'modify_user')
    if (modifierFilter) {
      const data = await userLookup(modifierFilter.values.map(item => item.name).join(','))
      modifierFilter.values = data.map(item => ({
        id: item.login_name,
        username: item.bk_username,
        name: item.display_name
      }))
    }
  }, { immediate: true })

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
</script>

<template>
  <bk-search-select
    ref="searchSelectComp"
    class="search-select"
    :clearable="false"
    :placeholder="$t('请输入关键字或选择条件搜索')"
    :filter="true"
    :show-condition="false"
    :show-popover-tag-change="false"
    :data="displayFilterMenus"
    v-model="filter"
    :remote-method="fetchOptions"
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
