package catu

func (r *App) Can(permission string, userRoles []string) bool {
	// first check if user is administrator
	for i := range userRoles {
		if userRoles[i] == "administrator" {
			return true
		}
	}

	for j := range userRoles {
		R := r.RolesList[userRoles[j]]
		if R.Can(permission) {
			return true
		}
	}

	return false
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
