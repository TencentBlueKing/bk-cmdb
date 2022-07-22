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

/**
 * usage:
 * node ./builder/utils/add-license.js
 * node ./builder/utils/add-license.js --dir dirName
 * node ./builder/utils/add-license.js --debug 1
 */

const fs = require('fs')
const path = require('path')

const parseArgs = require('minimist')

const argv = parseArgs(process.argv.slice(2))

const { dir, debug } = argv

const baseDir = dir || fs.realpathSync(process.cwd())

const searchDir = path.resolve(baseDir, './')

const htmlTypeExts = ['.vue', '.html', '.htm']
const searchFileExts = ['.js', ...htmlTypeExts, '.css', '.scss']

const newLicenseContent = `
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2022 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
`

const oldLicense = `/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and limitations under the License.
 */`

const newLicense = `/*${newLicenseContent} */`

const newLicenseForHTML = `<!--${newLicenseContent}-->`

const matched = []

function searchAndReplace(targetDir) {
  function find(dir) {
    const dirList = fs.readdirSync(dir, { withFileTypes: true })
    dirList.forEach((dirent) => {
      if (dirent.isDirectory() && dirent.name !== 'node_modules') {
        find(path.join(dir, dirent.name))
      } else if (dirent.isFile()) {
        const filepath = path.join(dir, dirent.name)
        const fileext = path.extname(filepath)

        if (!searchFileExts.includes(fileext)) {
          return
        }

        const content = fs.readFileSync(filepath, { encoding: 'utf8' })

        // 已添加则不处理
        if (content.startsWith(oldLicense) || content.startsWith(newLicense) || content.startsWith(newLicenseForHTML)) {
          return
        }

        const writeLicenseContent = htmlTypeExts.includes(fileext) ? newLicenseForHTML : newLicense

        // 将license内容添加到原内容的头部
        const newContent = `${writeLicenseContent}\n\n${content}`

        // 写回原文件
        fs.writeFileSync(filepath, newContent, { encoding: 'utf8' })

        if (debug) {
          matched.push(filepath)
        }
      }
    })
  }

  find(targetDir)
}

try {
  searchAndReplace(searchDir)
  console.log('✅ everything is ok')
} catch (err) {
  console.error(err)
} finally {
  if (debug) {
    console.log(matched, matched.length)
  }
}
