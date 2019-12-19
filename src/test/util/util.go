/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017,-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package util

import (
	"bytes"
	"math/rand"
	"strconv"
	"text/template"

	. "github.com/onsi/gomega"
	"k8s.io/client-go/util/jsonpath"
)

func JsonPathExtractInt(responses map[string]interface{}, key, name, path string) int64 {
	output := JsonPathExtract(responses, key, name, path)
	value, err := strconv.ParseInt(output, 10, 64)
	Expect(err).NotTo(HaveOccurred())
	return value
}

func JsonPathExtract(responses map[string]interface{}, key, name, path string) string {
	j := jsonpath.New(name)
	err := j.Parse(path)
	Expect(err).NotTo(HaveOccurred())
	Expect(responses).Should(HaveKey(key))

	buf := new(bytes.Buffer)
	err = j.Execute(buf, responses[key])
	Expect(err).NotTo(HaveOccurred(), key)
	out := buf.String()

	return out
}

func RenderTemplate(templateContext map[string]interface{}, urlTemplate string) string {
	var subPathBuffer bytes.Buffer
	t := template.Must(template.New("url template").Parse(urlTemplate))
	err := t.Execute(&subPathBuffer, templateContext)
	Expect(err).NotTo(HaveOccurred())
	subPath := subPathBuffer.String()
	return subPath

}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
