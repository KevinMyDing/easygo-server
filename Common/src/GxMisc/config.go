/**
作者:guangbo
模块：配置读写接口模块
说明：
创建时间：2015-10-30
**/
package GxMisc

import (
	. "github.com/bitly/go-simplejson"
	"os"
)

var Config *Json

func LoadConfig(filename string) error {
	buf := make([]byte, 1024)
	f, err := os.Open(filename)
	if err != nil {
		// fmt.Printf("Error: %s\n", err)
		return err
	}
	defer f.Close()
	f.Read(buf)
	//
	Config, _ = NewJson(buf)
	return nil
}
