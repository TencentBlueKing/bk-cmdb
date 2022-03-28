const storageKey = 'selectedBusinessSet'

// 上次是否使用（访问）了业务集
let recentlyUsed = false

export const setBizSetIdToStorage = (id) => {
  window.localStorage.setItem(storageKey, id)
}

export const getBizSetIdFromStorage = () => window.localStorage.getItem(storageKey)

export const setBizSetRecentlyUsed = used => recentlyUsed = used

export const getBizSetRecentlyUsed = () => recentlyUsed
