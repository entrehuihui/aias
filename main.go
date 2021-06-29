package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"test/swagger"
)

func main() {
	dir := flag.String("d", "./proto", "Input folder")
	aiasfile := flag.String("o", "aias.go", "out file name")
	flag.Parse()
	readGPRCJSON(*dir, *aiasfile)
	getFileName()
}

// 复制文件
func getFileName() {
	fmt.Println("创建swagger-ui文件")
	fileNameList := swagger.AssetNames()
	for _, v := range fileNameList {
		err := saveFile(v)
		if err != nil {
			log.Println(err)
		}
	}
}
func saveFile(filename string) error {
	data, err := swagger.Asset(filename)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, data, 0666)
	return err
}

func readGPRCJSON(dir, aiasfile string) {
	_, fileName := filepath.Split(dir)

	jfiles := []string{}
	protoName := "*.json"
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, fi := range files {
		if fi.IsDir() {
			continue
		}
		ok := strings.HasSuffix(fi.Name(), ".json")
		if !ok {
			continue
		}
		jfiles = append(jfiles, dir+"/"+fi.Name())
	}
	retmap := make(map[string]interface{})
	ifmap := true
	pathsmap := make(map[string]interface{})
	definitionsmap := make(map[string]interface{})
	for _, jfile := range jfiles {
		jmap := make(map[string]interface{})
		buf, err := ioutil.ReadFile(jfile)
		if err != nil {
			continue
		}
		err = json.Unmarshal(buf, &jmap)
		if err != nil {
			fmt.Println(err)
			continue
		}

		if ifmap {
			ifmap = false
			retmap = jmap
			v, ok := jmap["info"].(map[string]interface{})
			if ok {
				v["title"] = protoName
			}
		}
		paths, ok := jmap["paths"].(map[string]interface{})
		if ok {
			for k1, v1 := range paths {
				method, ok := v1.(map[string]interface{})
				if ok {
					for _, v2 := range method {
						value, ok := v2.(map[string]interface{})
						if ok {
							value["security"] = []interface{}{
								map[string][]interface{}{
									"authorization": make([]interface{}, 0),
								},
							}
						}
					}
				}
				pathsmap[k1] = v1
			}
		}

		definitions, ok := jmap["definitions"].(map[string]interface{})
		if ok {
			// fmt.Println(v["/GetCommodityLists"])
			for k, v1 := range definitions {
				definitionsmap[k] = v1
			}
		}
	}

	retmap["paths"] = pathsmap
	retmap["definitions"] = definitionsmap
	// 添加http/https选择
	retmap["schemes"] = []string{"https",
		"http"}
	// 验证密钥
	retmap["securityDefinitions"] = map[string]interface{}{
		"authorization": map[string]string{
			"type": "apiKey",
			"name": "authorization",
			"in":   "header",
		},
	}
	// fmt.Println(retmap)
	retbuf, err := json.Marshal(retmap)
	if err != nil {
		log.Fatal(err)
	}
	s := `package ` + fileName + `

// AiasJSON --
var AiasJSON =` + "`" + string(retbuf) + "`"

	// return bytes.NewReader(retbuf)
	ioutil.WriteFile(dir+"/"+aiasfile, []byte(s), 0777)
}
