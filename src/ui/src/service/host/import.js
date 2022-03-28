import http from '@/api'

const importFactory = (type, { file, params, config }) => {
  const urlSuffix = {
    create: 'hosts/import',
    update: 'hosts/update'
  }
  const form = new FormData()
  form.append('file', file)
  form.append('params', JSON.stringify(params))
  return http.post(`${window.API_HOST}${urlSuffix[type]}`, form, config)
}
export const create = options => importFactory('create', options)
export const update = options => importFactory('update', options)
export default {
  create,
  update
}
