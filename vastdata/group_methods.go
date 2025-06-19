package vastdata

func (rs *Group) GetRestResource(rest *VMSRest) VastResourceAPIWithContext {
	return rest.Groups
}
