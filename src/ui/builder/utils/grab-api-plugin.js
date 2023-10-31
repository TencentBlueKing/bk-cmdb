/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

const fs = require('node:fs')
const { resolveBase } = require('./index')

class GrabAPIPlugin {
  constructor() {
    this.pathfile = 'apipaths.txt'
    this.varfile = 'apivars.txt'
  }
  apply(compiler) {
    compiler.hooks.initialize.tap('GrabAPIPlugin', () => {
      fs.rmSync(this.pathfile, { force: true })
      fs.rmSync(this.varfile, { force: true })
      fs.rmSync(resolveBase('node_modules/.cache'), { force: true, recursive: true })
    })

    compiler.hooks.done.tapAsync('GrabAPIPlugin', (stats, callback) => {
      if (stats.hasErrors()) {
        return callback()
      }

      if (fs.existsSync(this.pathfile)) {
        const ignorePaths = [
          '${url}',
          '${window.API_HOST}${urlSuffix[type]}',
          'hosts/snapshot/asstdetail',
          'object/statistics',
          'objectatt/group/property',
          'hosts/add',
          'hosts/${rootGetters.supplierAccount}/${hostId}',
          'hosts/snapshot/${hostId}',
          'topoinstchild/object/${bkObjId}/biz/${bkBizId}/inst/${bkInstId}',
          'find/objecttopology',
          'find/objectclassification',
          'find/topoassociationtype',
          'biz/${bkSupplierAccount}/${bkBizId}',

          // 重复
          'biz/search/${rootGetters.supplierAccount}',

          'create/proc/service_instance/preview',
          'deletemany/proc/service_instance/preview',
          'find/proc/service_instance/difference',
          'findmany/topo/service_template_sync_status/bk_biz_id/${params.bk_biz_id}',

          // 重复
          'host/transfer_with_auto_clear_service_instance/bk_biz_id/{this.bizId}',

          'delete/objectattgroupasst/object/${objId}/property/${propertyId}/group/${groupId}',
          '${window.API_HOST}object/owner/${rootGetters.supplierAccount}/object/${objId}/export',
          '${window.API_HOST}object/owner/${rootGetters.supplierAccount}/object/${objId}/import',

          // 重复
          '${window.API_HOST}object/object/${this.activeModel.bk_obj_id}/import',
          'delete/resource/directory/${moduleId}',
          'find/objectattgroup/object/${options.bk_obj_id}',
          'find/objectunique/object/${objId}',
          'host/transfer_with_auto_clear_service_instance/bk_biz_id/${this.bizId}',
          'update/resource/directory/${moduleId}'
        ]
        const appendPaths = [
          'hosts/import',
          'hosts/update',
          'findmany/inst/association/object/${currentModelId}/inst_id/${currentInstId}/offset/${offset}/limit/${limit}/web',
          'findmany/hosts/search/with_biz',
          'findmany/hosts/search/resource',
          'findmany/hosts/search/noauth'
        ]
        const orders = {
          拓扑: ['topoinstnode', '^topo', 'topomodelmainline'],
          权限: ['^auth'],
          业务查询: ['^biz'],
          项目: ['/project'],
          业务集: ['/biz_set'],
          '管控区域/云账户/云资源发现': ['cloud/', '/cloudarea'],
          动态分组: ['^dynamicgroup'],
          '服务模板/服务实例/进程/主机自动应用': ['proc/proc_template', 'proc/service_category', 'proc/service_instance', 'service_template', 'host_apply_plan', 'host_apply_rule'],
          容器数据纳管: ['kube'],
          字段组合模板: ['field_template']
        }
        const orderEntries = Object.entries(orders)
        const group = {
          其它: []
        }

        const content = fs.readFileSync('apipaths.txt', 'utf-8')
        const list = content.split(/\n/)
          .map(x => x.replaceAll(/[`']/g, ''))
          .filter(x => x.length > 0 && !ignorePaths.includes(x))
          .concat(appendPaths)
          .sort()
        const uniqueList = [...new Set(list)]

        uniqueList.forEach((path) => {
          for (const [key, searchs] of orderEntries) {
            if (searchs.some(x => (x.startsWith('^') ? path.indexOf(x.substring(1)) === 0 : path.indexOf(x) > -1))) {
              if (group[key]) {
                group[key].push(path)
              } else {
                group[key] = [path]
              }
            }
          }
        })

        const grouped = Object.values(group).reduce((acc, cur) => acc.concat(cur), [])
        uniqueList.forEach((path) => {
          if (!grouped.includes(path)) {
            group['其它'].push(path)
          }
        })

        const newlist = []
        Object.keys(orders).concat('其它')
          .forEach((key) => {
            newlist.push(`${key}：`)
            newlist.push(...group[key])
            newlist.push('')
          })

        fs.writeFileSync(this.pathfile, newlist.join('\n'))
      }

      if (fs.existsSync(this.varfile)) {
        const content = fs.readFileSync(this.varfile, 'utf-8')
        const list = content.split(/\n/).filter(x => x.length > 0)
          .sort()

        fs.writeFileSync(this.varfile, [...new Set(list)].join('\n'))
      }

      callback()
    })
  }
}

module.exports = GrabAPIPlugin
