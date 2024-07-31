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

const { extname } = require('path')
const { sources } = require('webpack')

const dealAssets = (assets, compilation, config) => {
  Object.entries(assets).forEach(([pathname, source]) => {
    const ext = extname(pathname)
    if (ext === '.css' || ext === '.js') {
      const replacement = {
        '.css': '../',
        // webworker importScript
        '.js': '../'
      }

      const newContent = source.source().replace(
        new RegExp(config.assetsPublicPath, 'g'),
        () => replacement[ext],
      )

      compilation.updateAsset(pathname, new sources.RawSource(newContent))
    }
  })
}

class ReplaceStaticUrlPlugin {
  constructor(buildConfig) {
    this.config = buildConfig
  }

  apply(compiler) {
    compiler.hooks.compilation.tap('ReplaceStaticUrlPlugin', (compilation) => {
      compilation.hooks.processAssets.tap({
        name: 'ReplaceStaticUrlPlugin',
        stage: compiler.webpack.Compilation.PROCESS_ASSETS_STAGE_OPTIMIZE_INLINE,
      }, (assets) => {
        dealAssets(assets, compilation, this.config)
      })
    })
  }
}

module.exports = ReplaceStaticUrlPlugin
