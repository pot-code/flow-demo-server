package rbac

allow {
	input.roles[_] == "admin"
}

allow {
	input.roles[_] == input.allow[_]
}
