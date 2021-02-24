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

package configures

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strconv"
	"strings"

	"configcenter/src/common/blog"
	"configcenter/src/common/types"

	"github.com/spf13/viper"
)

// checkFile check that the configuration in the file is legal
func (cc *ConfCenter) checkFile(confFilePath string) error {
	file := path.Base(confFilePath)
	split := strings.Split(file, ".")
	fileName := split[0]
	v := viper.New()
	if runtime.GOOS == "windows" {
		nameWithPointList := strings.Split(file, `\`)
		nameWithPoint := nameWithPointList[len(nameWithPointList)-1]
		name := strings.Split(nameWithPoint, ".")[0]
		filePath := strings.TrimSuffix(confFilePath, nameWithPoint)
		v.SetConfigName(name)
		v.AddConfigPath(filePath)
	} else {
		v.SetConfigName(fileName)
		v.AddConfigPath(path.Dir(confFilePath))
	}
	err := v.ReadInConfig()
	if err != nil {
		blog.Errorf("fail to read configure from %s ", file)
		return errors.New("fail to read configure from file.")
	}

	switch fileName {
	case types.CCConfigureRedis:
		if err := cc.isRedisConfigOK("redis", file, v); err != nil {
			return err
		}
		if err := cc.isRedisConfigOK("redis.snap", file, v); err != nil {
			return err
		}
		if err := cc.isRedisConfigOK("redis.discover", file, v); err != nil {
			return err
		}
		if err := cc.isRedisConfigOK("redis.netcollect", file, v); err != nil {
			return err
		}

	case types.CCConfigureMongo:
		if err := cc.isMongoConfigOK(v, file); err != nil {
			return err
		}

	case types.CCConfigureCommon:
		// check es config
		if err := cc.isEsConfigOK(v, file); err != nil {
			return err
		}

		// check datacollection config
		if err := cc.isDatacollectionConfigOK(v, file); err != nil {
			return err
		}

		// check operation server config
		if err := cc.isOperationConfigOK(v, file); err != nil {
			return err
		}

		// check monitor config
		if err := cc.isMonitorConfigOK(v, file); err != nil {
			return err
		}

	}
	return nil
}

func (cc *ConfCenter) isRedisConfigOK(prefix, fileName string, v *viper.Viper) error {
	if err := cc.isConfigEmpty(prefix+".host", fileName, v); err != nil {
		return err
	}
	if err := cc.isConfigNotIntVal(prefix+".database", fileName, v); err != nil {
		return err
	}
	if v.IsSet(prefix + ".maxOpenConns") {
		if err := cc.isConfigNotIntVal(prefix+".maxOpenConns", fileName, v); err != nil {
			return err
		}
	}
	if v.IsSet(prefix + ".maxIDleConns") {
		if err := cc.isConfigNotIntVal(prefix+".maxIDleConns", fileName, v); err != nil {
			return err
		}
	}
	return nil
}

func (cc *ConfCenter) isMongoConfigOK(v *viper.Viper, fileName string) error {
	if err := cc.isConfigEmpty("mongodb.host", fileName, v); err != nil {
		return err
	}
	if err := cc.isConfigEmpty("mongodb.database", fileName, v); err != nil {
		return err
	}
	if err := cc.isConfigEmpty("mongodb.rsName", fileName, v); err != nil {
		return err
	}
	if v.IsSet("mongodb.maxOpenConns") {
		if err := cc.isConfigNotIntVal("mongodb.maxOpenConns", fileName, v); err != nil {
			return err
		}
	}
	if v.IsSet("mongodb.maxIdleConns") {
		if err := cc.isConfigNotIntVal("mongodb.maxIdleConns", fileName, v); err != nil {
			return err
		}
	}
	if v.IsSet("mongodb.socketTimeoutSeconds") {
		if err := cc.isConfigNotIntVal("mongodb.socketTimeoutSeconds", fileName, v); err != nil {
			return err
		}
	}
	return nil
}

func (cc *ConfCenter) isEsConfigOK(v *viper.Viper, fileName string) error {
	if v.IsSet("es.fullTextSearch") {
		fullTextSearch := v.GetString("es.fullTextSearch")
		if fullTextSearch != "on" && fullTextSearch != "off" {
			blog.Errorf("The configuration file is %s, the es.fullTextSearch should be on or off !", fileName)
			return fmt.Errorf("The configuration file is %s, the es.fullTextSearch should be on or off !", fileName)
		}
		if fullTextSearch == "on" {
			if err := cc.isConfigEmpty("es.url", fileName, v); err != nil {
				return err
			}
		}
	}
	return nil
}

func (cc *ConfCenter) isDatacollectionConfigOK(v *viper.Viper, fileName string) error {
	if v.IsSet("datacollection.hostsnap.changeRangePercent") {
		if err := cc.isConfigNotIntVal("datacollection.hostsnap.changeRangePercent", fileName, v); err != nil {
			return err
		}
	}
	if v.IsSet("datacollection.hostsnap.rateLimiter.qps") {
		if err := cc.isConfigNotIntVal("datacollection.hostsnap.rateLimiter.qps", fileName, v); err != nil {
			return err
		}
	}
	if v.IsSet("datacollection.hostsnap.rateLimiter.burst") {
		if err := cc.isConfigNotIntVal("datacollection.hostsnap.rateLimiter.burst", fileName, v); err != nil {
			return err
		}
	}
	if v.IsSet("datacollection.hostsnap.timeWindow.atTime") {
		if err := cc.isTimeFormat("datacollection.hostsnap.timeWindow.atTime", fileName, v); err != nil {
			return err
		}
	}
	if v.IsSet("datacollection.hostsnap.timeWindow.checkIntervalHours") {
		if err := cc.isConfigNotIntVal("datacollection.hostsnap.timeWindow.checkIntervalHours", fileName, v); err != nil {
			return err
		}
	}

	if v.IsSet("datacollection.hostsnap.timeWindow.windowMinutes") {
		if err := cc.isConfigNotIntVal("datacollection.hostsnap.timeWindow.windowMinutes", fileName, v); err != nil {
			return err
		}
	}
	return nil
}

func (cc *ConfCenter) isOperationConfigOK(v *viper.Viper, fileName string) error {
	if !v.IsSet("operationServer.disableOperationStatistic") {
		return nil
	}
	if err := cc.isConfigNotBoolVal("operationServer.disableOperationStatistic", fileName, v); err != nil {
		return err
	}
	return nil
}

func (cc *ConfCenter) isMonitorConfigOK(v *viper.Viper, fileName string) error {
	if err := cc.isConfigKeyExist("monitor", fileName, v); err != nil {
		return err
	}

	if v.IsSet("monitor.enableMonitor") {
		if err := cc.isConfigNotBoolVal("monitor.enableMonitor", fileName, v); err != nil {
			return err
		}
	}

	if v.IsSet("monitor.pluginName") {
		pluginName := v.GetString("monitor.pluginName")
		if pluginName == "blueking" {
			if err := cc.isConfigEmpty("monitor.dataID", fileName, v); err != nil {
				return err
			}
			if dataID := v.GetInt64("monitor.dataID"); dataID <= 0 {
				blog.Errorf("The configuration file is %s, the %s must be an integer >0 when plugin name is blueking !", fileName, "monitor.dataID")
				return fmt.Errorf("The configuration file is %s, the %s must be an integer >0 when plugin name is blueking !", fileName, "monitor.dataID")
			}
		}
	}

	return nil
}

func (cc *ConfCenter) isTimeFormat(configName, fileName string, v *viper.Viper) error {
	atTime := v.GetString(configName)
	timeVal := strings.Split(atTime, ":")
	if len(timeVal) != 2 {
		blog.Errorf("The configuration file is %s, the format of %s is wrong !", fileName, configName)
		return fmt.Errorf("The configuration file is %s, the format of %s is wrong !", fileName, configName)
	}
	hour, err := strconv.Atoi(timeVal[0])
	if err != nil || hour < 0 || hour > 24 {
		blog.Errorf("The configuration file is %s, the format of %s is wrong !", fileName, configName)
		return fmt.Errorf("The configuration file is %s, the format of %s is wrong !", fileName, configName)
	}
	min, err := strconv.Atoi(timeVal[1])
	if err != nil || min < 0 || min > 60 {
		blog.Errorf("The configuration file is %s, the format of %s is wrong !", fileName, configName)
		return fmt.Errorf("The configuration file is %s, the format of %s is wrong !", fileName, configName)
	}
	if min > 0 && hour == 24 {
		blog.Errorf("The configuration file is %s, the format of %s is wrong !", fileName, configName)
		return fmt.Errorf("The configuration file is %s, the format of %s is wrong !", fileName, configName)
	}
	return nil
}

func (cc *ConfCenter) isConfigKeyExist(configName, fileName string, v *viper.Viper) error {
	if !v.InConfig(configName) {
		blog.Errorf("The configuration file is %s, the %s must exist !", fileName, configName)
		return fmt.Errorf("The configuration file is %s, the %s must exist !", fileName, configName)
	}
	return nil
}

func (cc *ConfCenter) isConfigEmpty(configName, fileName string, v *viper.Viper) error {
	if !v.IsSet(configName) {
		blog.Errorf("The configuration file is %s, the %s can not be empty !", fileName, configName)
		return fmt.Errorf("The configuration file is %s, the %s can not be empty !", fileName, configName)
	}
	return nil
}

func (cc *ConfCenter) isConfigNotIntVal(configName, fileName string, v *viper.Viper) error {
	val := v.GetString(configName)
	_, err := strconv.Atoi(val)
	if err != nil {
		blog.Errorf("The configuration file is %s, the %s should be a string of numbers !", fileName, configName)
		return fmt.Errorf("The configuration file is %s, the %s should be a string of numbers !", fileName, configName)
	}
	return nil
}

func (cc *ConfCenter) isConfigNotBoolVal(configName, fileName string, v *viper.Viper) error {
	val := v.GetString(configName)
	if val != "true" && val != "false" {
		blog.Errorf("The configuration file is %s, the %s should be true or false !", fileName, configName)
		return fmt.Errorf("The configuration file is %s, the %s should be true or false !", fileName, configName)
	}
	return nil
}
