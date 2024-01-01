package container

func In[V comparable](elem V, slice []V) bool {
	for _, v := range slice {
		if v == elem {
			return true
		}
	}

	return false
}
