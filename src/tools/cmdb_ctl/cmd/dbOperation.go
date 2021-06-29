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

package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"configcenter/src/common"
	"configcenter/src/common/mapstr"
	"configcenter/src/storage/dal/mongo/local"
	"configcenter/src/tools/cmdb_ctl/app/config"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

const (
	maxDeleteNum      = 1000
	maxDeleteBatchNum = 300
)
const (
	defaultClientMaxOpenConns     = 10
	minimumClientMaxIdleOpenConns = 5
)

func init() {
	rootCmd.AddCommand(NewDbOperationCommand())
}

type delDbConf struct {
	colName   string
	condition string
}

type findDbConf struct {
	colName   string
	condition string
	resfilter string
	num       int32
	bPretty   bool
}

type dbOperationConf struct {
	service   *config.Service
	delParam  delDbConf
	findParam findDbConf
}
type delData struct {
	MongoID primitive.ObjectID `bson:"_id"`
}

//  db
//    --find
//             --colName(collection name) --condition（查询的条件） --resfilter（结果是否需要过滤指定字段） --pretty（是否需要采用json pretty格式返回） --num（返回的数量默认值是5）
//    --delete
//             --colName（collection name）--condition（删除的条件）
//    --show

func NewDbOperationCommand() *cobra.Command {

	conf := new(dbOperationConf)
	cmd := &cobra.Command{
		Use:   "db",
		Short: "db operations",
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	findCmds := make([]*cobra.Command, 0)
	findCmd := &cobra.Command{
		Use:   "find",
		Short: "find eligible data from the db",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runFindDbDataCmd(conf)
		},
	}
	findCmd.Flags().StringVar(&conf.findParam.colName, "collection", "", "collection name,the param must be assigned")
	findCmd.Flags().StringVar(&conf.findParam.condition, "condition", "", "query conditions ,the parameter must be json format string")
	findCmd.Flags().StringVar(&conf.findParam.resfilter, "resfilter", "", "display the required fields ,the fieds link with comma")
	findCmd.Flags().Int32Var(&conf.findParam.num, "num", 5, "numbers of result to show ,default num is 5 ")
	findCmd.Flags().BoolVar(&conf.findParam.bPretty, "pretty", false, "query result are displayed in json pretty format")

	findCmds = append(findCmds, findCmd)
	for _, fCmd := range findCmds {
		cmd.AddCommand(fCmd)
	}

	delCmds := make([]*cobra.Command, 0)
	delCmd := &cobra.Command{
		Use:   "delete",
		Short: "delete eligible data from the db",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelDbDataCmd(conf)
		},
	}

	delCmd.Flags().StringVar(&conf.delParam.colName, "collection", "", "collection name,the parameter must be assigned")
	delCmd.Flags().StringVar(&conf.delParam.condition, "condition", "", "conditions for deletion,the parameter must be json format string")
	delCmds = append(delCmds, delCmd)
	for _, dCmd := range delCmds {
		cmd.AddCommand(dCmd)
	}

	showCmds := make([]*cobra.Command, 0)
	showCmd := &cobra.Command{
		Use:   "show",
		Short: "show all collections",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runShowDbDataCmd()
		},
	}

	showCmds = append(showCmds, showCmd)
	for _, cCmd := range showCmds {
		cmd.AddCommand(cCmd)
	}

	return cmd
}

func runDelDbDataCmd(conf *dbOperationConf) error {

	s, err := newMongo(config.Conf.MongoURI, config.Conf.MongoRsName)
	if err != nil {
		fmt.Printf("connect mongo db fail ,err: %v\n", err)
		return err
	}
	defer s.DbProxy.Close()

	cond, err := mapstr.NewFromInterface(conf.delParam.condition)
	if err != nil {
		fmt.Printf("condition  convert to MapStr fail ,err: %v\n", err)
		return err
	}

	ctx := context.Background()
	total, err := s.DbProxy.Table(conf.delParam.colName).Find(cond).Count(ctx)
	if err != nil {
		fmt.Printf("connect mongo db fail ,err: %v\n", err)
		return err
	}

	if total > maxDeleteNum {
		errInfo := fmt.Sprintf("number of data to delete is %d,over the max delete number 1000.", total)
		return errors.New(errInfo)
	}

	if total < maxDeleteBatchNum {
		if err = s.DbProxy.Table(conf.delParam.colName).Delete(ctx, cond); err != nil {
			fmt.Printf("delete data failed, err: %s \n", err.Error())
			return err
		}
	} else {
		var (
			start   int
			delCond map[string]interface{}
		)
		dataArr := make([]delData, 0)
		err := s.DbProxy.Table(conf.delParam.colName).Find(cond).Sort("_id").Fields("_id").Start(0).
			Limit(common.BKMaxPageSize).All(ctx, &dataArr)
		if err != nil {
			fmt.Printf("find previous del archive data failed, err: %v \n", err)
			return err
		}
		if len(dataArr) == 0 {
			fmt.Printf("no eligible data was found to be deleted .\n")
			return nil
		}

		delMongoIDs := make([]primitive.ObjectID, len(dataArr))
		for index, data := range dataArr {
			delMongoIDs[index] = data.MongoID
		}

		for {
			if start >= len(delMongoIDs) {
				break
			}
			if start+maxDeleteBatchNum > len(delMongoIDs) {
				delCond = map[string]interface{}{
					"_id": map[string]interface{}{common.BKDBIN: delMongoIDs[start:]},
				}
			} else {
				delCond = map[string]interface{}{
					"_id": map[string]interface{}{common.BKDBIN: delMongoIDs[start : start+maxDeleteBatchNum]},
				}
			}

			if err := s.DbProxy.Table(conf.delParam.colName).Delete(ctx, delCond); err != nil {
				fmt.Printf("delete previous del archive data failed, err: %v \n", err)
				return err
			}
			time.Sleep(50 * time.Millisecond)
			start = start + maxDeleteBatchNum
		}

	}

	fmt.Printf(" delete total data num is %d\n", total)

	return nil
}

