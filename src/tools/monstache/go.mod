module monstache

go 1.16

require (
	configcenter v0.0.0
	github.com/BurntSushi/toml v1.1.0
	github.com/rwynn/monstache v4.12.3+incompatible
	github.com/tidwall/gjson v1.14.1
	go.mongodb.org/mongo-driver v1.9.1
)

replace (
	configcenter v0.0.0 => ../../../../configcenter
	github.com/rwynn/monstache v4.12.3+incompatible => github.com/ZQHcode/monstache v1.0.0
)
