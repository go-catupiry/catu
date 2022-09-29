package acl_test

import (
	"reflect"
	"testing"

	"github.com/go-catupiry/catu/acl"
	"github.com/stretchr/testify/assert"
)

func TestRole(t *testing.T) {
	r, _ := acl.NewRole(&acl.NewRoleOpts{Name: "editor"})
	assert.Equal(t, 0, len(r.Permissions))
	assert.False(t, r.Can("find_image"))

	r.AddPermission("upload_image")
	r.AddPermission("find_image")

	assert.Equal(t, 2, len(r.Permissions))
	assert.True(t, r.Can("find_image"))

	r.RemovePermission("find_image")

	assert.Equal(t, 1, len(r.Permissions))
	assert.False(t, r.Can("find_image"))
}

func TestNewRole(t *testing.T) {
	type args struct {
		opts *acl.NewRoleOpts
	}
	tests := []struct {
		name    string
		args    args
		want    *acl.Role
		wantErr bool
	}{
		{
			"success empty",
			args{opts: &acl.NewRoleOpts{
				Name: "faxineira",
			}},
			&acl.Role{
				Name: "faxineira",
			},
			false,
		},
		{
			"error no name",
			args{opts: &acl.NewRoleOpts{}},
			nil,
			true,
		},
		{
			"success with permissions",
			args{opts: &acl.NewRoleOpts{
				Name:        "porteiro",
				Permissions: []string{"block-user-access"},
			}},
			&acl.Role{
				Name:        "porteiro",
				Permissions: []string{"block-user-access"},
			},
			false,
		},
		{
			"success with systemRole",
			args{opts: &acl.NewRoleOpts{
				Name:         "editor",
				IsSystemRole: true,
			}},
			&acl.Role{
				Name:         "editor",
				IsSystemRole: true,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := acl.NewRole(tt.args.opts)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRole() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_Can(t *testing.T) {
	editoRole := acl.Role{
		Name:          "editor",
		Permissions:   []string{"update_image", "find_image", "create_content"},
		CanAddInUsers: true,
		IsSystemRole:  true,
	}

	type args struct {
		permission string
	}
	tests := []struct {
		name string
		r    *acl.Role
		args args
		want bool
	}{
		{
			"can find_image",
			&editoRole,
			args{permission: "find_image"},
			true,
		},
		{
			"cant delete_content",
			&editoRole,
			args{permission: "delete_content"},
			false,
		},
		{
			"cant create_content",
			&editoRole,
			args{permission: "create_content"},
			true,
		},
		{
			"lixeiro cant jogar_lixo_na_rua",
			&acl.Role{
				Name:          "lixeiro",
				Permissions:   []string{"pegar_lixo"},
				CanAddInUsers: true,
				IsSystemRole:  false,
			},
			args{permission: "jogar_lixo_na_rua"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.r.Can(tt.args.permission); got != tt.want {
				t.Errorf("Role.Can() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRole_AddPermission(t *testing.T) {
	type args struct {
		permission string
	}
	tests := []struct {
		name string
		r    *acl.Role
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.AddPermission(tt.args.permission)
		})
	}
}

func TestRole_RemovePermission(t *testing.T) {
	type args struct {
		permission string
	}
	tests := []struct {
		name string
		r    *acl.Role
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.r.RemovePermission(tt.args.permission)
		})
	}
}

func TestLoadRoles(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := acl.LoadRoles()
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadRoles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("LoadRoles() = %v, want %v", got, tt.want)
			}
		})
	}
}
