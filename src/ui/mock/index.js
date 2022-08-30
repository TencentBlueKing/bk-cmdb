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

const fs = require('fs')
const path = require('path')

const baseDir = path.resolve(__dirname, './')
const selfpath = path.resolve(__filename, './')

function load(dir) {
  // 所有mock定义的集合
  const allDefs = {}

  const dirList = fs.readdirSync(dir, { withFileTypes: true })
  dirList.forEach((dirent) => {
    if (dirent.isDirectory()) {
      const defs = load(path.join(dir, dirent.name))
      Object.assign(allDefs, defs)
    } else if (dirent.isFile()) {
      const filepath = path.join(dir, dirent.name)

      // 忽略本文件
      if (filepath === selfpath) {
        return
      }

      // 忽略非index.js入口文件
      if (dirent.name !== 'index.js') {
        return
      }

      try {
        delete require.cache[filepath]
        const def = require(filepath)

        // 解析处理def
        Object.keys(def).forEach((key) => {
          if (def[key].path) {
          // 限制path只能相对于入口文件
            def[key].fullpath = path.resolve(dir, './', def[key].path)
          }
        })

        // 合并到一起
        Object.assign(allDefs, def)
      } catch (err) {
        console.error(err)
      }
    }
  })

  return allDefs
}

function getDefs() {
  // 每一次重新load，确保为最新的定义
  return load(baseDir)
}

module.exports = {
  getDefs
}
