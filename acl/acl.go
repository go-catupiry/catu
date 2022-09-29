package acl

import (
	"io/ioutil"

	"github.com/pkg/errors"
)

type NewRoleOpts struct {
	Name          string
	Permissions   []string
	CanAddInUsers bool
	IsSystemRole  bool
}

func NewRole(opts *NewRoleOpts) (*Role, error) {
	if opts.Name == "" {
		return nil, errors.New("NewRole name is required")
	}

	r := Role{
		Name:          opts.Name,
		CanAddInUsers: opts.CanAddInUsers,
		Permissions:   opts.Permissions,
		IsSystemRole:  opts.IsSystemRole,
	}

	return &r, nil
}

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

func (r *Role) AddPermission(permission string) {
	for i := range r.Permissions {
		if permission == r.Permissions[i] {
			return
		}
	}

	r.Permissions = append(r.Permissions, permission)
}

func (r *Role) RemovePermission(permission string) {
	for i := range r.Permissions {
		if permission == r.Permissions[i] {
			r.Permissions = append(r.Permissions[:i], r.Permissions[i+1:]...)
			return
		}
	}

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
