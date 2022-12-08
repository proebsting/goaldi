//  utils.go -- general-purpose utility routines

package runtime

import (
	"reflect"
	"sort"
)

// AllKeys generates (over a channel) the keys of a map[string].
// usage:  for k := range AllKeys(mymap) { ... }
func AllKeys(m interface{}) chan string {
	return genKeys(m, false)
}

// SortedKeys generates in order (over a channel) the keys of a map[string].
// usage:  for k := range SortedKeys(mymap) { ... }
func SortedKeys(m interface{}) chan string {
	return genKeys(m, true)
}

// genKeys does the actual work for AllKeys and SortedKeys.
func genKeys(m interface{}, doSort bool) chan string {
	vlist := reflect.ValueOf(m).MapKeys()
	n := len(vlist)
	slist := make([]string, n)
	for i, k := range vlist {
		slist[i] = k.String()
	}
	if doSort {
		sort.Strings(slist)
	}
	ch := make(chan string, n)
	go func() {
		for _, k := range slist {
			ch <- k
		}
		close(ch)
	}()
	return ch
}
