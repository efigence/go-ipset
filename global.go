package ipset

import "sync"

type renameList struct {
	sync.Mutex
	renameHooks map[string]func(newName string)
}

var registry = renameList{
	renameHooks: map[string]func(newName string){},
}

func (r *renameList) Swap(name1, name2 string) {
	n1 := r.renameHooks[name1]
	n2 := r.renameHooks[name2]
	r.renameHooks[name1] = n2
	r.renameHooks[name2] = n1
}
