package vastdata

func (rs *User) GetRestResource(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Users
}
