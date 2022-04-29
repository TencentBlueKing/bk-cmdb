import store from '@/store'
import router from '@/router'
import { t } from '@/i18n'

/**
 * 更改文档标题
 * @param {Array} [appendTitles] 追加的标题，会展示在默认名称之后。不传入时会根据当前路由重新生成路径。
 */
export const changeDocumentTitle = (appendTitles = []) => {
  const { name, separator } = store.state.globalConfig.config.site
  const { matched } = router.currentRoute
  let matchedNames = [name]
  matched.forEach((match) => {
    if (match?.meta?.menu?.i18n) {
      matchedNames.push(t(match.meta.menu.i18n))
    }
  })

  if (appendTitles?.length) {
    matchedNames = matchedNames.concat(appendTitles)
  }

  document.title = matchedNames?.join(` ${separator} `)
}
