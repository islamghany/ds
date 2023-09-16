package main

import (
	"fmt"
	hashfunctions "hashing/hashFunctions"
	hashtables "hashing/hashTables"
)

func print(x ...interface{}) {
	fmt.Println(x)
}

type MapItem string

func main() {
	map1 := hashtables.NewSeperateChaining[MapItem](10, hashfunctions.NewFolding())
	map1.Insert("1", "one")
	map1.Insert("2", "two")
	map1.Insert("3", "three")
	print(map1.Size())
	item, ok := map1.Search("5")
	print(item, ok)
	item, ok = map1.Search("1")
	print(*item, ok)
	item, ok = map1.Search("2")
	print(*item, ok)
	item, ok = map1.Search("3")
	print(*item, ok)
	map1.Delete("2")
	item, ok = map1.Search("2")
	print(item, ok)

}
