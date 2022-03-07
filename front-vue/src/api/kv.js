import request from '@/utils/request'

const kvApi = {
  kvSet: 'kv/set',
  kvGet: 'kv/get',
  kvDelete: 'kv/delete'
}

export function kvSet (parameter) {
  return request({
    url: kvApi.kvSet,
    method: 'post',
    data: parameter
  })
}

export function kvGet (parameter) {
  return request({
    url: kvApi.kvGet,
    method: 'post',
    data: parameter
  })
}

export function kvDelete (parameter) {
  return request({
    url: kvApi.kvDelete,
    method: 'post',
    data: parameter
  })
}
