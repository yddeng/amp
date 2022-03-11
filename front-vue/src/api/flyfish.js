import request from '@/utils/request'

const flyfishApi = {
  getMeta: 'flyfish/getMeta',
  addTable: 'flyfish/addTable',
  addField: 'flyfish/addField',
  getSetStatus: 'flyfish/getSetStatus'
}

export function getMeta (parameter) {
  return request({
    url: flyfishApi.getMeta,
    method: 'post',
    data: parameter
  })
}

export function addTable (parameter) {
  return request({
    url: flyfishApi.addTable,
    method: 'post',
    data: parameter
  })
}

export function addField (parameter) {
  return request({
    url: flyfishApi.addField,
    method: 'post',
    data: parameter
  })
}

export function getSetStatus (parameter) {
  return request({
    url: flyfishApi.getSetStatus,
    method: 'post',
    data: parameter
  })
}

export function setMarkClear (parameter) {
  return request({
    url: '/flyfish/setMarkClear',
    method: 'post',
    data: parameter
  })
}

export function addSet (parameter) {
  return request({
    url: '/flyfish/addSet',
    method: 'post',
    data: parameter
  })
}

export function remSet (parameter) {
  return request({
    url: '/flyfish/remSet',
    method: 'post',
    data: parameter
  })
}

export function addNode (parameter) {
  return request({
    url: '/flyfish/addNode',
    method: 'post',
    data: parameter
  })
}

export function remNode (parameter) {
  return request({
    url: '/flyfish/remNode',
    method: 'post',
    data: parameter
  })
}

export function addLeaderStoreToNode (parameter) {
  return request({
    url: '/flyfish/addLeaderStoreToNode',
    method: 'post',
    data: parameter
  })
}

export function removeNodeStore (parameter) {
  return request({
    url: '/flyfish/removeNodeStore',
    method: 'post',
    data: parameter
  })
}

export function clearDBData (parameter) {
  return request({
    url: '/flyfish/clearDBData',
    method: 'post',
    data: parameter
  })
}

export function drainKv (parameter) {
  return request({
    url: '/flyfish/drainKv',
    method: 'post',
    data: parameter
  })
}
