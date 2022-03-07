import request from '@/utils/request'

const cmdApi = {
  cmdList: 'cmd/list',
  cmdCreate: 'cmd/create',
  cmdDelete: 'cmd/delete',
  cmdUpdate: 'cmd/update',
  cmdExec: 'cmd/exec',
  cmdLog: 'cmd/log'
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
export function cmdList (parameter) {
  return request({
    url: cmdApi.cmdList,
    method: 'post',
    data: parameter
  })
}

export function cmdCreate (parameter) {
  return request({
    url: cmdApi.cmdCreate,
    method: 'post',
    data: parameter
  })
}

export function cmdDelete (parameter) {
  return request({
    url: cmdApi.cmdDelete,
    method: 'post',
    data: parameter
  })
}

export function cmdUpdate (parameter) {
  return request({
    url: cmdApi.cmdUpdate,
    method: 'post',
    data: parameter
  })
}

export function cmdExec (timeout, parameter) {
  return request({
    url: cmdApi.cmdExec,
    method: 'post',
    timeout: timeout,
    data: parameter
  })
}

export function cmdLog (parameter) {
  return request({
    url: cmdApi.cmdLog,
    method: 'post',
    data: parameter
  })
}
