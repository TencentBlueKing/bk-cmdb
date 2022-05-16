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

export const layout = {
  name: 'dagre',
  rankDir: 'TB',
  ranker: 'tight-tree',
  padding: 60,
  nodeSep: 80,
  rankSep: 100
}

const HIGHLIGHT_COLOR = '#699df4'
const HOVER_COLOR = '#3a84ff'
export const style = [{
  // grabbed画布时
  selector: 'core',
  style: {
    'active-bg-color': '#3c96ff',
    'active-bg-size': '18px'
  }
}, {
  // 有关node样式配置
  selector: 'node',
  style: {
    // 点击时不显示overlay
    'overlay-opacity': 0
  }
}, {
  selector: 'node',
  style: {
    width: 36,
    height: 36,
    // 设置label文本
    label: 'data(name)',
    // label
    color: '#868b97',
    'text-valign': 'bottom',
    'text-halign': 'center',
    'font-size': '14px',
    'text-margin-y': '9px',
    // label换行
    'text-wrap': 'ellipsis',
    'text-max-width': '110px',
    'text-overflow-wrap': 'anywhere',
    // 背景图
    'background-image': 'data(background.unselected)',
    'background-color': '#ffffff',
    'background-fit': 'cover cover',
    'border-width': 1,
    'border-color': '#939393',
    'border-opacity': 0.5,
  }
}, {
  // edge样式配置
  selector: 'edge',
  style: {
    'curve-style': 'bezier',
    'target-arrow-shape': 'triangle-backcurve',
    opacity: 1,
    'arrow-scale': 1.5,
    'line-color': '#c3cdd7',
    'target-arrow-color': '#c3cdd7',
    width: 2,
    // 点击时overlay
    color: '#979ba5',
    'overlay-opacity': 0,
    'font-size': '10px',
    'text-background-opacity': 0.7,
    'text-background-color': '#ffffff',
    'text-background-shape': 'roundrectangle',
    'text-background-padding': 2,
    'text-border-opacity': 0.7,
    'text-border-width': 1,
    'text-border-style': 'solid',
    'text-border-color': '#dcdee5',
    'loop-direction': '45deg',
    'loop-sweep': '90deg',
  }
}, {
  selector: 'edge[direction="none"]', // 无方向
  style: {
    'source-arrow-shape': 'none',
    'target-arrow-shape': 'none'
  }
}, {
  selector: 'edge[direction="bidirectional"]', // 双向
  style: {
    'source-arrow-shape': 'triangle-backcurve',
    'source-arrow-color': '#c3cdd7',
    'target-arrow-shape': 'triangle-backcurve',
    'target-arrow-color': '#c3cdd7'
  }
}, {
  selector: 'edge[direction="src_to_dest"]', // 源指向目标
  style: {
    'source-arrow-shape': 'triangle-backcurve',
    'source-arrow-color': '#c3cdd7',
    'target-arrow-shape': 'none'
  }
}, {
  selector: 'edge[direction="dest_to_src"]', // 源指向目标
  style: {
    'target-arrow-shape': 'triangle-backcurve',
    'target-arrow-color': '#c3cdd7',
    'source-arrow-shape': 'none'
  }
}, {
  selector: '.weaken',
  style: {
    opacity: 0.6
  }
}, {
  selector: 'node.highlight',
  style: {
    opacity: 1,
    'border-color': HIGHLIGHT_COLOR
  }
}, {
  selector: 'edge.highlight',
  style: {
    opacity: 1,
    'line-color': HIGHLIGHT_COLOR,
    'source-arrow-color': HIGHLIGHT_COLOR,
    'target-arrow-color': HIGHLIGHT_COLOR
  }
}, {
  selector: 'node.hover',
  style: {
    'border-width': 1.5,
    'border-color': HOVER_COLOR,
    'font-weight': 'bold',
    'background-image': 'data(background.hover)'
  }
}, {
  selector: 'edge.hover',
  style: {
    width: 3,
    label: 'data(label)',
    'line-color': HOVER_COLOR,
    'source-arrow-color': HOVER_COLOR,
    'target-arrow-color': HOVER_COLOR,
    'font-weight': 'bold'
  }
}, {
  selector: 'node.current',
  style: {
    width: 56,
    height: 56,
    'background-image': 'data(background.selected)',
    'border-color': HIGHLIGHT_COLOR,
    'font-weight': 'bold'
  }
}]

export default {
  layout,
  style
}
