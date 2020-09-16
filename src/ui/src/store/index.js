import Vue from 'vue'
import Vuex from 'vuex'
import global from './modules/global.js'
import request from './modules/request.js'
import index from './modules/view/index.js'
import hostDetails from './modules/view/host-details.js'
import serviceProcess from './modules/view/service-process.js'
import businessSync from './modules/view/business-sync.js'
import hosts from './modules/view/hosts.js'
import setFeatures from './modules/view/set-features.js'
import auth from './modules/api/auth.js'
import eventSub from './modules/api/event-sub.js'
import hostBatch from './modules/api/host-batch.js'
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
import procConfig from './modules/api/proc-config.js'
import userCustom from './modules/api/user-custom.js'
import userPrivilege from './modules/api/user-privilege.js'
import globalModels from './modules/api/global-models.js'
import cloudDiscover from './modules/api/cloud-discover'
import netCollectDevice from './modules/api/net-collect-device.js'
import netCollectProperty from './modules/api/net-collect-property.js'
import netDataCollection from './modules/api/net-data-collection.js'
import netDiscovery from './modules/api/net-discovery.js'
import serviceTemplate from './modules/api/service-template.js'
import serviceClassification from './modules/api/service-classification.js'
import processTemplate from './modules/api/process-template.js'
import businessSynchronous from './modules/api/business-synchronous.js'
import serviceInstance from './modules/api/service-instance.js'
import processInstance from './modules/api/process-instance.js'
import operationChart from './modules/api/operation-chart'
import instanceLabel from './modules/api/instance-label.js'
import fullTextSearch from './modules/api/full-text-search.js'
import setSync from './modules/api/set-sync.js'
import setTemplate from './modules/api/set-template.js'
import cloud from './modules/api/cloud.js'
import hostApply from './modules/api/host-apply'

import resourceDirectory from './modules/api/resource-directory.js'
import resource from './modules/api/resource.js'

import organization from './modules/api/organization'

import businessHost from './modules/view/business-host.js'
import resourceHost from './modules/view/resource-host.js'
import cloudarea from './modules/api/cloudarea'
import audit from './modules/api/audit.js'
import dynamicGroup from './modules/api/dynamic-group'

Vue.use(Vuex)

export default new Vuex.Store({
    ...global,
    modules: {
        index,
        hostDetails,
        serviceProcess,
        businessSync,
        hosts,
        setFeatures,
        auth,
        request,
        eventSub,
        hostBatch,
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
        procConfig,
        userCustom,
        userPrivilege,
        globalModels,
        cloudDiscover,
        netCollectDevice,
        netCollectProperty,
        netDataCollection,
        netDiscovery,
        serviceTemplate,
        serviceClassification,
        processTemplate,
        businessSynchronous,
        serviceInstance,
        processInstance,
        operationChart,
        instanceLabel,
        fullTextSearch,
        setSync,
        setTemplate,
        businessHost,
        cloud,
        hostApply,
        resourceHost,
        resourceDirectory,
        resource,
        cloudarea,
        organization,
        audit,
        dynamicGroup
    }
})
