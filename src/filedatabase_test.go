package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeMarshallable struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
	Blop SubStruct
}

type SubStruct struct {
	Blop string
}

func TestSet(t *testing.T) {

	os.WriteFile("../var/dbfile/dbtest", []byte{}, 0777)

	db := FileDatabase[FakeMarshallable]{
		filename: "dbtest",
	}

	db.Set("123", FakeMarshallable{"Michel", 42, SubStruct{"Yop"}})
	db.Set("1234", FakeMarshallable{"Miguelinhos", 43, SubStruct{}})
	db.Set("12345", FakeMarshallable{"Micheline", 8, SubStruct{"Yopla"}})
	db.Set("abvcfd", FakeMarshallable{"Mickel", 78, SubStruct{}})
}

func TestGet(t *testing.T) {

	os.WriteFile("../var/dbfile/dbtest", []byte{}, 0777)

	assert := assert.New(t)

	db := FileDatabase[FakeMarshallable]{
		filename: "dbtest",
	}

	db.Set("123", FakeMarshallable{"Michel", 42, SubStruct{"Yop"}})
	db.Set("1234", FakeMarshallable{"Miguelinhos", 43, SubStruct{"Yap"}})

	f1 := db.Get("123")
	f2 := db.Get("1234")

	assert.Equal("Michel", f1.Name)
	assert.Equal(42, f1.Age)
	assert.Equal("Yop", f1.Blop.Blop)
	assert.Equal("Miguelinhos", f2.Name)
	assert.Equal(43, f2.Age)
	assert.Equal("Yap", f2.Blop.Blop)
}

func TestExist(t *testing.T) {

	os.WriteFile("../var/dbfile/dbtest", []byte{}, 0777)

	assert := assert.New(t)

	db := FileDatabase[FakeMarshallable]{
		filename: "dbtest",
	}

	db.Set("123", FakeMarshallable{"Michel", 42, SubStruct{"Yop"}})
	db.Set("1234", FakeMarshallable{"Miguelinhos", 43, SubStruct{"Yap"}})

	assert.True(db.Exist("123"))
	assert.True(db.Exist("1234"))
	assert.False(db.Exist("coucou"))
	assert.False(db.Exist(""))
}
