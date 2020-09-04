package hubsub

// MapFromPairs returns map[string]string from paires key-values.
func MapFromPairs(paires ...string) map[string]string {
	if len(paires) == 0 {
		return nil
	}
	if len(paires)%2 != 0 {
		return nil
	}
	res := make(map[string]string, len(paires)/2)
	for i := 0; i < len(paires); i += 2 {
		res[paires[i]] = paires[i+1]
	}
	return res
}

// FilterFromPairs returns filter function from paires key-values.
func FilterFromPairs(paires ...string) func(map[string]string) bool {
	if len(paires) == 0 {
		return nil
	}
	if len(paires)%2 != 0 {
		return nil
	}
	return func(in map[string]string) bool {
		for i := 0; i < len(paires); i += 2 {
			if in[paires[i]] != paires[i+1] {
				return false
			}
		}
		return true
	}
}
