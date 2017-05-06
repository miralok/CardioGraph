package main

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"errors"
	"fmt"
)

type Entity struct{
	Stub shim.ChaincodeStubInterface
	TableName string
}

type EntityInitializer interface{
	InitForCreate() []*shim.ColumnDefinition
	InitForInsertAndReplace() []*shim.Column
	InitForDeleteAndGet() []shim.Column
}

func createTable(e Entity, ei EntityInitializer) (error) {
	var err error

	deleteTable(e)

	err = e.Stub.CreateTable(e.TableName, ei.InitForCreate())

	if err != nil {
		return fmt.Errorf("Failed creating "+e.TableName+" table, [%v]", err)
	}

	return nil
}

func deleteTable(e Entity) (error) {
	var err error

	err = e.Stub.DeleteTable(e.TableName)

	return err
}


func insert(e Entity, ei EntityInitializer)(error){
	var err error
	var ok bool

	ok, err = e.Stub.InsertRow(e.TableName, shim.Row{Columns: ei.InitForInsertAndReplace(),})

	if !ok && err == nil {
		return errors.New("Row already exists in "+e.TableName+" table.")
	}else
	if err != nil {
		return errors.New("Error inserting row in "+e.TableName+" table.")
	}

	return nil
}

func replace(e Entity, ei EntityInitializer)(error){
	var err error
	var ok bool

	ok, err = e.Stub.ReplaceRow(e.TableName, shim.Row{Columns: ei.InitForInsertAndReplace(),})

	if !ok && err == nil {
		return errors.New("Row does not exist in "+e.TableName+" table.")
	}else
	if err != nil {
		return errors.New("Error updating row in "+e.TableName+" table.")
	}

	return nil
}

func get(e Entity, ei EntityInitializer)(shim.Row, error){
	var err error
	var row shim.Row

	row, err = e.Stub.GetRow(e.TableName, ei.InitForDeleteAndGet(),)

	if err != nil {
		return row, errors.New("Error getting row from "+e.TableName+" table.")
	}

	return row, nil
}

func getAll(e Entity, ei EntityInitializer, keyArr []shim.Column)(<-chan shim.Row, error){
	var err error
	var rowChan <-chan shim.Row

	rowChan, err = e.Stub.GetRows(e.TableName, keyArr)

	if err != nil {
		return rowChan, errors.New("Error getting row from "+e.TableName+" table.")
	}

	return rowChan, nil
}

func delete(e Entity, ei EntityInitializer)(error){
	var err error

	err = e.Stub.DeleteRow(e.TableName, ei.InitForDeleteAndGet(),)

	if err != nil {
		return errors.New("Error deleting row from "+e.TableName+" table.")
	}

	return nil
}

