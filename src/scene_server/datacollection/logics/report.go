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

package logics

import (
	"fmt"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/condition"
	"configcenter/src/common/metadata"
)

func (lgc *Logics) SearchReportSummary(param metadata.ParamSearchNetcollectReport) ([]*metadata.NetcollectReportSummary, error) {
	// search reports
	param.CloudID = -1
	_, reports, err := lgc.SearchReport(param)
	if err != nil {
		blog.Errorf("[SearchReportSummary] SearchReport failed %v", err)
		return nil, err
	}
	summarym := map[int64]*metadata.NetcollectReportSummary{}

	// combine details
	for _, report := range reports {
		if summary, ok := summarym[report.CloudID]; !ok {
			summarym[report.CloudID] = &metadata.NetcollectReportSummary{
				CloudID:    report.CloudID,
				CloudName:  report.CloudName,
				Statistics: map[string]int{},
			}
		} else {
			summary.Statistics["associations"] += len(report.Detail.Associations)
			for _, association := range report.Detail.Associations {
				if association.LastTime.Time.Sub(summary.LastTime.Time) > 0 {
					summary.LastTime = association.LastTime
				}
			}
			for _, attribute := range report.Detail.Attributes {
				summary.Statistics[attribute.ObjectName]++
				if attribute.LastTime.Time.Sub(summary.LastTime.Time) > 0 {
					summary.LastTime = attribute.LastTime
				}
			}
		}
	}

	summarys := []*metadata.NetcollectReportSummary{}
	for key := range summarym {
		summarys = append(summarys, summarym[key])
	}

	return summarys, nil
}

func (lgc *Logics) SearchReport(param metadata.ParamSearchNetcollectReport) (int64, []*metadata.NetcollectReport, error) {
	cond := condition.CreateCondition()
	if param.CloudID >= 0 {
		cond.Field(common.BKCloudIDField).Eq(param.CloudID)
	}
	if param.Action != "" {
		cond.Field("action").Eq(param.Action)
	}
	if param.ObjectID != "" {
		cond.Field(common.BKObjIDField).Eq(param.ObjectID)
	}
	if param.InnerIP != "" {
		cond.Field(common.BKHostInnerIPField).Like(param.InnerIP)
	}

	count, err := lgc.Instance.GetCntByCondition(common.BKTableNameNetcollectReport, cond)
	if err != nil {
		blog.Errorf("[SearchReport] GetMutilByCondition failed %v", err)
		return 0, nil, err
	}

	// search reports
	reportm := map[string]*metadata.NetcollectReport{}
	reports := []*metadata.NetcollectReport{}
	err = lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectReport, nil, cond, &reports, param.Page.Sort, param.Page.Start, param.Page.Limit)
	if err != nil {
		blog.Errorf("[SearchReport] GetMutilByCondition failed %v", err)
		return 0, nil, err
	}

	// search details
	ips := []string{}
	for _, report := range reports {
		ips = append(ips, report.InnerIP)
		reportm[fmt.Sprintf("%v:%v", report.CloudID, report.InstName)] = report
	}
	attributes, associations, err := lgc.searchDetail(param.CloudID, ips)
	if err != nil {
		return 0, nil, err
	}

	// combine details
	for _, attribute := range attributes {
		reportm[fmt.Sprintf("%v:%v", attribute.CloudID, attribute.InstName)].Detail.Attributes =
			append(reportm[fmt.Sprintf("%v:%v", attribute.CloudID, attribute.InstName)].Detail.Attributes,
				attribute)
	}
	for _, association := range associations {
		reportm[fmt.Sprintf("%v:%v", association.CloudID, association.InstName)].Detail.Associations =
			append(reportm[fmt.Sprintf("%v:%v", association.CloudID, association.InstName)].Detail.Associations,
				association)
	}

	return int64(count), reports, nil
}

func (lgc *Logics) searchDetail(cloudID int64, ips []string) ([]metadata.NetcollectReportAttribute, []metadata.NetcollectReportAssociation, error) {
	detailCond := condition.CreateCondition()
	detailCond.Field(common.BKCloudIDField).Eq(cloudID)
	detailCond.Field(common.BKHostInnerIPField).In(ips)

	attributes := []metadata.NetcollectReportAttribute{}
	associations := []metadata.NetcollectReportAssociation{}
	err := lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectAssociation, nil, detailCond, &attributes, "", 0, 0)
	if err != nil {
		blog.Errorf("[SearchReport] GetMutilByCondition failed %v", err)
		return nil, nil, err
	}
	err = lgc.Instance.GetMutilByCondition(common.BKTableNameNetcollectAssociation, nil, detailCond, &associations, "", 0, 0)
	if err != nil {
		blog.Errorf("[SearchReport] GetMutilByCondition failed %v", err)
		return nil, nil, err
	}

	return attributes, associations, nil
}

func (lgc *Logics) ComfirmReport(metadata.RspNetcollectReport) {

}
