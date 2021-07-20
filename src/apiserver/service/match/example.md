### 第三方服务使用apiserver 路由


##### 使用服务发现发现服务

```
	RegisterService(服务名字, 实现URL匹配方法)

```

##### 实现URL匹配方法


```
func MatchMigrate(req *restful.Request) (from, to string, isHit bool) 
```




eg: 

``` golang
package match




import (
	"strings"

	"configcenter/src/common/types"

	"github.com/emicklei/go-restful"
)

func init() {
	RegisterService(types.CC_MODULE_MIGRATE, MatchMigrate)
}

func MatchMigrate(req *restful.Request) (from, to string, isHit bool) {
	uri := req.Request.RequestURI
	from, to = RootPath, "/migrate/v3"
	switch {
		
	case strings.HasPrefix(uri, RootPath+"/authcenter/init"):
	// 将请求URL未/api/v3/authcenter/init，转发到当前服务
		isHit = true
	case strings.HasPrefix(uri, RootPath+"/find/system/config_admin"):
		// 将请求URL未/api/v3/find/system/config_admin，转发到当前服务
		isHit = true

	default:
		isHit = false
	}

	return from, to, isHit
}
```