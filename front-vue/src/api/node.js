import request from '@/utils/request'

const nodeApi = {
  nodeList: 'node/list',
  nodeRemove: 'node/remove'
}

export function nodeList (parameter) {
  return request({
    url: nodeApi.nodeList,
    method: 'post',
    data: parameter
  })
}

export function nodeRemove (parameter) {
  return request({
    url: nodeApi.nodeRemove,
    method: 'post',
    data: parameter
  })
}
