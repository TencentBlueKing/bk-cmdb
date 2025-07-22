/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017 Tencent. All rights reserved.
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
const crypto = require('crypto')

class BuildHashPlugin {
  constructor(options = {}) {
    this.options = Object.assign(
      {
        outputFile: 'build-hash.txt', // 默认输出文件名
        algorithm: 'sha256', // 哈希算法
        includeTimestamp: false, // 是否包含时间戳
        customContent: null, // 自定义内容，会添加到哈希中
      },
      options,
    )
  }

  apply(compiler) {
    // 在完成构建后执行
    compiler.hooks.done.tap('BuildHashPlugin', (stats) => {
      const { compilation } = stats
      const outputPath = compilation.outputOptions.path || compiler.outputPath
      const hashFile = path.resolve(outputPath, this.options.outputFile)

      // 1. 创建哈希生成器
      const hash = crypto.createHash(this.options.algorithm)

      // 2. 添加构建时间戳（可选）
      if (this.options.includeTimestamp) {
        hash.update(Date.now().toString())
      }

      // 3. 添加自定义内容（可选）
      if (typeof this.options.customContent === 'string') {
        hash.update(this.options.customContent)
      }

      // 4. 添加编译信息
      hash.update(compilation.hash)

      // 5. 添加所有依赖模块的修改时间
      compilation.fileDependencies.forEach((filePath) => {
        if (fs.existsSync(filePath)) {
          const stats = fs.statSync(filePath)
          hash.update(stats.mtimeMs.toString())
        }
      })

      // 6. 生成最终哈希
      const finalHash = hash.digest('hex')

      // 7. 确保输出目录存在
      if (!fs.existsSync(outputPath)) {
        fs.mkdirSync(outputPath, { recursive: true })
      }

      // 8. 写入文件
      fs.writeFileSync(hashFile, finalHash)

      // 9. 在控制台显示信息（可选）
      if (this.options.verbose !== false) {
        console.log(`\nBuild hash generated: ${finalHash}`)
        console.log(`Hash file written to: ${hashFile}`)
      }
    })
  }
}

module.exports = BuildHashPlugin
