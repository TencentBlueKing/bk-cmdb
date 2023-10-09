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

package command

import (
	"fmt"

	"configcenter/src/common/metadata"
	"configcenter/src/common/util"
)

type processPropertyOpF func(processProperty *metadata.ProcessProperty, key string, val interface{}) error

func procNumOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		procNum, err := util.GetInt64ByInterface(val)
		if err != nil {
			return fmt.Errorf("%s not integer. val:%s", key, val)
		}
		processProperty.ProcNum.Value = &procNum
		blTrue := true
		processProperty.ProcNum.AsDefaultValue = &blTrue
		if err := processProperty.ProcNum.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func stopCmdOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		stopCmd, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.StopCmd.Value = &stopCmd
		blTrue := true
		processProperty.StopCmd.AsDefaultValue = &blTrue
		if err := processProperty.StopCmd.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func restartCmdOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		restartCmd, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.RestartCmd.Value = &restartCmd
		blTrue := true
		processProperty.RestartCmd.AsDefaultValue = &blTrue
		if err := processProperty.RestartCmd.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func faceStopCmdOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		restartCmd, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.RestartCmd.Value = &restartCmd
		blTrue := true
		processProperty.RestartCmd.AsDefaultValue = &blTrue
		if err := processProperty.RestartCmd.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func bkFuncNameOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		funcName, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.FuncName.Value = &funcName
		blTrue := true
		processProperty.FuncName.AsDefaultValue = &blTrue
		if err := processProperty.FuncName.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func workPathOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		workPath, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.WorkPath.Value = &workPath
		blTrue := true
		processProperty.WorkPath.AsDefaultValue = &blTrue
		if err := processProperty.WorkPath.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func bindIpOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		return nil
	}
}

func priorityOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		priority, err := util.GetInt64ByInterface(val)
		if err != nil {
			return fmt.Errorf("%s not integer. val:%s", key, val)
		}
		processProperty.Priority.Value = &priority
		blTrue := true
		processProperty.Priority.AsDefaultValue = &blTrue
		if err := processProperty.Priority.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func reloadCmdOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		reloadCmd, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.ReloadCmd.Value = &reloadCmd
		blTrue := true
		processProperty.ReloadCmd.AsDefaultValue = &blTrue
		if err := processProperty.ReloadCmd.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func bkProcessNameOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		procName, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.ProcessName.Value = &procName
		blTrue := true
		processProperty.ProcessName.AsDefaultValue = &blTrue
		if err := processProperty.ProcessName.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func portOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		return nil
	}
}

func pidFileOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		pidFile, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.PidFile.Value = &pidFile
		blTrue := true
		processProperty.PidFile.AsDefaultValue = &blTrue
		if err := processProperty.PidFile.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func autoStartOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		autoStart, ok := val.(bool)
		if !ok {
			return fmt.Errorf("%s not boolean. val:%s", key, val)
		}
		processProperty.AutoStart.Value = &autoStart
		blTrue := true
		processProperty.AutoStart.AsDefaultValue = &blTrue
		if err := processProperty.AutoStart.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func bkStartCheckSecsOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		startCheckSecs, err := util.GetInt64ByInterface(val)
		if err != nil {
			return fmt.Errorf("%s not integer. val:%s", key, val)
		}
		processProperty.StartCheckSecs.Value = &startCheckSecs
		blTrue := true
		processProperty.StartCheckSecs.AsDefaultValue = &blTrue
		if err := processProperty.StartCheckSecs.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func startCmdOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		startCmd, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.StartCmd.Value = &startCmd
		blTrue := true
		processProperty.StartCmd.AsDefaultValue = &blTrue
		if err := processProperty.StartCmd.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func userOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		user, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.User.Value = &user
		blTrue := true
		processProperty.User.AsDefaultValue = &blTrue
		if err := processProperty.User.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func timeoutOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		timeout, err := util.GetInt64ByInterface(val)
		if err != nil {
			return fmt.Errorf("%s not integer. val:%s", key, val)
		}
		processProperty.TimeoutSeconds.Value = &timeout
		blTrue := true
		processProperty.TimeoutSeconds.AsDefaultValue = &blTrue
		if err := processProperty.TimeoutSeconds.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func protocolOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		return nil
	}
}

func descriptionOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		desc, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.Description.Value = &desc
		blTrue := true
		processProperty.Description.AsDefaultValue = &blTrue
		if err := processProperty.Description.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}

func bkStartParamRegexOpFunc() processPropertyOpF {
	return func(processProperty *metadata.ProcessProperty, key string, val interface{}) error {
		regex, ok := val.(string)
		if !ok {
			return fmt.Errorf("%s not string. val:%s", key, val)
		}
		processProperty.StartParamRegex.Value = &regex
		blTrue := true
		processProperty.StartParamRegex.AsDefaultValue = &blTrue
		if err := processProperty.StartParamRegex.Validate(); err != nil {
			return fmt.Errorf("%s illegal. val:%s. err:%s", key, val, err.Error())
		}

		return nil
	}
}
