package acl

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
