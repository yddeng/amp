package web

import "initial-sever/util"

type Nav struct {
	Name     string `json:"name,omitempty"`
	Path     string `json:"path,omitempty"`
	Id       int    `json:"id,omitempty"`
	ParentId int    `json:"parentId,omitempty"`
	Meta     struct {
		Title        string `json:"title,omitempty"`
		Icon         string `json:"icon,omitempty"`
		Show         bool   `json:"show,omitempty"`
		HideHeader   bool   `json:"hideHeader,omitempty"`
		HideChildren bool   `json:"hideChildren,omitempty"`
	} `json:"meta,omitempty"`
	Redirect  string `json:"redirect,omitempty"`
	Component string `json:"component,omitempty"`
}

var defNav []Nav

func LoadNav(filename string) error {
	return util.DecodeJsonFromFile(&defNav, filename)
}
