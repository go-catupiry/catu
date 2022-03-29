package acl

import (
	"io/ioutil"
)

type Role struct {
	Name          string   `json:"name"`
	Permissions   []string `json:"permissions"`
	CanAddInUsers bool     `json:"canAddInUsers"`
	IsSystemRole  bool     `json:"isSystemRole"`
}

func (r *Role) Can(permission string) bool {
	for i := range r.Permissions {
		if permission == r.Permissions[i] {
			return true
		}
	}

	return false
}

func LoadRoles() (string, error) {
	aclFileName := "acl.json"

	b, err := ioutil.ReadFile(aclFileName)
	if err != nil {
		return defaultRoles, nil
	} else {
		return string(b), nil
	}
}

var defaultRoles = `{
	"administrator": {
		"name": "administrator",
		"permissions": [],
		"canAddInUsers": true,
		"isSystemRole": true
	},
	"authenticated": {
		"name": "authenticated",
		"permissions": [],
		"isSystemRole": true
	},
	"unAuthenticated": {
		"name": "unAuthenticated",
		"permissions": [],
		"isSystemRole": true
	},
	"owner": {
		"name": "owner",
		"permissions": [],
		"isSystemRole": true
	},
	"premium": {
		"name": "premium",
		"permissions": [],
		"canAddInUsers": true,
		"isSystemRole": false
	},
	"blog_editor": {
		"name": "blog_editor",
		"canAddInUsers": false,
		"permissions": []
	}}`
