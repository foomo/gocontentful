package erm

// SliceIncludes returns true if slice includes string
func SliceIncludes(slice []string, key string) bool {
	for _, val := range slice {
        if val == key {
            return true
        }
    }
    return false
}