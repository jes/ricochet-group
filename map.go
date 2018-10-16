package main

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
	"sync"
)

// XXX: just like in list.go: could have one per map, but we probably don't care
var mapsLock sync.Mutex

func GetMap(name string) map[string]string {
	mapsLock.Lock()
	defer mapsLock.Unlock()

	return _unsafe_GetMap(name)
}

func _unsafe_GetMap(name string) map[string]string {
	bytes, err := ioutil.ReadFile(MapFilename(name))
	if err != nil {
		return make(map[string]string)
	}

	lines := strings.Split(string(bytes), "\n")
	if len(lines) > 0 && len(lines[len(lines)-1]) == 0 {
		// remove empty string after trailing endline
		lines = lines[:len(lines)-1]
	}

	m := make(map[string]string)
	for _, l := range lines {
		parts := strings.SplitN(l, "=", 2)
		m[parts[0]] = parts[1]
	}

	return m
}

func MapFilename(name string) string {
	return viper.GetString("datadir") + "/" + name + ".map"
}

func _unsafe_WriteMap(name string, m map[string]string) {
	filename := MapFilename(name)
	contents := ""
	for key, value := range m {
		contents = contents + key + "=" + value + "\n"
	}

	err := ioutil.WriteFile(filename, []byte(contents), 0644)
	if err != nil {
		fmt.Printf("error writing to %s: %v", filename, err)
	}
}

// XXX: don't use a key that has an "=" sign in it; the "=" is used as a separator and is not escaped
func AddToMap(name string, key string, value string) {
	mapsLock.Lock()
	defer mapsLock.Unlock()

	m := _unsafe_GetMap(name)
	m[key] = value
	_unsafe_WriteMap(name, m)
}

// return the value (or "" if none), and true if there was a value and false if not
func RetrieveFromMap(name string, key string) (string, bool) {
	mapsLock.Lock()
	defer mapsLock.Unlock()

	m := _unsafe_GetMap(name)
	if value, ok := m[key]; ok {
		return value, true
	} else {
		return "", false
	}
}

// return true if it already existed and false if it didn't
func RemoveFromMap(name string, key string) bool {
	mapsLock.Lock()
	defer mapsLock.Unlock()

	m := _unsafe_GetMap(name)
	_, existed := m[key]
	delete(m, key)
	_unsafe_WriteMap(name, m)

	return existed
}
