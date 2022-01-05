import { translateAuth } from '@/setup/permission'
import store from '@/store'

export default function applyPermission(auth, action) {
  return new Promise(async (resolve, reject) => {
    try {
      const permission = translateAuth(auth)
      const url = await store.dispatch('auth/getSkipUrl', { params: permission })
      if (!action) {
        window.open(url)
      } else {
        action(url)
      }
      resolve(url)
    } catch (e) {
      reject(e)
    }
  })
}
