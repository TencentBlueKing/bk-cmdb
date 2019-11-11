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

package upgrader

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"configcenter/src/common/blog"
)

/*
	migration版本号对比
*/

const VersionNgPrefix = "y"

func VersionCmp(version1, version2 string) int {
	if strings.HasPrefix(version1, VersionNgPrefix) && strings.HasPrefix(version2, VersionNgPrefix) {
		return V36VersionCmp(version1, version2)
	}
	if strings.HasPrefix(version1, VersionNgPrefix) || strings.HasPrefix(version2, VersionNgPrefix) {
		// y > x > v
		return StringCompare(version1, version2)
	}
	// legacy compare scheme
	return StringCompare(version1, version2)
}

func Int64Compare(val1, val2 int64) int {
	if val1 == val2 {
		return 0
	} else if val1 > val2 {
		return 1
	} else {
		return -1
	}
}

func StringCompare(s1, s2 string) int {
	if s1 == s2 {
		return 0
	} else if s1 > s2 {
		return 1
	} else {
		return -1
	}
}

func V36VersionCmp(version1, version2 string) int {
	// version format should be validate before compare
	ngVersion1, err := ParseNgVersion(version1)
	if err != nil {
		blog.Fatalf(err.Error())
	}
	ngVersion2, err := ParseNgVersion(version2)
	if err != nil {
		blog.Fatalf(err.Error())
	}
	result := Int64Compare(ngVersion1.Major, ngVersion2.Major)
	if result != 0 {
		return result
	}
	result = Int64Compare(ngVersion1.Minor, ngVersion2.Minor)
	if result != 0 {
		return result
	}
	return StringCompare(ngVersion1.Patch, ngVersion2.Patch)
}

type NgVersion struct {
	Major int64
	Minor int64
	Patch string
}

var PatchRegex = regexp.MustCompile(`^\d{12}$`)

func ParseNgVersion(version string) (NgVersion, error) {
	ngVersion := NgVersion{}
	invalidMessage := fmt.Errorf("invalid version [%s]", version)
	version = strings.TrimLeft(version, VersionNgPrefix)
	fields := strings.Split(version, ".")
	if len(fields) != 3 {
		return ngVersion, invalidMessage
	}

	major, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return ngVersion, invalidMessage
	}
	ngVersion.Major = major

	minor, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return ngVersion, invalidMessage
	}
	ngVersion.Minor = minor

	patch := fields[2]
	match := PatchRegex.MatchString(patch)
	if match == false {
		return ngVersion, invalidMessage
	}
	ngVersion.Patch = patch
	return ngVersion, nil
}
