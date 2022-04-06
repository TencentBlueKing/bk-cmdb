import { ref } from '@vue/composition-api'
import { find } from '@/service/association/index.js'

/**
 * 加载所有关联关系类型
 * @returns {boolean} loading 类型加载状态
 * @returns {object} associations 所有关联关系类型
 */
export const useAssociations = () => {
  const loading = ref(false)
  const associations = ref([])

  loading.value = true

  find()
    .then(({ info }) => {
      associations.value = info
    })
    .finally(() => {
      loading.value = false
    })

  return {
    loading,
    associations
  }
}