func runFindDbDataCmd(conf *dbOperationConf) error {

	s, err := newMongo(config.Conf.MongoURI, config.Conf.MongoRsName)
	if err != nil {
		fmt.Printf("connect mongo db fail ,err: %v\n", err)
		return err
	}
	defer s.DbProxy.Close()

	cond, err := mapstr.NewFromInterface(conf.findParam.condition)
	if err != nil {
		fmt.Printf("condition  convert to MapStr fail ,err: %v\n", err)
		return err
	}

	filter := strings.Split(conf.findParam.resfilter, ",")
	resultMany := make([]map[string]interface{}, 0)

	if err = s.DbProxy.Table(conf.findParam.colName).Find(cond).Fields(filter...).Limit(uint64(conf.findParam.num)).
		Sort("create_time").All(context.Background(), &resultMany); err != nil {
		return fmt.Errorf("find the result from db failed, %+v", err)
	}

	dbJSON, err := json.Marshal(resultMany)
	if err != nil {
		fmt.Printf("condition  convert to MapStr fail ,err: %v\n", err)
		return err
	}

	if conf.findParam.bPretty {
		var out bytes.Buffer
		err = json.Indent(&out, dbJSON, "", "    ")
		if err != nil {
			fmt.Printf("condition  convert to MapStr fail ,err: %v\n", err)
			return err
		}
		out.WriteTo(os.Stdout)
		fmt.Printf("\n")

	} else {
		fmt.Printf("%s\n", dbJSON)
	}

	total, totalerr := s.DbProxy.Table(conf.findParam.colName).Find(cond).Fields(filter...).Count(context.Background())
	if totalerr != nil {
		fmt.Printf("find the total data num is something wrong err: %v \n", totalerr)
	} else {
		fmt.Printf("total data num is %d \n", total)
	}

	return nil
}

func runShowDbDataCmd() error {

	s, err := newMongo(config.Conf.MongoURI, config.Conf.MongoRsName)
	if err != nil {
		fmt.Printf("connect mongo db fail ,err: %v\n", err)
		return err
	}
	defer s.DbProxy.Close()

	connStr, err := connstring.Parse(config.Conf.MongoURI)
	if nil != err {
		fmt.Printf("parse dbname fail ,err: %v\n", err)
		return err
	}

	mongo, ok := s.DbProxy.(*local.Mongo)
	if !ok {
		return fmt.Errorf("db is not local.Mongo type")
	}

	cols, err := mongo.GetDBClient().Database(connStr.Database).ListCollectionNames(context.Background(), bson.D{})
	if err != nil {
		fmt.Printf("get collections fail ,err: %v\n", err)
		return err
	}
	if len(cols) == 0 {
		fmt.Printf("the db have no collections !\n")
		return nil
	}
	for _, col := range cols {
		fmt.Printf("%s\n", col)
	}

	fmt.Printf("total collection num is %d \n", len(cols))

	return nil
}

func newMongo(mongoURI string, mongoRsName string) (*config.Service, error) {

	if mongoURI == "" {
		return nil, errors.New("mongo-uri must set via flag or environment variable")
	}

	mongoConfig := local.MongoConf{
		MaxOpenConns: defaultClientMaxOpenConns,
		MaxIdleConns: minimumClientMaxIdleOpenConns,
		URI:          mongoURI,
		RsName:       mongoRsName,
	}
	
	db, err := local.NewMgo(mongoConfig, time.Minute)
	if err != nil {
		return nil, err
	}

	return &config.Service{
		DbProxy: db,
	}, nil
}
