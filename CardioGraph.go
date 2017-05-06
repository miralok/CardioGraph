package main

import (
	"errors"
	"github.com/hyperledger/fabric/core/chaincode/shim"
//	"strings"
	"strconv"
)

const (
	CARDIO_GRAPH_VALUE_NONE = "NONE"
)

type CardioGraph struct {
	Entity
	PatientID	string //unique
	Name	string
	Gender	string
	CreationDateTime	string //unique
	Age	int32
	Height	int32
	Weight	int32
	HeartRate int32
	PPInterval int32
	STInterval	int32
	QRSDuration	int32
	ErrMsg	string
}

type CardioGraphCompare func(CardioGraph, CardioGraph) (bool)

func createCardioGraph(stub shim.ChaincodeStubInterface) CardioGraph {
	cg := CardioGraph{}
	e := Entity{Stub: stub, TableName: "CardioGraph"}

	cg.Entity = e
	return cg
}

func (cg CardioGraph) createTable() error {
	return createTable(cg.Entity, cg)
}

func (cg CardioGraph) deleteTable() error {
	return deleteTable(cg.Entity)
}

func (cg CardioGraph) insert() error {
	return insert(cg.Entity, cg)
}

func (cg CardioGraph) replace() error {
	return replace(cg.Entity, cg)
}

func (cg CardioGraph)delete()(error){
	return delete(cg.Entity, cg)
}

func (cg CardioGraph) get() (CardioGraph, error) {
	row, err := get(cg.Entity, cg)

	return cg.Convert(row), err
}

func (cg CardioGraph) InitForCreate() []*shim.ColumnDefinition {
	return []*shim.ColumnDefinition{
		{Name: "PatientID", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Name", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Age", Type: shim.ColumnDefinition_INT32, Key: false},
		{Name: "Gender", Type: shim.ColumnDefinition_STRING, Key: true},
		{Name: "Height", Type: shim.ColumnDefinition_INT32, Key: false},
		{Name: "Weight", Type: shim.ColumnDefinition_INT32, Key: false},
		{Name: "HeartRate", Type: shim.ColumnDefinition_INT32, Key: false},
		{Name: "PPInterval", Type: shim.ColumnDefinition_INT32, Key: false},
		{Name: "STInterval", Type: shim.ColumnDefinition_INT32, Key: false},
		{Name: "QRSDuration", Type: shim.ColumnDefinition_INT32, Key: false},
		{Name: "CreationDateTime", Type: shim.ColumnDefinition_STRING, Key: true},
	}
}

func (cg CardioGraph) InitForInsertAndReplace() []*shim.Column {
	return []*shim.Column{
		{Value: &shim.Column_String_{String_: cg.PatientID}},
		{Value: &shim.Column_String_{String_: cg.Name}},
		{Value: &shim.Column_Int32{Int32: cg.Age}},
		{Value: &shim.Column_String_{String_: cg.Gender}},
		{Value: &shim.Column_Int32{Int32: cg.Height}},
		{Value: &shim.Column_Int32{Int32: cg.Weight}},
		{Value: &shim.Column_Int32{Int32: cg.HeartRate}},
		{Value: &shim.Column_Int32{Int32: cg.PPInterval}},
		{Value: &shim.Column_Int32{Int32: cg.STInterval}},
		{Value: &shim.Column_Int32{Int32: cg.QRSDuration}},
		{Value: &shim.Column_String_{String_: cg.CreationDateTime}},
	}
}

func (cg CardioGraph) InitForDeleteAndGet() []shim.Column {
	return []shim.Column{
		{Value: &shim.Column_String_{String_: cg.PatientID}},
		{Value: &shim.Column_String_{String_: cg.CreationDateTime}},
	}
}

func (cg CardioGraph) Convert(row shim.Row) CardioGraph {
	cg.PatientID = row.Columns[0].GetString_()
	cg.Name = row.Columns[1].GetString_()
	cg.Age = row.Columns[2].GetInt32()
	cg.Gender = row.Columns[3].GetString_()
	cg.Height = row.Columns[4].GetInt32()
	cg.Weight = row.Columns[5].GetInt32()
	cg.HeartRate = row.Columns[6].GetInt32()
	cg.PPInterval = row.Columns[7].GetInt32()
	cg.STInterval = row.Columns[8].GetInt32()
	cg.QRSDuration = row.Columns[9].GetInt32()
	cg.CreationDateTime = row.Columns[10].GetString_()
	return cg
}

func (cg CardioGraph)generateKeyArray(num int8)([]shim.Column) {
	var col shim.Column

	switch(num){
	case 0:	return []shim.Column{}
	case 1:	col = shim.Column{Value: &shim.Column_String_{String_: cg.PatientID}}
	case 2:	col = shim.Column{Value: &shim.Column_String_{String_: cg.CreationDateTime}}
	}

	return append(cg.generateKeyArray(num-1), col)
}

