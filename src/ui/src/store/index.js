import Vue from 'vue'
import Vuex from 'vuex'

import global from './modules/global.js'
import request from './modules/request.js'

import eventSub from './modules/api/event-sub.js'
import hostCustomApi from './modules/api/host-custom-api.js'
import hostDelete from './modules/api/host-delete.js'
import hostFavorites from './modules/api/host-favorites.js'
import hostRelation from './modules/api/host-relation.js'
import hostSearchHistory from './modules/api/host-search-history.js'
import hostSearch from './modules/api/host-search.js'
import hostUpdate from './modules/api/host-update.js'
import objectAssociation from './modules/api/object-association.js'
import objectBatch from './modules/api/object-batch.js'
import objectBiz from './modules/api/object-biz.js'
import objectCommonInst from './modules/api/object-common-inst.js'
import objectMainLineModule from './modules/api/object-main-line-module.js'
import objectModelClassify from './modules/api/object-model-classify.js'
import objectModelFieldGroup from './modules/api/object-model-field-group.js'
import objectModelProperty from './modules/api/object-model-property.js'
import objectModel from './modules/api/object-model.js'
import objectModule from './modules/api/object-module.js'
import objectRelation from './modules/api/object-relation.js'
import objectSet from './modules/api/object-set.js'
import objectUnique from './modules/api/object-unique.js'
import operationAudit from './modules/api/operation-audit.js'
import procConfig from './modules/api/proc-config.js'
import userCustom from './modules/api/user-custom.js'
import userPrivilege from './modules/api/user-privilege.js'
import globalModels from './modules/api/global-models.js'
import cloudDiscover from './modules/api/cloud-discover'
import netCollectDevice from './modules/api/net-collect-device.js'
import netCollectProperty from './modules/api/net-collect-property.js'
import netDataCollection from './modules/api/net-data-collection.js'
import netDiscovery from './modules/api/net-discovery.js'

Vue.use(Vuex)

export default new Vuex.Store({
    ...global,
    modules: {
        request,
        eventSub,
        hostCustomApi,
        hostDelete,
        hostFavorites,
        hostRelation,
        hostSearchHistory,
        hostSearch,
        hostUpdate,
        objectAssociation,
        objectBatch,
        objectBiz,
        objectCommonInst,
        objectMainLineModule,
        objectModelClassify,
        objectModelFieldGroup,
        objectModelProperty,
        objectModel,
        objectModule,
        objectRelation,
        objectSet,
        objectUnique,
        operationAudit,
        procConfig,
        userCustom,
        userPrivilege,
        globalModels,
        cloudDiscover,
        netCollectDevice,
        netCollectProperty,
        netDataCollection,
        netDiscovery
    }
})
