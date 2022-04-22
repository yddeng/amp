import request from '@/utils/request'

export function post (url, parameter) {
  return request({
    url: url,
    method: 'post',
    data: parameter
  })
}
