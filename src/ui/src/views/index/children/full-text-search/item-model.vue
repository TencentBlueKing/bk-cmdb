<template>
  <div class="result-item">
    <div class="result-title" @click="data.linkTo(data.source)">
      <span v-html="`${data.typeName} - ${data.title}`"></span>
    </div>
    <div class="result-desc" @click="data.linkTo(data.source)">
      <div class="desc-item"
        v-html="`${$t('模型ID')}：${getHighlightValue(data.source.bk_obj_id, data, 'bk_obj_id')}`"></div>
      <div class="desc-item">{{$t('所属模型分组')}}：{{classificationName}}</div>
      <dl class="model-group-list">
        <div class="group" v-for="(group, index) in groupedProperties" :key="index">
          <dt class="group-name">{{$t('模型字段')}}</dt>
          <dd class="property-list">
            <div class="property-item" v-for="(property, childIndex) in group.properties" :key="childIndex">
              {{property.bk_property_name}}（{{fieldTypeMap[property.bk_property_type]}}）
            </div>
          </dd>
        </div>
      </dl>
    </div>
  </div>
</template>

<script>
  import { defineComponent, toRefs, computed, watchEffect, ref } from '@vue/composition-api'
  import { getText, getHighlightValue } from './use-item.js'

  export default defineComponent({
    name: 'item-model',
    props: {
      data: {
        type: Object,
        default: () => ({})
      },
      propertyMap: {
        type: Object,
        default: () => ({})
      }
    },
    setup(props, { root }) {
      const { $store } = root
      const { data, propertyMap } = toRefs(props)

      const objId = computed(() => data.value.source.bk_obj_id)
      const properties = computed(() => propertyMap.value[objId.value])

      const classificationName = computed(() => {
        const classifications = $store.getters['objectModelClassify/classifications']
        const id = data.value.source.bk_classification_id
        return (classifications.find(item => item.bk_classification_id === id) || {}).bk_classification_name
      })

      const propertyGroups = ref([])
      watchEffect(async () => {
        propertyGroups.value = await $store.dispatch('objectModelFieldGroup/searchGroup', {
          objId: objId.value,
          config: {
            requestId: `get_searchGroup_${objId.value}`,
            fromCache: true,
            cancelPrevious: true
          }
        })
      })

      // eslint-disable-next-line max-len
      const sortProperties = computed(() => (properties.value || []).sort((propertyA, propertyB) => propertyA.bk_property_index - propertyB.bk_property_index))

      // eslint-disable-next-line max-len
      const sortedPropertyGroups = computed(() => propertyGroups.value.sort((groupA, groupB) => groupA.bk_group_index - groupB.bk_group_index))

      const groupedProperties = computed(() => sortedPropertyGroups.value.map(group => ({
        group,
        properties: sortProperties.value.filter((property) => {
          if (['default', 'none'].includes(property.bk_property_group) && group.bk_group_id === 'default') {
            return true
          }
          return property.bk_property_group === group.bk_group_id
        })
      })))

      const fieldTypeMap = {
        singlechar: root.$t('短字符'),
        int: root.$t('数字'),
        float: root.$t('浮点'),
        enum: root.$t('枚举'),
        date: root.$t('日期'),
        time: root.$t('时间'),
        longchar: root.$t('长字符'),
        objuser: root.$t('用户'),
        timezone: root.$t('时区'),
        bool: 'bool',
        list: root.$t('列表'),
        organization: root.$t('组织')
      }

      return {
        properties,
        groupedProperties,
        fieldTypeMap,
        classificationName,
        getText,
        getHighlightValue
      }
    }
  })
</script>

<style lang="scss" scoped>
  .model-group-list {
    display: flex;
    flex-wrap: wrap;
    .group {
      display: flex;
      margin-bottom: 6px;

      .group-name {
        flex: none;
        &::after {
          content: '：';
        }
      }

      .property-list {
        display: flex;
        flex-wrap: wrap;

        .property-item {
          flex: none;
          margin-right: 8px;
        }
      }
    }
  }
</style>
