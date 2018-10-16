package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
	"sync"
)

// XXX: really, we could have a separate lock for each file, but this will probably do
var listsLock sync.Mutex

func IsInList(s string, list []string) bool {
	for _, item := range list {
		if item == s {
			return true
		}
	}

	return false
}

func GetList(name string) []string {
	listsLock.Lock()
	defer listsLock.Unlock()

	return _unsafe_GetList(name)
}

func _unsafe_GetList(name string) []string {
	bytes, err := ioutil.ReadFile(ListFilename(name))
	if err != nil {
		return make([]string, 0)
	}

	items := strings.Split(string(bytes), "\n")
	if len(items) > 0 && len(items[len(items)-1]) == 0 {
		// remove empty string after trailing endline
		items = items[:len(items)-1]
	}
	return items
}

func ListFilename(name string) string {
	return viper.GetString("datadir") + "/" + name + ".list"
}

func _unsafe_WriteList(name string, list []string) {
	filename := ListFilename(name)
	contents := strings.Join(list, "\n")
	if len(contents) > 0 {
		// add trailing endline unless the list is empty
		contents += "\n"
	}

	err := ioutil.WriteFile(filename, []byte(contents), 0644)
	if err != nil {
		fmt.Printf("error writing to %s: %v", filename, err)
	}
}

func AddToList(name string, onion string) {
	listsLock.Lock()
	defer listsLock.Unlock()

	l := _unsafe_GetList(name)
	if IsInList(onion, l) {
		return
	}

	l = append(l, onion)
	_unsafe_WriteList(name, l)
}

func RemoveFromList(name string, onion string) {
	listsLock.Lock()
	defer listsLock.Unlock()

	l := _unsafe_GetList(name)
	i := 0
	for _, s := range l {
		// XXX: it's plausible that the entry will appear more than once; we want to delete all of them
		if s != onion {
			l[i] = s
			i++
		}
	}
	l = l[:i]
	_unsafe_WriteList(name, l)
}
