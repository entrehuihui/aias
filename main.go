package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	dir := flag.String("d", "./proto", "Input folder")
	aiasfile := flag.String("o", "aias.go", "out file name")
	flag.Parse()
	readGPRCJSON(*dir, *aiasfile)
}

func readGPRCJSON(dir, aiasfile string) {
	_, fileName := filepath.Split(dir)

	jfiles := []string{}
	protoName := ""
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
		ok = strings.HasSuffix(fi.Name(), ".json")
		if ok {
			protoName = protoName + "/" + fi.Name()
		}
	}
	retmap := make(map[string]interface{}, 0)
	ifmap := true
	pathsmap := make(map[string]interface{}, 0)
	definitionsmap := make(map[string]interface{}, 0)
	for _, jfile := range jfiles {
		jmap := make(map[string]interface{}, 0)
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

		v, ok := jmap["paths"].(map[string]interface{})
		if ok {
			// fmt.Println(v["/GetCommodityLists"])
			for k, v1 := range v {
				pathsmap[k] = v1
			}
		}

		v, ok = jmap["definitions"].(map[string]interface{})
		if ok {
			// fmt.Println(v["/GetCommodityLists"])
			for k, v1 := range v {
				definitionsmap[k] = v1
			}
		}
	}

	retmap["paths"] = pathsmap
	retmap["definitions"] = definitionsmap
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
