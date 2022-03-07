import request from '@/utils/request'

const processApi = {
  processList: 'process/list',
  processCreate: 'process/create',
  processUpdate: 'process/update',
  processDelete: 'process/delete',
  processStart: 'process/start',
  processStop: 'process/stop',
  groupList: 'process/glist',
  groupAdd: 'process/gadd',
  groupRemove: 'process/gremove'
}

export function groupList (parameter) {
  return request({
    url: processApi.groupList,
    method: 'post',
    data: parameter
  })
}

export function groupAdd (parameter) {
  return request({
    url: processApi.groupAdd,
    method: 'post',
    data: parameter
  })
}

export function groupRemove (parameter) {
  return request({
    url: processApi.groupRemove,
    method: 'post',
    data: parameter
  })
}

export function processList (parameter) {
  return request({
    url: processApi.processList,
    method: 'post',
    data: parameter
  })
}

export function processCreate (parameter) {
  return request({
    url: processApi.processCreate,
    method: 'post',
    data: parameter
  })
}

export function processUpdate (parameter) {
  return request({
    url: processApi.processUpdate,
    method: 'post',
    data: parameter
  })
}

export function processDelete (parameter) {
  return request({
    url: processApi.processDelete,
    method: 'post',
    data: parameter
  })
}

export function processStart (parameter) {
  return request({
    url: processApi.processStart,
    method: 'post',
    data: parameter
  })
}

export function processStop (parameter) {
  return request({
    url: processApi.processStop,
    method: 'post',
    data: parameter
  })
}
