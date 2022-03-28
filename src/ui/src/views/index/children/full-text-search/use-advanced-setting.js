import { reactive, shallowReactive, ref } from '@vue/composition-api'

export const targetMap = {
  model: '模型',
  instance: '实例'
}

const defaultSetting = {
  targets: Object.keys(targetMap),
  models: [],
  instances: []
}

export const currentSetting = shallowReactive(defaultSetting)

export const allModelIds = ref([])

export const handleReset = () => {
  currentSetting.targets = Object.keys(targetMap),
  currentSetting.models = []
  currentSetting.instances = []
}

export default function useAdvancedSetting(options = {}, root) {
  const { onConfirm, onShow, onCancel } = options

  const activeSetting = reactive(defaultSetting)

  const handleShow = () => {
    activeSetting.targets = currentSetting.targets
    activeSetting.models = currentSetting.targets.includes('model') ? currentSetting.models : []
    activeSetting.instances = currentSetting.targets.includes('instance') ? currentSetting.instances : []
    onShow && onShow()
  }

  const handleConfirm = () => {
    currentSetting.targets = activeSetting.targets
    currentSetting.models = activeSetting.models
    currentSetting.instances = activeSetting.instances
    onConfirm && onConfirm()
  }

  const handleCancel = () => {
    onCancel && onCancel()
  }

  const handleTargetClick = (value) => {
    if (activeSetting.targets.includes(value)) {
      activeSetting.targets.length > 1 && activeSetting.targets.splice(activeSetting.targets.indexOf(value), 1)
    } else {
      activeSetting.targets.push(value)
    }
  }

  // 获取所有的模型ID值
  allModelIds.value = []
  const classifications = root.$store.getters['objectModelClassify/classifications']
  const displayModelList = []
  classifications.forEach((classification) => {
    displayModelList.push({
      ...classification,
      bk_objects: classification.bk_objects.filter(model => !model.bk_ispaused && !model.bk_ishidden)
    })
  })
  displayModelList.forEach(model => allModelIds.value.push(...model.bk_objects.map(m => m.bk_obj_id)))

  return {
    activeSetting,
    handleShow,
    handleConfirm,
    handleCancel,
    handleTargetClick
  }
}
