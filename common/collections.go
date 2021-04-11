package common

type EmptySetVal struct{}

var member EmptySetVal

type Set struct {
	data map[string]EmptySetVal
}

func NewSet() *Set {
	return &Set{
		data: map[string]EmptySetVal{},
	}
}

func (set *Set) Add(val string) {
	set.data[val] = member
}

func (set *Set) Del(val string) {
	delete(set.data, val)
}

func (set *Set) Len() int {
	return len(set.data)
}

func (set *Set) HasKey(key string) bool {
	_, exists := set.data[key]
	return exists
}

func (set *Set) Keys() []string {
	keys := make([]string, 0, len(set.data))
	for key := range set.data {
		keys = append(keys, key)
	}
	return keys
}

func Includes(slice []string, elem string) bool {
	for _, el := range slice {
		if el == elem {
			return true
		}
	}
	return false
}
