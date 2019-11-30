/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package authsynchronizer

import (
	"context"
	"time"

	"configcenter/src/apimachinery"
	"configcenter/src/auth/extensions"
	"configcenter/src/common"
	"configcenter/src/common/backbone"
	"configcenter/src/common/blog"
	"configcenter/src/common/mapstr"
	"configcenter/src/common/metadata"
	"configcenter/src/scene_server/admin_server/authsynchronizer/meta"
	"configcenter/src/scene_server/admin_server/authsynchronizer/utils"
)

// Producer producer WorkRequest and enqueue it
type Producer struct {
	clientSet           apimachinery.ClientSetInterface
	authManager         *extensions.AuthManager
	ID                  int
	WorkerQueue         chan meta.WorkRequest
	QuitChan            chan bool
	Engine              *backbone.Engine
	SyncIntervalMinutes int
}

// NewProducer make a producer
func NewProducer(clientSet apimachinery.ClientSetInterface, authManager *extensions.AuthManager,
	workerQueue chan meta.WorkRequest, engine *backbone.Engine, syncIntervalMinutes int) *Producer {
	// Create, and return the producer.
	producer := Producer{
		clientSet:           clientSet,
		authManager:         authManager,
		ID:                  0,
		WorkerQueue:         workerQueue,
		QuitChan:            make(chan bool),
		Engine:              engine,
		SyncIntervalMinutes: syncIntervalMinutes,
	}

	return &producer
}

// Start do main loop
func (p *Producer) Start() {
	// then tick and loop
	if p.SyncIntervalMinutes < meta.MinSyncIntervalMinutes {
		blog.Warnf("SyncIntervalMinutes min value is: %d, config is: %d", meta.MinSyncIntervalMinutes, p.SyncIntervalMinutes)
		p.SyncIntervalMinutes = meta.MinSyncIntervalMinutes
	}
	blog.Infof("start producer with SyncIntervalMinutes value: %d", p.SyncIntervalMinutes)
	duration := time.Duration(p.SyncIntervalMinutes) * time.Minute
	ticker := time.NewTicker(duration)
	go func(producer *Producer) {
		// loop immediately at first.
		time.Sleep(1 * time.Minute)
		p.loop()

		for {
			select {
			case <-ticker.C:
				if len(p.WorkerQueue) > 0 {
					blog.Infof("workerQueue not empty, skip generate new job, current length: %d", len(p.WorkerQueue))
					continue
				}
				p.loop()
			case <-p.QuitChan:
				ticker.Stop()
				return
			}
		}
	}(p)
}

func (p *Producer) loop() {
	if isMaster := p.Engine.ServiceManageInterface.IsMaster(); !isMaster {
		blog.Info("not master, don't generate iam sync job")
		return
	}

	// get jobs
	jobs := p.generateJobs()

	blog.Infof("generate auth synchronize jobs, count: %d", len(*jobs))

	for _, job := range *jobs {
		p.WorkerQueue <- job
	}
}

func (p *Producer) generateJobs() *[]meta.WorkRequest {
	// split all jobs
	jobs := make([]meta.WorkRequest, 0)

	businessSyncJob := meta.WorkRequest{
		ResourceType: meta.BusinessResource,
		Data:         map[string]interface{}{},
	}
	jobs = append(jobs, businessSyncJob)

	// list all business
	header := utils.NewListBusinessAPIHeader()
	// find business, but not contains resource pool business
	condition := metadata.QueryCondition{Condition: mapstr.MapStr{"default": 0}}
	result, err := p.clientSet.CoreService().Instance().ReadInstance(context.TODO(), *header, common.BKInnerObjIDApp, &condition)
	if err != nil {
		blog.Errorf("list business failed, err: %v", err)
		return &jobs
	}

	businessList := make([]extensions.BusinessSimplify, 0)
	for _, business := range result.Data.Info {
		businessSimplify := extensions.BusinessSimplify{}
		if _, err := businessSimplify.Parse(business); err != nil {
			blog.Errorf("parse business %+v simplify information failed, err: %+v", business, err)
			continue
		}
		businessList = append(businessList, businessSimplify)
	}
	blog.V(4).Infof("list business, count:%d", len(businessList))

	// job of synchronize business scope resources to iam
	resourceTypes := []meta.ResourceType{
		meta.HostBizResource,
		meta.SetResource,
		meta.ModuleResource,
		meta.ModelResource,
		meta.ProcessResource,
		meta.DynamicGroupResource,
		// meta.AuditCategory,
		meta.ClassificationResource,
		// meta.UserGroupSyncResource,
		meta.ServiceTemplateResource,
		meta.SetTemplateResource,
	}
	for _, resourceType := range resourceTypes {
		for _, businessSimplify := range businessList {
			jobs = append(jobs, meta.WorkRequest{
				ResourceType: resourceType,
				Data:         businessSimplify,
			})
		}
	}

	globalBusiness := extensions.BusinessSimplify{
		BKAppIDField:      0,
		BKAppNameField:    "",
		BKSupplierIDField: 0,
		BKOwnerIDField:    "0",
		IsDefault:         0,
	}
	instanceBizList := append(businessList, globalBusiness)
	for _, business := range instanceBizList {
		header := utils.NewListBusinessAPIHeader()
		objects, err := p.authManager.CollectObjectsByBusinessID(context.Background(), *header, business.BKAppIDField)
		if err != nil {
			blog.Errorf("get models by business id: %d failed, err: %+v", business.BKAppIDField, err)
			continue
		}
		for _, object := range objects {
			header := utils.NewListBusinessAPIHeader()
			jobs = append(jobs, meta.WorkRequest{
				ResourceType: meta.InstanceResource,
				Data:         object,
				Header:       *header,
			})
		}
	}

	// some resource type need sync global resource
	resourceTypes = []meta.ResourceType{
		// meta.AuditCategory,
		meta.ClassificationResource,
		meta.PlatResource,
		meta.ModelResource,
		meta.HostResourcePool,
	}
	for _, resourceType := range resourceTypes {
		header := utils.NewListBusinessAPIHeader()
		jobs = append(jobs, meta.WorkRequest{
			ResourceType: resourceType,
			Data:         globalBusiness,
			Header:       *header,
		})
	}

	jobs = FilterJobs(jobs)

	blog.Infof("jobs: count: %d", len(jobs))
	return &jobs
}

func FilterJobs(jobs []meta.WorkRequest) []meta.WorkRequest {
	debugSync := false
	if debugSync == false {
		return jobs
	}
	debugResourceType := make([]meta.ResourceType, 0)
	finalJobs := make([]meta.WorkRequest, 0)
	for _, job := range jobs {
		found := false
		for _, resourceType := range debugResourceType {
			if resourceType == job.ResourceType {
				found = true
				break
			}
		}
		if found {
			finalJobs = append(finalJobs, job)
		}
	}

	return finalJobs
}
