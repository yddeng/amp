import router from './router'
import storage from 'store'
import { ACCESS_TOKEN } from '@/store/mutation-types'

const loginRoutePath = '/user/login'

router.beforeEach((to, from, next) => {
  if (storage.get(ACCESS_TOKEN)) {
    next()
  } else {
    if (to.path === loginRoutePath) {
      next()
    } else {
      next({ path: '/user/login', query: { redirect: to.fullPath } })
    }
  }
})
