//  namespace.go -- named and unnamed global variable collections

package goaldi

import ()

type Namespace struct {
	Name    string           // actual name, possibly empty
	Qname   string           // identifier:: or empty
	Entries map[string]Value // mapping of names to variables
}

var allSpaces = make(map[string]*Namespace)

//  GetSpace(name) -- get or create a global namespace
//  The name may be blank to specify the default unnamed space
func GetSpace(name string) *Namespace {
	ns := allSpaces[name]
	if ns == nil {
		ns = &Namespace{}
		ns.Name = name
		ns.Entries = make(map[string]Value)
		if name != "" {
			ns.Qname = name + "::"
		}
		allSpaces[name] = ns
	}
	return ns
}

//  Namespace.Declare(name, contents) -- initialize a namespace entry
func (ns *Namespace) Declare(name string, contents Value) {
	if ns.Entries[name] != nil {
		panic(Malfunction("duplicate entry " + ns.Qname + name))
	}
	ns.Entries[name] = contents
}

//  Namespace.GetQual() -- return "" if default space else name + "::"
func (ns *Namespace) GetQual() string {
	return ns.Qname
}

//  Namespace.Get(name) -- retrieve namespace entry (or nil)
func (ns *Namespace) Get(name string) Value {
	return ns.Entries[name]
}

//  Namespace.All() -- generate all names over a channel.
//  usage:  for k := range ns.All() {...}
func (ns *Namespace) All() chan string {
	return SortedKeys(ns.Entries)
}

//  AllSpaces() -- generate names of all namespaces, in sorted order
//  usage:  for k := range AllSpaces() {...}
func AllSpaces() chan string {
	return SortedKeys(allSpaces)
}
