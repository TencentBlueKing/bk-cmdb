/*
 * Tencent is pleased to support the open source community by making
 * 蓝鲸智云 - 配置平台 (BlueKing - Configuration System) available.
 * Copyright (C) 2017 THL A29 Limited,
 * a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 * We undertake not to change the open source license (MIT license) applicable
 * to the current version of the project delivered to anyone in the future.
 */

package transfer

import (
	"fmt"
	"regexp"
	"strconv"
)

// EncodeDot encode the dot in the string as the unicode value
func EncodeDot(input string) string {
	re := regexp.MustCompile(`\.`)
	encodedStr := re.ReplaceAllStringFunc(input, func(s string) string {
		return fmt.Sprintf("\\u%04x", rune('.'))
	})
	return encodedStr
}

// DecodeDot decode the unicode value of dot in a string to dot
func DecodeDot(input string) string {
	re := regexp.MustCompile(`\\u([0-9a-fA-F]{4})`)
	decodedStr := re.ReplaceAllStringFunc(input, func(s string) string {
		unicodePoint, _ := strconv.ParseInt(s[2:], 16, 32)
		return string(rune(unicodePoint))
	})
	return decodedStr
}
