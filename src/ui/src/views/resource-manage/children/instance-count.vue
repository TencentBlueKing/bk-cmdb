<template>
  <loading :loading="loading">{{count}}</loading>
</template>

<script>
  import { defineComponent, computed, ref, watchEffect } from '@vue/composition-api'
  import loading from '@/components/loading/index.vue'
  import { instanceCounts } from './use-instance-count.js'

  export default defineComponent({
    components: { loading },
    props: {
      objId: {
        type: String,
        requried: true
      }
    },
    setup(props) {
      const { objId } = props

      // 累积每一次的结果
      const list = computed(() => {
        const list = []
        instanceCounts.value.forEach((item) => {
          list.push(...item)
        })
        return list
      })

      // 找到当前模型实例数据
      const instance = computed(() => list.value.find(inst => inst.bk_obj_id === objId))

      // 确定loading态
      const loading = ref(true)
      watchEffect(() => {
        loading.value = list.value?.findIndex(item => item.bk_obj_id === objId) === -1
      })

      // 当前模型实例count
      const count = computed(() => {
        if (instance.value?.error) {
          return '--'
        }
        return instance.value?.inst_count
      })

      return {
        loading,
        count
      }
    }
  })
</script>
