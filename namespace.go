//  namespace.go -- named and unnamed global variable collections

package goaldi

import ()

type Namespace struct {
	name    string           // actual name, possibly empty
	qname   string           // name:: or empty
	entries map[string]Value // mapping of names to variables
}

var allSpaces = make(map[string]*Namespace)

//  GetSpace(name) -- get or create a global namespace
//  The name may be blank to specify the default unnamed space
func GetSpace(name string) *Namespace {
	ns := allSpaces[name]
	if ns == nil {
		ns = &Namespace{}
		ns.name = name
		ns.entries = make(map[string]Value)
		if ns.name == "" {
			ns.name = ""
		} else {
			ns.qname = ns.name + "::"
		}
		allSpaces[name] = ns
	}
	return ns
}

//  Namespace.Declare(name, contents) -- initialize a namespace entry
func (ns *Namespace) Declare(name string, contents Value) {
	if ns.entries[name] != nil {
		panic(Malfunction("duplicate entry " + ns.name + "::" + name))
	}
	ns.entries[name] = contents
}

//  Namespace.GetQual() -- return "" if default space else name + "::"
func (ns *Namespace) GetQual() string {
	return ns.qname
}

//  Namespace.Get(name) -- retrieve namespace entry (or nil)
func (ns *Namespace) Get(name string) Value {
	return ns.entries[name]
}

//  Namespace.All() -- generate all names over a channel.
//  usage:  for k := range ns.All() {...}
func (ns *Namespace) All() chan string {
	return SortedKeys(ns.entries)
}

//  AllSpaces() -- generate names of all namespaces, in sorted order
//  usage:  for k := range AllSpaces() {...}
func AllSpaces() chan string {
	return SortedKeys(allSpaces)
}
