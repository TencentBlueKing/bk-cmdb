import http from '@/api'

const create = ({ file, params, config, bk_obj_id: objId }) => {
  const form = new FormData()
  form.append('file', file)
  form.append('params', JSON.stringify(params))
  return http.post(`${window.API_HOST}insts/object/${objId}/import`, form, config)
}

const update = create

export default {
  create,
  update
}