func (cg CardioGraph)getAll(keyArr []shim.Column, cgCompare CardioGraphCompare)([]CardioGraph, error) {
	var err error
	var rowChan <-chan shim.Row
	var cgArr []CardioGraph

	rowChan, err = getAll(cg.Entity, cg, keyArr)

	if err != nil {
		return cgArr, err
	}

	for row := range rowChan {
		cgTemp :=  createCardioGraph(cg.Stub)
		cgTemp = cgTemp.Convert(row)
		if cgCompare(cg, cgTemp) {
			cgArr = append(cgArr, cgTemp)
		}
	}

	return cgArr, nil
}

func (cg CardioGraph)getAllByAge()([]CardioGraph, error){
	var keyArr []shim.Column
	var cgCompare CardioGraphCompare

	keyArr = cg.generateKeyArray(0)

	cgCompare = func(cg1 CardioGraph, cg2 CardioGraph) (bool) {
		return cg1.Age == cg2.Age
	}

	cgArr, err := cg.getAll(keyArr, cgCompare)
	
	return cgArr, err
}

func (cg CardioGraph)getAllByPatientID()([]CardioGraph, error){
	var keyArr []shim.Column
	var cgCompare CardioGraphCompare

	keyArr = cg.generateKeyArray(0)

	cgCompare = func(cg1 CardioGraph, cg2 CardioGraph) (bool) {
		return cg1.PatientID == cg2.PatientID;
	}

	cgArr, err := cg.getAll(keyArr, cgCompare)

	return cgArr, err
}


/*******************************

CHAINCODE HELPER FUNCTIONS

 *******************************/




func insertCardioGraph(stub shim.ChaincodeStubInterface, args []string) (CardioGraph, error) {
	cg := createCardioGraph(stub)

	if len(args) != 11 {
		return cg, errors.New("Function insertCardioGraph expects 11 arguments.")
	}

	patientID := args[0];
	name := args[1];
	age, _ := strconv.Atoi(args[2])
	gender := args[3];
	height, _ := strconv.Atoi(args[4])
	weight, _ := strconv.Atoi(args[5])
	heartRate, _ := strconv.Atoi(args[6])
	ppInterval, _ := strconv.Atoi(args[7])
	stInterval, _ := strconv.Atoi(args[8])
	qrsDuration, _ := strconv.Atoi(args[9])
	creationDateTime := args[10];
	
	cg.PatientID = patientID
	cg.Name = name
	cg.Age = int32(age)
	cg.Gender = gender
	cg.Height = int32(height)
	cg.Weight = int32(weight)
	cg.HeartRate = int32(heartRate)
	cg.PPInterval = int32(ppInterval)
	cg.STInterval = int32(stInterval)
	cg.QRSDuration = int32(qrsDuration)
	cg.CreationDateTime = creationDateTime

	err := cg.insert()
	if err != nil {
		cg.get()
		return cg, err
	}
	return cg, nil
}

func deleteCardioGraph(stub shim.ChaincodeStubInterface, args []string) (cgRet CardioGraph, errRet error) {

	cg := createCardioGraph(stub)


	if len(args) != 1 {
		return cg, errors.New("Function deleteCardioGraph expects 1 argument.")
	}

	patientID := args[0]

	cg.PatientID = patientID

	err := cg.delete()

	if err != nil {
		return cg, err
	}

	return cg, err
}


func getCardioGraph(stub shim.ChaincodeStubInterface, args []string) (cgRet CardioGraph, errRet error) {

	cg := createCardioGraph(stub)

	if len(args) != 1 {
		return cg, errors.New("Function getCardioGraph expects 1 argument.")
	}

	patientID := args[0]

	cg.PatientID = patientID
	cg, err := cg.get()

	if err != nil {
		return cg, err
	}

	return cg, err
}

func getAllCardioGraphByAge(stub shim.ChaincodeStubInterface, args []string) ([]CardioGraph, error){
	var err error
	var cgArr []CardioGraph

	cg := createCardioGraph(stub)

	if len(args) != 1 {
		return cgArr, errors.New("Function getAllCardioGraph expects 1 argument.")
	}

	age, _ := strconv.Atoi(args[1])

	cg.Age = int32(age)
	cgArr, err = cg.getAllByAge()

	return cgArr, err
}

func getCardioGraphByPatientID(stub shim.ChaincodeStubInterface, args []string) ([]CardioGraph, error){
	var err error
	var cgArr []CardioGraph

	cg := createCardioGraph(stub)

	if len(args) != 1 {
		return cgArr, errors.New("Function getAllCardioGraph expects 1 argument.")
	}

	patientID := args[0]

	cg.PatientID = patientID
	cgArr, err = cg.getAllByPatientID()

	return cgArr, err
}