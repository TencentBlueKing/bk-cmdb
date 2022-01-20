export default {
  namespaced: true,
  state: {
    /**
     * 业务集 ID
     * @type {number}
     */
    bizSetId: null,
    /**
     * 业务集名称
     * @type {string}
     */
    bizSetName: null,
    /**
     * 业务集列表
     * @type {Array}
     */
    bizSetList: [],
    /**
     * 当前所选节点的业务 ID
     * @type {number}
     */
    bizId: null,
  },
  mutations: {
    setBizSetId(state, bizSetId) {
      state.bizSetId = Number(bizSetId)
      state.bizSetName = state.bizSetList.find(item => item.bk_biz_set_id === state.bizSetId)?.bk_biz_set_name
    },
    setBizSetList(state, bizSetList) {
      state.bizSetList = bizSetList
    },
    setBizId(state, bizId) {
      state.bizId = Number(bizId)
    }
  }
}
