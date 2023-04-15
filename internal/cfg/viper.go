package cfg

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
)

func LoadConfig(conf interface{}, path ...string) {
	var config string
	if len(path) == 0 {
		flag.StringVar(&config, "config", "", "choose config file.")
		flag.Parse()
	} else { // 函数传递的可变参数的第一个值赋值于config
		config = path[0]
	}
	v := viper.New()
	v.SetConfigFile(config)
	v.SetConfigType("yaml")
	err := v.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	if err = v.Unmarshal(conf); err != nil {
		panic(err)
	}
}
