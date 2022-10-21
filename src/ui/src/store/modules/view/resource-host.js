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

const state = {
  activeDirectory: null,
  directoryList: []
}

const getters = {
  activeDirectory: state => state.activeDirectory,
  directoryList: state => state.directoryList,
  defaultDirectory: state => state.directoryList.find(directory => directory.default === 1),
  directorySortedList: state => (topList, directoryList) => {
    const list = directoryList || [...state.directoryList]
    const count = topList.length
    list.sort((dirA, dirB) => {
      const stickyIndexA = topList.indexOf(dirA.bk_module_id) + 1
      const stickyIndexB = topList.indexOf(dirB.bk_module_id) + 1

      return (stickyIndexA || (count + 1)) - (stickyIndexB || (count + 1))
    })
    return list
  }
}

const mutations = {
  setActiveDirectory(state, active) {
    state.activeDirectory = active
  },
  setDirectoryList(state, list) {
    state.directoryList = list
  },
  addDirectory(state, directory) {
    state.directoryList.splice(1, 0, directory)
  },
  updateDirectory(state, directory) {
    const index = state.directoryList.findIndex(data => data.bk_module_id === directory.bk_module_id)
    if (index > -1) {
      state.directoryList.splice(index, 1, directory)
    }
  },
  deleteDirectory(state, id) {
    const index = state.directoryList.findIndex(target => target.bk_module_id === id)
    if (index > -1) {
      state.directoryList.splice(index, 1)
    }
  },
  refreshDirectoryCount(state, newList = []) {
    state.directoryList.forEach((directory) => {
      const newDirectory = newList.find(newDirectory => newDirectory.bk_module_id === directory.bk_module_id)
      Object.assign(directory, newDirectory)
    })
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations
}
