import userQueue from '@/components/ui/user/user-queue.js'

const setTextContent = (el, binding) => {
    const user = binding.value
    const modifiers = binding.modifiers
    if ([undefined, null, '', '--'].includes(user)) {
        el.textContent = user
        modifiers.title && (el.title = user)
        return
    }
    const userInfo = userQueue.getUserInfo(user)
    if (userInfo) {
        el.textContent = userInfo
        modifiers.title && (el.title = user)
    } else {
        userQueue.addUser({
            user,
            node: el,
            options: modifiers
        })
    }
}

const user = {
    inserted (el, binding, vnode) {
        setTextContent(el, binding)
    },
    update (el, binding) {
        if (binding.value === binding.oldValue) return
        setTextContent(el, binding)
    }
}
export default {
    install: Vue => Vue.directive('user', user)
}
