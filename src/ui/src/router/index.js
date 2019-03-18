import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

const router = new Router({
    mode: 'history',
    routes: []
})

router.beforeEach(async (to, from, next) => {
    console.log(router)
    next()
})

router.afterEach((to, from) => {
})

export default router
