package main

import (
	"container/list"
	"fmt"
	"strconv"
	"testing"
	// "container/list"
	// "errors"
	// "fmt"
	// "strconv"
)

func TestCar(t *testing.T) {
	car := Car{uint64(1), "Wolfswagen", uint8(4), uint16(2008)}

	t.Run("GetColValueByColIndex", func(t *testing.T) {
		value, _ := car.GetColValueByColIndex(0)

		if value.(uint64) == uint64(1) {

		} else {
			t.Errorf("got %d, want %d", value, 1)
		}
	})

	t.Run("GetColValueByColName", func(t *testing.T) {
		value, _ := car.GetColValueByColName("id")

		if value.(uint64) == uint64(1) {

		} else {
			t.Errorf("got %d, want %d", value, 1)
		}
	})

	t.Run("GetStringRow", func(t *testing.T) {
		str := "1 Wolfswagen 4 2008"
		value := car.GetStringRow()

		if value == str {

		} else {
			t.Errorf("got %s, want %s", value, str)
		}
	})
}

func TestTable(t *testing.T) {

	carTable := Table{
		colNames: carColNames,
		rowList:  list.New(),
	}

	var outputInsert [99]Car

	for i := 1; i < 100; i++ {
		car := Car{
			uint64(i), "Car" + strconv.FormatUint(uint64(i), 10), 4, uint16(2000 + i),
		}
		carTable.Insert(car)
		outputInsert[i-1] = car
	}

	t.Run("Insert test", func(t *testing.T) {
		fmt.Println("Insert test start")
		for i, el := 0, carTable.rowList.Front(); i < 99; i++ {

			if el.Value.(Row).GetStringRow() != outputInsert[i].GetStringRow() {
				t.Errorf("get %s, want %s, interation %d", el.Value.(Row).GetStringRow(),
				 outputInsert[i].GetStringRow(), i+1)
			}
			el = el.Next()

		}
		fmt.Println("Insert test end")
	})

	var findElemetsTest [99]struct {
		colName   string
		value     any
		wantedRow string
	}

	for i := range findElemetsTest {
		colIndex := 0
		findElemetsTest[i].colName = carTable.colNames[colIndex]
		findElemetsTest[i].value, _ = outputInsert[i].GetColValueByColIndex(uint(colIndex))
		findElemetsTest[i].wantedRow = outputInsert[i].GetStringRow()
	}

	t.Run("FindElemets test", func(t *testing.T) {
		fmt.Println("FindElemets test start")
		for _, el := range findElemetsTest {
			//fmt.Println("it:", i)
			elFind, _ := carTable.findElemets(el.colName, el.value)
			if el.wantedRow != elFind.Value.(Row).GetStringRow() {

				t.Errorf("get %s, want %s", elFind.Value.(Row).GetStringRow(), el.wantedRow)
			}

		}
		fmt.Println("FindElemets test end")
	})

	var deleteElemetsTest [20]struct {
		colName   string
		value     any
		err       error
		wantedRow string
	}

	for i := 0; i < 20; i++ {
		colIndex := 0
		deleteElemetsTest[i].colName = carTable.colNames[colIndex]
		deleteElemetsTest[i].value = uint64(i + 1)
		deleteElemetsTest[i].err = nil
		deleteElemetsTest[i].wantedRow = outputInsert[i].GetStringRow()
	}

	t.Run("Delete test", func(t *testing.T) {
		fmt.Println("Delete test start")
		for _, el := range deleteElemetsTest {
			err := carTable.Delete(el.colName, el.value)

			if err != el.err {
				t.Errorf("get err %s, want err %s", err.Error(), el.err.Error())
			}
		}
		fmt.Println("Delete test end")
	})

	var updateElemetsTest [20]struct {
		colName   string
		oldValue  any
		newRow    Car
		wantedRow string
	}

	for i := 0; i < 20; i++ {
		colIndex := 0
		updateElemetsTest[i].colName = carTable.colNames[colIndex]
		updateElemetsTest[i].oldValue = uint64(i + 21)
		updateElemetsTest[i].newRow = Car{uint64(i + 100), "Toyta", uint8(4), uint16(2010)}
		updateElemetsTest[i].wantedRow = outputInsert[i].GetStringRow()
	}
	t.Run("Update test", func(t *testing.T) {
		fmt.Println("Update test start")
		for _, el := range updateElemetsTest {

			err := carTable.Update(el.colName, el.oldValue, el.newRow)

			rowValue, _ := el.newRow.GetColValueByColIndex(0)
			findRow, _ := carTable.findElemets(el.colName, rowValue)
			findValue, _ := findRow.Value.(Row).GetColValueByColIndex(0)
			if findValue != rowValue {
				t.Errorf("get %s, want %s", findValue, rowValue)
			}

			if err != nil {
				t.Errorf("get err %s, want err %s", err.Error(), "nil")
			}
		}
		fmt.Println("Update test end")
	})

	cleanCarTable := Table{
		colNames: carColNames,
		rowList:  list.New(),
	}

	for i := 1; i < 100; i++ {
		car := Car{
			uint64(i), "Car" + strconv.FormatUint(uint64(i), 10), 4, uint16(2000 + i),
		}
		cleanCarTable.Insert(car)
		outputInsert[i-1] = car
	}

	var getRowElemetsTest [99]struct {
		colName         string
		value           any
		wantedRowString string
	}

	for i := 0; i < 99; i++ {
		colIndex := 0
		getRowElemetsTest[i].colName = carTable.colNames[colIndex]
		getRowElemetsTest[i].value = uint64(i + 1)
		getRowElemetsTest[i].wantedRowString = outputInsert[i].GetStringRow()
	}
	t.Run("GetRow test", func(t *testing.T) {
		fmt.Println("GetRow test start")
		for i, el := range getRowElemetsTest {
			row, err := cleanCarTable.GetRow(el.colName, el.value)

			if err != nil {
				t.Errorf("get error: %s, want: %s", err, el.wantedRowString)
			}

			if row.GetStringRow() != outputInsert[i].GetStringRow() {
				t.Errorf("get %s, want %s", row.GetStringRow(), 
				outputInsert[i].GetStringRow())
			}

		}
		fmt.Println("GetRow test end")
	})
}

func asdasd(){

	var findElemetsTest [99]struct {
		colName   string
		value     any
		wantedRow string
	}

	var deleteElemetsTest [20]struct {
		colName   string
		value     any
		err       error
		wantedRow string
	}

	var updateElemetsTest [20]struct {
		colName   string
		oldValue  any
		newRow    Car
		wantedRow string
	}

	fmt.Print(findElemetsTest, deleteElemetsTest, updateElemetsTest)
}