package utils

type ResponseConversionFunc func(map[string]interface{}) map[string]interface{}

func EntityMergeToUserQuotas(m map[string]interface{}) map[string]interface{} {
	/*This function should handle the case of the Quota object where sending is defferant than reading sturctue
	  to move the fields from the entity object into the user quotas
	*/
	for _, key := range []string{"user_quotas", "group_quotas"} {
		quotas, exists := m[key]
		if exists {
			old_quotas := quotas.([]interface{})
			new_quotas := []map[string]interface{}{}
			for _, quota := range old_quotas {
				new_quota := make(map[string]interface{})
				_quota := quota.(map[string]interface{})
				entity, entity_exists := _quota["entity"]
				if entity_exists {
					for k, v := range entity.(map[string]interface{}) {
						new_quota[k] = v
					}
				}
				for k, v := range _quota {
					if k == "entity" {
						continue
					}
					new_quota[k] = v
				}

				new_quotas = append(new_quotas, new_quota)
			}
			m[key] = new_quotas
		}
	}
	return m
}

func EnabledMustBeSet(m map[string]interface{}) map[string]interface{} {
	_, exists := m["enabled"]

	if !exists {
		m["enabled"] = false
	}
	return m

}
