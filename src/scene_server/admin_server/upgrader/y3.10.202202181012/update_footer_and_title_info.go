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

package y3_10_202202181012

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/blog"
	"configcenter/src/common/metadata"
	"configcenter/src/common/version"
	"configcenter/src/scene_server/admin_server/upgrader"
	"configcenter/src/storage/dal"
)

const (
	// NOCC:tosa/linelength(忽略长度)
	oldFooterCn = "[QQ咨询(800802001)](http://wpa.b.qq.com/cgi/wpa.php?ln=1&key=XzgwMDgwMjAwMV80NDMwOTZfODAwODAyMDAxXzJf)|[蓝鲸论坛](https://bk.tencent.com/s-mart/community/)|[蓝鲸官网](https://bk.tencent.com/index/)|[蓝鲸桌面](http://your-bk-desktop.com)"
	// NOCC:tosa/linelength(忽略长度)
	oldFooterEn = "[QQ(800802001)](http://wpa.b.qq.com/cgi/wpa.php?ln=1&key=XzgwMDgwMjAwMV80NDMwOTZfODAwODAyMDAxXzJf)|[Blueking Forum](https://bk.tencent.com/s-mart/community/)|[BlueKing Official](https://bk.tencent.com/index/)|[Blueking Desktop](http://your-bk-desktop.com)"
	// NOCC:tosa/linelength(忽略长度)
	footerCn = "[技术支持](https://wpa1.qq.com/KziXGWJs?_type=wpa&qidian=true) | [社区论坛](https://bk.tencent.com/s-mart/community/) | [产品官网](https://bk.tencent.com/index/)"
	// NOCC:tosa/linelength(忽略长度)
	footerEn          = "[Support](https://wpa1.qq.com/KziXGWJs?_type=wpa&amp;qidian=true) | [Forum](https://bk.tencent.com/s-mart/community/) | [Official](https://bk.tencent.com/index/)"
	communityCn       = "配置平台 | 蓝鲸智云社区版"
	tencentBlueKingCn = "配置平台 | 腾讯蓝鲸智云"
	enterpriseCn      = "配置平台 | 蓝鲸智云企业版"
	versionEn         = "CMDB | BlueKing"
	tencentVersion    = "CMDB | Tencent BlueKing"
)

func getDBPlatformSetting(ctx context.Context, db dal.RDB) (*metadata.PlatformSettingConfig, error) {
	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}

	ret := make(map[string]interface{})
	err := db.Table(common.BKTableNameSystem).Find(cond).Fields(common.ConfigAdminValueField).One(ctx, &ret)
	if err != nil {
		blog.Errorf("get platform setting failed, err: %v", err)
		return nil, err
	}

	value := ret[common.ConfigAdminValueField]
	if value == nil {
		blog.Errorf("failed to get platform management config field, err: %v", err)
		return nil, errors.New("get platform config failed")
	}
	setting, ok := value.(string)
	if !ok {
		blog.Errorf("platform management config field type validation error, type: %v", reflect.TypeOf(value))
		return nil, errors.New("platform config type error")
	}

	platformSetting := new(metadata.PlatformSettingConfig)
	if err := json.Unmarshal([]byte(setting), platformSetting); err != nil {
		blog.Errorf("platform setting unmarshal failed, config: %v, err: %v", setting, err)
		return nil, fmt.Errorf("config unmarshal failed, err: %v", err)
	}
	return platformSetting, nil
}

// getFinalConfig 1、在更改默认title配置时需要根据不同版本的编译参数判断系统的title此时是否是原初始化配置，只有是初始化配置时下才可以对title
// 的默认文案进行更新，如果用户已经改过则保留用户改后的配置。
// 2、只有当初始化配置时才将联系方式改成新的默认值，如果footer中的值已经被用户改过，则需要保留用户配置。
func getFinalConfig(dbConfig *metadata.PlatformSettingConfig) {

	switch version.CCDistro {

	case version.CCDistrEnterprise:
		if dbConfig.SiteConfig.SiteName.I18N.CN == enterpriseCn {
			dbConfig.SiteConfig.SiteName.Value = tencentBlueKingCn
			dbConfig.SiteConfig.SiteName.I18N.CN = tencentBlueKingCn
		}

	default:
		if dbConfig.SiteConfig.SiteName.I18N.CN == communityCn {
			dbConfig.SiteConfig.SiteName.Value = tencentBlueKingCn
			dbConfig.SiteConfig.SiteName.I18N.CN = tencentBlueKingCn
		}
	}
	if dbConfig.SiteConfig.SiteName.I18N.EN == versionEn {
		dbConfig.SiteConfig.SiteName.I18N.EN = tencentVersion
	}

	if dbConfig.FooterConfig.ContactInfo.I18N.CN == oldFooterCn {
		dbConfig.FooterConfig.ContactInfo.Value = footerCn
		dbConfig.FooterConfig.ContactInfo.I18N.CN = footerCn
	}
	if dbConfig.FooterConfig.ContactInfo.I18N.EN == oldFooterEn {
		dbConfig.FooterConfig.ContactInfo.I18N.EN = footerEn
	}
	return
}

func updateTitleAndFooterInfo(ctx context.Context, db dal.RDB, conf *upgrader.Config) error {

	dbConfig, err := getDBPlatformSetting(ctx, db)
	if err != nil {
		return err
	}

	getFinalConfig(dbConfig)

	if err := updatePlatformConfig(ctx, db, dbConfig); err != nil {
		blog.Errorf("upgrade platform config encode base64 failed, cfg: %#v, err: %v", *dbConfig, err)
		return err
	}
	return nil
}

// updatePlatformConfig update configuration to database.
func updatePlatformConfig(ctx context.Context, db dal.RDB, config *metadata.PlatformSettingConfig) error {
	bytes, err := json.Marshal(config)
	if err != nil {
		return err
	}

	cond := map[string]interface{}{
		"_id": common.ConfigAdminID,
	}
	data := map[string]interface{}{
		common.ConfigAdminValueField: string(bytes),
		common.LastTimeField:         time.Now(),
	}

	err = db.Table(common.BKTableNameSystem).Update(ctx, cond, data)
	if err != nil {
		return err
	}
	return nil
}
