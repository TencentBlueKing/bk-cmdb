<template>
  <div class="setting-tags" v-show="customized">
    <div class="tag-item target">检索对象：{{targetScopes}}</div>
    <div class="tag-item" v-for="(item, index) in targetModels" :key="index" v-bk-overflow-tips>
      {{item.targetName}}：{{item.models}}
    </div>
    <bk-button class="reset" text size="small"
      @click="handleClear">
      {{$t('清空所有')}}
    </bk-button>
  </div>
</template>

<script>
  import { computed, defineComponent } from '@vue/composition-api'
  import { targetMap, currentSetting, handleReset } from './use-advanced-setting.js'
  import useRoute, { pickQuery } from './use-route.js'

  export default defineComponent({
    setup(props, { root }) {
      const { $store, $routerActions } = root
      const { route } = useRoute(root)

      const getModelById = $store.getters['objectModelClassify/getModelById']
      const getModelName = id => getModelById(id)?.bk_obj_name ?? '--'

      const targetModels = computed(() => {
        const targetModels = []
        currentSetting.targets.forEach((target) => {
          const modelIds = currentSetting[`${target}s`]
          targetModels.push({
            targetName: targetMap[target],
            models: modelIds.length ? modelIds.map(id => getModelName(id)).join(' | ') : root.$t('全部')
          })
        })
        return targetModels
      })

      const customized = computed(() => {
        const changedModels = []
        currentSetting.targets.forEach((target) => {
          const modelIds = currentSetting[`${target}s`]
          changedModels.push(modelIds.length > 0)
        })
        return changedModels.some(changed => changed)
      })

      const targetScopes = computed(() => currentSetting.targets.map(target => targetMap[target]).join(' | '))

      const handleClear = () => {
        handleReset()
        const query = pickQuery(route.value.query, ['tab', 'keyword'])
        $routerActions.redirect({
          query: {
            ...query,
            t: Date.now()
          }
        })
      }

      return {
        customized,
        targetMap,
        targetScopes,
        targetModels,
        currentSetting,
        handleClear
      }
    }
  })
</script>

<style lang="scss" scoped>
  .setting-tags {
    display: flex;
    font-size: 12px;
    align-items: center;

    .tag-item {
      background: #dcdee5;
      border-radius: 2px;
      padding: 3px 8px;
      max-width: 500px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;

      & + .tag-item {
        margin-left: 6px;
      }

      &.target {
        flex: none;
      }

      .sub-item {
        & + .sub-item {
          &::before {
            content: "|";
            padding: 0 3px;
          }
        }
      }
    }

    .reset {
      flex: none;
    }
  }
</style>
