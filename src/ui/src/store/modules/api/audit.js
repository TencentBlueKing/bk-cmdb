import $http from '@/api'
const actions = {
  getList(context, { params, config }) {
    return $http.post('findmany/audit_list', params, config)
  },
  getDictionary(context, config) {
    return $http.get('find/audit_dict', config)
  },
  getDetails(context, { id, config }) {
    return $http.post('find/audit', { id: [id] }, config).then(([detail]) => detail)
  },
  getInstList(context, { params, config }) {
    return $http.post('find/inst_audit', params, config)
  },
  getInstDetails(context, { params, config }) {
    return $http.post('find/inst_audit', params, config).then(({ info: [detail] }) => detail)
  }
}

export default {
  namespaced: true,
  actions
}
