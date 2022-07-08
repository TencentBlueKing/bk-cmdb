// 通用
import en from './en.json'
import cn from './cn.json'

// 全局配置
import globalConfigZhCN from '@/views/global-config/i18n/zh-CN.json'
import globalConfigEn from '@/views/global-config/i18n/en.json'

// 模型管理
import modelManageZhCN from '@/views/model-manage/i18n/zh-CN.json'
import modelManageEn from '@/views/model-manage/i18n/en.json'

export default {
  en: {
    ...en,
    ...globalConfigEn,
    ...modelManageEn
  },
  zh_CN: {
    ...cn,
    ...globalConfigZhCN,
    ...modelManageZhCN
  }
}
