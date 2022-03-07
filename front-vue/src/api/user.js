import request from '@/utils/request'

const userApi = {
  userList: 'user/list',
  userAdd: 'user/add',
  userDelete: 'user/delete'
}

/**
 * login func
 * parameter: {
 *     username: '',
 *     password: '',
 *     remember_me: true,
 *     captcha: '12345'
 * }
 * @param parameter
 * @returns {*}
 */
export function userList (parameter) {
  return request({
    url: userApi.userList,
    method: 'post',
    data: parameter
  })
}

export function userAdd (parameter) {
  return request({
    url: userApi.userAdd,
    method: 'post',
    data: parameter
  })
}

export function userDelete (parameter) {
  return request({
    url: userApi.userDelete,
    method: 'post',
    data: parameter
  })
}
