import { ref, watch } from 'vue'
import { t } from '@/i18n'
import { $bkInfo } from '@/magicbox/index.js'

export default function useSideslider(data, options = {}) {
  const isChanged = ref(false)
  const { watchOnce = true } = options

  const unwatch = watch(data, () => {
    isChanged.value = true
    watchOnce && unwatch()
  }, { deep: true })


  const beforeClose = (confirmCallback, cancelCallback) => new Promise((resolve, reject) => {
    if (!isChanged.value) {
      confirmCallback && confirmCallback?.()
      resolve(true)
      return
    }
    $bkInfo({
      title: t('确认离开当前页？'),
      subTitle: t('离开将会导致未保存信息丢失'),
      clsName: 'custom-info-confirm default-info',
      okText: t('离开'),
      cancelText: t('取消'),
      confirmFn() {
        confirmCallback && confirmCallback?.()
        resolve(true)
      },
      cancelFn() {
        cancelCallback && cancelCallback?.()
        reject(false)
      },
    })
  })

  const reset = () => {
    setTimeout(() => {
      isChanged.value = false
    })
  }

  const setChanged = (v) => {
    isChanged.value = v
  }

  return {
    beforeClose,
    isChanged,
    reset,
    setChanged
  }
}
