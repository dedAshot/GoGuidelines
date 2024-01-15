package main

import (
	"container/list"
	"errors"
	"fmt"
	"strconv"
)

type Row interface {
	GetStringRow() string
	GetColValueByColIndex(uint) (any, error)
	GetColValueByColName(string) (any, error)
}

type Table struct {
	colNames []string
	rowList  *list.List
}

func (Table *Table) findElemets(colName string, value any) (*list.Element, error) {
	var err error = errors.New("no such element")

	for e := Table.rowList.Front(); e != nil; e = e.Next() {

		currentValue, _ := e.Value.(Row).GetColValueByColName(colName)

		if currentValue == value {
			err = nil
			return e, err
		}
	}

	return nil, err
}

func (Table *Table) Insert(Row Row) {

	Table.rowList.PushBack(Row)
	fmt.Println("inserted: ", Row)
}

func (Table *Table) Delete(colName string, value any) error {
	var err error = nil
	var element *list.Element


	element, err = Table.findElemets(colName, value)
	if err != nil {
		return err
	}

	removedElem := Table.rowList.Remove(element)
	fmt.Println("removed: ", removedElem)
	return err
}

func (Table *Table) Update(colName string, oldValue any, newRow Row) error {
	var err error = nil
	var oldRow *list.Element

	oldRow, err = Table.findElemets(colName, oldValue)

	if err != nil {
		return err
	}

	Table.rowList.InsertBefore(newRow, oldRow)
	removedElem := Table.rowList.Remove(oldRow)

	fmt.Println("updated: ", removedElem, " -> ", newRow)

	return err
}

func (Table *Table) GetRow(colName string, value any) (Row, error) {
	var err error = nil

	rowEl, err := (Table).findElemets(colName, value)

	if err != nil {
		return nil, err
	}

	row := rowEl.Value.(Row)

	fmt.Println("get row: ", row)

	return row, err
}

type Car struct {
	id         uint64
	model      string
	seatsCount uint8
	buildYear  uint16
}

var carColNames []string

func (car Car) GetStringRow() string {
	return strconv.FormatUint(uint64(car.id), 10) + " " + car.model + " " +
		strconv.FormatUint(uint64(car.seatsCount), 10) + " " +
		strconv.FormatUint(uint64(car.buildYear), 10)
}

func (car Car) GetColValueByColIndex(index uint) (any, error) {
	var err error = nil

	switch index {
	case 0:
		return car.id, err
	case 1:
		return car.model, err
	case 2:
		return car.seatsCount, err
	case 3:
		return car.buildYear, err

	}

	err = errors.New("too large index (>3)")
	return nil, err
}

func (car Car) GetColValueByColName(colName string) (any, error) {
	var err error = nil
	//fmt.Println("GetColValueByColName")

	switch colName {
	case carColNames[0]:
		return car.id, err
	case carColNames[1]:
		return car.model, err
	case carColNames[2]:
		return car.seatsCount, err
	case carColNames[3]:
		return car.buildYear, err
	}

	err = errors.New("no such column")
	return nil, err
}

func init() {
	carColNames = []string{
		"id",
		"model",
		"seatsCount",
		"buildYear",
	}
}

func main() {
	fmt.Println("Hello, World!")
	defer fmt.Println("Goodbye World")

	var carTable Table = Table{
		colNames: carColNames,
		rowList:  list.New(),
	}

	var row Car = Car{1, "Toyta", 4, 2000}

	carTable.Insert(row)

	carTable.Delete(carTable.colNames[0], uint64(1))

	fmt.Println(carTable.rowList.Len())

	for i := 1; i < 100; i++ {
		carTable.Insert(Car{
			uint64(i), "Car" + strconv.FormatUint(uint64(i), 10), 4, uint16(2000 + i),
		})
	}

	el, _ := carTable.findElemets("id", uint64(5))
	fmt.Println((el).Value.(Row).GetStringRow())
	fmt.Println(carTable.rowList.Back())
}
