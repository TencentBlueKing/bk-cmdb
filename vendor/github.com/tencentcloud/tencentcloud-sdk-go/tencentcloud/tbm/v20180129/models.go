// Copyright (c) 2017-2018 THL A29 Limited, a Tencent company. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v20180129

import (
    "encoding/json"

    tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

type AgePortrait struct {

	// 年龄区间
	AgeRange *string `json:"AgeRange" name:"AgeRange"`

	// 百分比
	Percent *float64 `json:"Percent" name:"Percent"`
}

type AgePortraitInfo struct {

	// 用户年龄画像数组
	PortraitSet []*AgePortrait `json:"PortraitSet" name:"PortraitSet" list`
}

type BrandReportArticle struct {

	// 文章标题
	Title *string `json:"Title" name:"Title"`

	// 文章url地址
	Url *string `json:"Url" name:"Url"`

	// 文章来源
	FromSite *string `json:"FromSite" name:"FromSite"`

	// 文章发表日期
	PubTime *string `json:"PubTime" name:"PubTime"`

	// 文章标识
	Flag *uint64 `json:"Flag" name:"Flag"`

	// 文章热度值
	Hot *uint64 `json:"Hot" name:"Hot"`

	// 文章来源等级
	Level *uint64 `json:"Level" name:"Level"`

	// 文章摘要
	Abstract *string `json:"Abstract" name:"Abstract"`

	// 文章ID
	ArticleId *string `json:"ArticleId" name:"ArticleId"`
}

type Comment struct {

	// 评论的日期
	Date *string `json:"Date" name:"Date"`

	// 差评的个数
	NegCommentCount *uint64 `json:"NegCommentCount" name:"NegCommentCount"`

	// 好评的个数
	PosCommentCount *uint64 `json:"PosCommentCount" name:"PosCommentCount"`
}

type CommentInfo struct {

	// 用户评论内容
	Comment *string `json:"Comment" name:"Comment"`

	// 评论的时间
	Date *string `json:"Date" name:"Date"`
}

type DateCount struct {

	// 统计日期
	Date *string `json:"Date" name:"Date"`

	// 统计值
	Count *uint64 `json:"Count" name:"Count"`
}

type DescribeBrandCommentCountRequest struct {
	*tchttp.BaseRequest

	// 品牌ID
	BrandId *string `json:"BrandId" name:"BrandId"`

	// 查询开始日期
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 查询结束日期
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeBrandCommentCountRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandCommentCountRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandCommentCountResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 按天统计好评/差评数
		CommentSet []*Comment `json:"CommentSet" name:"CommentSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBrandCommentCountResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandCommentCountResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandExposureRequest struct {
	*tchttp.BaseRequest

	// 品牌ID
	BrandId *string `json:"BrandId" name:"BrandId"`

	// 查询开始时间
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 查询结束时间
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeBrandExposureRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandExposureRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandExposureResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 累计曝光量
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 按天计算的统计数据
		DateCountSet []*DateCount `json:"DateCountSet" name:"DateCountSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBrandExposureResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandExposureResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandMediaReportRequest struct {
	*tchttp.BaseRequest

	// 品牌ID
	BrandId *string `json:"BrandId" name:"BrandId"`

	// 查询开始时间
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 查询结束时间
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeBrandMediaReportRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandMediaReportRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandMediaReportResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 查询范围内文章总数
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 按天计算的每天文章数
		DateCountSet []*DateCount `json:"DateCountSet" name:"DateCountSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBrandMediaReportResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandMediaReportResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandNegCommentsRequest struct {
	*tchttp.BaseRequest

	// 品牌ID
	BrandId *string `json:"BrandId" name:"BrandId"`

	// 查询开始时间
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 查询结束时间
	EndDate *string `json:"EndDate" name:"EndDate"`

	// 查询条数上限，默认20
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 查询偏移，默认从0开始
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeBrandNegCommentsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandNegCommentsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandNegCommentsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 评论列表
		BrandCommentSet []*CommentInfo `json:"BrandCommentSet" name:"BrandCommentSet" list`

		// 总的差评个数
		TotalComments *uint64 `json:"TotalComments" name:"TotalComments"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBrandNegCommentsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandNegCommentsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandPosCommentsRequest struct {
	*tchttp.BaseRequest

	// 品牌ID
	BrandId *string `json:"BrandId" name:"BrandId"`

	// 查询开始时间
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 查询结束时间
	EndDate *string `json:"EndDate" name:"EndDate"`

	// 查询条数上限，默认20
	Limit *uint64 `json:"Limit" name:"Limit"`

	// 查询偏移，从0开始
	Offset *uint64 `json:"Offset" name:"Offset"`
}

func (r *DescribeBrandPosCommentsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandPosCommentsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandPosCommentsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 评论列表
		BrandCommentSet []*CommentInfo `json:"BrandCommentSet" name:"BrandCommentSet" list`

		// 总的好评个数
		TotalComments *uint64 `json:"TotalComments" name:"TotalComments"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBrandPosCommentsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandPosCommentsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandSocialOpinionRequest struct {
	*tchttp.BaseRequest

	// 品牌ID
	BrandId *string `json:"BrandId" name:"BrandId"`

	// 检索开始时间
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 检索结束时间
	EndDate *string `json:"EndDate" name:"EndDate"`

	// 查询偏移，默认从0开始
	Offset *int64 `json:"Offset" name:"Offset"`

	// 查询条数上限，默认20
	Limit *int64 `json:"Limit" name:"Limit"`

	// 列表显示标记，若为true，则返回文章列表详情
	ShowList *bool `json:"ShowList" name:"ShowList"`
}

func (r *DescribeBrandSocialOpinionRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandSocialOpinionRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandSocialOpinionResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 文章总数
		ArticleCount *uint64 `json:"ArticleCount" name:"ArticleCount"`

		// 来源统计总数
		FromCount *uint64 `json:"FromCount" name:"FromCount"`

		// 疑似负面报道总数
		AdverseCount *uint64 `json:"AdverseCount" name:"AdverseCount"`

		// 文章列表详情
		ArticleSet []*BrandReportArticle `json:"ArticleSet" name:"ArticleSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBrandSocialOpinionResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandSocialOpinionResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandSocialReportRequest struct {
	*tchttp.BaseRequest

	// 品牌ID
	BrandId *string `json:"BrandId" name:"BrandId"`

	// 查询开始时间
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 查询结束时间
	EndDate *string `json:"EndDate" name:"EndDate"`
}

func (r *DescribeBrandSocialReportRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandSocialReportRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeBrandSocialReportResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 累计统计数据
		TotalCount *uint64 `json:"TotalCount" name:"TotalCount"`

		// 按天计算的统计数据
		DateCountSet []*DateCount `json:"DateCountSet" name:"DateCountSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeBrandSocialReportResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeBrandSocialReportResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeIndustryNewsRequest struct {
	*tchttp.BaseRequest

	// 行业ID
	IndustryId *string `json:"IndustryId" name:"IndustryId"`

	// 查询开始时间
	StartDate *string `json:"StartDate" name:"StartDate"`

	// 查询结束时间
	EndDate *string `json:"EndDate" name:"EndDate"`

	// 是否显示列表，若为 true，则返回文章列表
	ShowList *bool `json:"ShowList" name:"ShowList"`

	// 查询偏移，默认从0开始
	Offset *int64 `json:"Offset" name:"Offset"`

	// 查询条数上限，默认20
	Limit *int64 `json:"Limit" name:"Limit"`
}

func (r *DescribeIndustryNewsRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeIndustryNewsRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeIndustryNewsResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 总计文章数量
		NewsCount *uint64 `json:"NewsCount" name:"NewsCount"`

		// 总计来源数量
		FromCount *uint64 `json:"FromCount" name:"FromCount"`

		// 总计疑似负面数量
		AdverseCount *uint64 `json:"AdverseCount" name:"AdverseCount"`

		// 文章列表
		NewsSet []*IndustryNews `json:"NewsSet" name:"NewsSet" list`

		// 按天统计的数量列表
		DateCountSet []*DateCount `json:"DateCountSet" name:"DateCountSet" list`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeIndustryNewsResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeIndustryNewsResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeUserPortraitRequest struct {
	*tchttp.BaseRequest

	// 品牌ID
	BrandId *string `json:"BrandId" name:"BrandId"`
}

func (r *DescribeUserPortraitRequest) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeUserPortraitRequest) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type DescribeUserPortraitResponse struct {
	*tchttp.BaseResponse
	Response *struct {

		// 年龄画像
		Age *AgePortraitInfo `json:"Age" name:"Age"`

		// 性别画像
		Gender *GenderPortraitInfo `json:"Gender" name:"Gender"`

		// 省份画像
		Province *ProvincePortraitInfo `json:"Province" name:"Province"`

		// 电影喜好画像
		Movie *MoviePortraitInfo `json:"Movie" name:"Movie"`

		// 明星喜好画像
		Star *StarPortraitInfo `json:"Star" name:"Star"`

		// 唯一请求ID，每次请求都会返回。定位问题时需要提供该次请求的RequestId。
		RequestId *string `json:"RequestId" name:"RequestId"`
	} `json:"Response"`
}

func (r *DescribeUserPortraitResponse) ToJsonString() string {
    b, _ := json.Marshal(r)
    return string(b)
}

func (r *DescribeUserPortraitResponse) FromJsonString(s string) error {
    return json.Unmarshal([]byte(s), &r)
}

type GenderPortrait struct {

	// 性别
	Gender *string `json:"Gender" name:"Gender"`

	// 百分比
	Percent *uint64 `json:"Percent" name:"Percent"`
}

type GenderPortraitInfo struct {

	// 用户性别画像数组
	PortraitSet []*GenderPortrait `json:"PortraitSet" name:"PortraitSet" list`
}

type IndustryNews struct {

	// 行业报道ID
	IndustryId *string `json:"IndustryId" name:"IndustryId"`

	// 报道发表时间
	PubTime *string `json:"PubTime" name:"PubTime"`

	// 报道来源
	FromSite *string `json:"FromSite" name:"FromSite"`

	// 报道标题
	Title *string `json:"Title" name:"Title"`

	// 报道来源url
	Url *string `json:"Url" name:"Url"`

	// 报道来源等级
	Level *uint64 `json:"Level" name:"Level"`

	// 热度值
	Hot *uint64 `json:"Hot" name:"Hot"`

	// 报道标识
	Flag *uint64 `json:"Flag" name:"Flag"`

	// 报道摘要
	Abstract *string `json:"Abstract" name:"Abstract"`
}

type MoviePortrait struct {

	// 电影名称
	Name *string `json:"Name" name:"Name"`

	// 百分比
	Percent *float64 `json:"Percent" name:"Percent"`
}

type MoviePortraitInfo struct {

	// 用户喜好电影画像数组
	PortraitSet []*MoviePortrait `json:"PortraitSet" name:"PortraitSet" list`
}

type ProvincePortrait struct {

	// 省份名称
	Province *string `json:"Province" name:"Province"`

	// 百分比
	Percent *float64 `json:"Percent" name:"Percent"`
}

type ProvincePortraitInfo struct {

	// 用户省份画像数组
	PortraitSet []*ProvincePortrait `json:"PortraitSet" name:"PortraitSet" list`
}

type StarPortrait struct {

	// 喜欢的明星名字
	Name *string `json:"Name" name:"Name"`

	// 百分比
	Percent *float64 `json:"Percent" name:"Percent"`
}

type StarPortraitInfo struct {

	// 用户喜好的明星画像数组
	PortraitSet []*StarPortrait `json:"PortraitSet" name:"PortraitSet" list`
}
