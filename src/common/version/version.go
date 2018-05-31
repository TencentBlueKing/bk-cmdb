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

package version

import (
	"fmt"
)

//CCVersion discribes version
//CCTag show the git tag for this version
//CCBuildTime show the compile time
var (
	CCVersion       = "17.03.28"
	CCTag           = "2017-03-28 Release"
	CCBuildTime     = "2017-03-28 19:50:00"
	CCGitHash       = "unknown"
	CCRunMode       = "product" // product, test, dev
	CCDistro        = "enterprise"
	CCDistroVersion = "9999.9999.9999"
)

// CCRunMode enumeration
var (
	CCRunModeProduct = "product"
	CCRunModeTest    = "test"
	CCRunModeDev     = "dev"
)

//ShowVersion is the default handler which match the --version flag
func ShowVersion() {
	fmt.Printf("%s", GetVersion())
}

// GetVersion return the version info
func GetVersion() string {
	version := fmt.Sprintf("Version  :%s\nTag      :%s\nBuildTime:  %s\nGitHash:    %s\nRunMode:    %s\n", CCVersion, CCTag, CCBuildTime, CCGitHash, CCRunMode)
	return version
}
