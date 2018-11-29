package main

import (
	"fmt"
	"time"

	"github.com/ContinuumLLC/platform-common-lib/src/db/sqlite"
)

type TestTable struct {
	MessageID    string `gorm:"primary_key"`
	URL          string
	TimeStampUTC time.Time
	IsProcessed  bool
	Status       int
}

func main() {
	service := sqlite.GetService(&sqlite.Config{DBName: "Test.db"})
	err := service.Init()
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	err = service.CreateTable(&TestTable{})
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	err = service.Add(&TestTable{
		MessageID:    "11",
		URL:          "Test",
		TimeStampUTC: time.Now().UTC(),
		IsProcessed:  false,
		Status:       1,
	})
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	err = service.Update(&TestTable{
		MessageID:    "11",
		URL:          "Test",
		TimeStampUTC: time.Now().UTC(),
		IsProcessed:  false,
		Status:       11,
	})
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	data := &TestTable{}
	err = service.Get(1, data)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	fmt.Println(data)

	err = service.Delete(&TestTable{
		MessageID:    "11",
		URL:          "Test",
		TimeStampUTC: time.Now().UTC(),
		IsProcessed:  false,
		Status:       11,
	})
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	components := []TestTable{
		TestTable{
			MessageID:    "23",
			URL:          "Test",
			TimeStampUTC: time.Now().UTC(),
			IsProcessed:  false,
			Status:       1,
		}, TestTable{
			MessageID:    "24",
			URL:          "Test",
			TimeStampUTC: time.Now().UTC(),
			IsProcessed:  false,
			Status:       1,
		},
	}

	err = service.AddAll(a(components))
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}

	var data1 []TestTable
	err = service.Get(10, &data1)
	if err != nil {
		fmt.Printf("%+v", err)
		return
	}
	fmt.Println(data1)

	service.Close()
}

func a(c []TestTable) []interface{} {
	cmp := make([]interface{}, len(c))
	for index := 0; index < len(c); index++ {
		cmp[index] = c[index]
	}
	return cmp
}
