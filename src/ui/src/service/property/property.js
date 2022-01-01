import http from '@/api'

function createIdProperty(objId) {
  const keyMap = {
    biz: 'bk_biz_id',
    biz_set: 'bk_biz_set_id',
    host: 'bk_host_id'
  }
  return {
    id: Date.now(),
    bk_obj_id: objId,
    bk_property_id: keyMap[objId] || 'bk_inst_id',
    bk_property_name: 'ID',
    bk_property_index: -1,
    bk_property_type: 'int',
    isonly: true,
    ispre: true,
    bk_isapi: true,
    bk_issystem: true,
    isreadonly: true,
    editable: false,
    bk_property_group: null,
    _is_inject_: true
  }
}

export const find = async (params, config, injectId = false) => {
  try {
    const properties = await http.post('find/objectattr', params, config)

    if (!injectId) {
      return properties
    }

    // eslint-disable-next-line no-underscore-dangle
    if (properties.some(property => property._is_inject_)) {
      return properties
    }

    properties.unshift(createIdProperty(params.bk_obj_id))
    return properties
  } catch (error) {
    console.error(error)
    return []
  }
}

export const findBiz = (injectId = false, config) => {
  const params = {
    bk_obj_id: 'biz',
    bk_supplier_account: window.Supplier.account
  }
  return find(params, config, injectId)
}


export const findBizSet = (injectId = false, config) => {
  const params = {
    bk_obj_id: 'biz_set',
    bk_supplier_account: window.Supplier.account
  }
  return find(params, config, injectId)
}

export default {
  find,
  findBiz,
  findBizSet
}
