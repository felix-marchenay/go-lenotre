package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type FileDatabase[T any] struct {
	filename  string
	values    map[string]T
	initiated bool
}

func (f *FileDatabase[T]) init() {
	if f.initiated {
		return
	}

	data := getData(f.filename)

	if len(data) == 0 {
		f.values = make(map[string]T)
		f.initiated = true
		return
	}

	err := json.Unmarshal(data, &f.values)
	if err != nil {
		panic(err)
	}

	f.initiated = true
}

func (f FileDatabase[T]) Exist(key string) bool {
	f.init()

	_, exist := f.values[key]

	return exist
}

func (f FileDatabase[T]) Get(key string) T {
	f.init()

	if f.Exist(key) == false {
		panic("Aucune entrée pour la clef " + key)
	}

	return f.values[key]
}

func (f *FileDatabase[T]) Set(key string, value T) {
	f.init()

	f.values[key] = value

	json, err := json.Marshal(f.values)
	if err != nil {
		panic(err)
	}

	setData(f.filename, json)
}

func getData(filename string) []byte {
	f, err := os.ReadFile("../var/dbfile/" + filename)

	if err == nil {
		os.WriteFile(filename, []byte{}, 0777)
		return []byte{}
	}

	return f
}

func setData(filename string, data []byte) {
	err := os.WriteFile("../var/dbfile/"+filename, data, 0777)
	if err != nil {
		fmt.Println(err)
	}
}
