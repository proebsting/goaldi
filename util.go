//  utils.go -- general-purpose utility routines

package goaldi

import (
	"reflect"
	"sort"
)

//  SortedKeys generates in order (over a channel) the keys of a map[string].
//  usage:  for k := range SortedKeys(mymap) { ... }
func SortedKeys(m interface{}) chan string {
	vlist := reflect.ValueOf(m).MapKeys()
	n := len(vlist)
	slist := make([]string, n)
	for i, k := range vlist {
		slist[i] = k.String()
	}
	sort.Strings(slist)
	ch := make(chan string, n)
	go func() {
		for _, k := range slist {
			ch <- k
		}
		close(ch)
	}()
	return ch
}
