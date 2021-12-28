package web

import (
	"initial-sever/util"
	"os"
	"path"
)

type Nav struct {
	Name     string `json:"name,omitempty"`
	Path     string `json:"path,omitempty"`
	Id       int    `json:"id"`
	ParentId int    `json:"parentId"`
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

func findNav(routes map[string]struct{}) (navs []Nav) {
	navs = make([]Nav, 0, len(routes))
	for _, v := range defNav {
		if _, ok := routes[v.Name]; ok {
			navs = append(navs, v)
		}
	}
	return
}

const (
	rePath   = "/web"
	userFile = "user.json"
)

var (
	dataPath string
	defNav   []Nav
)

func LoadNav(filename string) error {
	return util.DecodeJsonFromFile(&defNav, filename)
}

func LoadData(root string, admin_ struct {
	Username string
	Password string
}) (err error) {
	dataPath = path.Join(root, rePath)
	_ = os.MkdirAll(dataPath, os.ModePerm)

	// user
	var data userData
	filename := path.Join(dataPath, userFile)
	if err = util.DecodeJsonFromFile(&data, filename); err != nil {
		if os.IsNotExist(err) {
			admin = &User{
				Username: admin_.Username,
				Password: admin_.Password,
			}
			userMap = map[string]*User{}
			return saveUser()
		} else {
			return
		}
	} else {
		admin = data.Admin
		userMap = data.UserMap
	}
	return
}

type userData struct {
	Admin   *User            `json:"admin"`
	UserMap map[string]*User `json:"user_map"`
}

func saveUser() error {
	filename := path.Join(dataPath, userFile)
	data := userData{
		Admin:   admin,
		UserMap: userMap,
	}
	return util.EncodeJsonToFile(data, filename)
}
