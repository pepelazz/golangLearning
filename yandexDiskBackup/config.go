package main

import (
	"github.com/pelletier/go-toml"
	"log"
)

type (
	Config struct {
		AuthToken   string
		ProjectList []Project
	}
)

func readConfig() {
	tree, err := toml.LoadFile("./config.toml")
	if err != nil {
		log.Fatalf("load config %s", err)
	}
	if !tree.Has("authToken") {
		log.Fatalf("missed authToken in config")
	}
	config.AuthToken = tree.Get("authToken").(string)
	if tree.Has("projectList") {
		data := tree.Get("projectList").([]interface{})
		config.ProjectList = []Project{}
		for i := range data {
			prj := Project{}
			for i, v := range data[i].([]interface{}) {
				switch i {
				case 0:
					prj.Name = v.(string)
				case 1:
					prj.Type = v.(string)
				case 2:
					prj.Path = v.(string)
				case 3:
					prj.DockerPgName = v.(string)
				case 4:
					prj.ExistFile = v.(string)
				}
			}
			config.ProjectList = append(config.ProjectList, prj)
		}
	}
}
