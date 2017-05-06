package main

import (
	"errors"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type CardioGraphChaincode struct {
}


type Attribute struct {
	AttributeName string
	AttributeVal string
	ErrMsg string
}

/*
args:
[0] - Attribute Name
 */
func (*CardioGraphChaincode) getAttribute(stub shim.ChaincodeStubInterface, args []string) (Attribute, error){
	var attr Attribute
	var err error

	if len(args) != 1 {
		return attr, errors.New("Function getAttribute expects 1 argument.")
	}


	index := -1

	index++
	attributeName := formatInput(args[index])

	attributeVal, err := getCertAttribute(stub, attributeName)
	if (err != nil){
		err = errors.New("getAttribute cannot get the attribute ["+ attributeName + "]")

		return attr, err
	}

	//attr  = Attribute{AttributeName: attributeName, AttributeVal: attributeVal}
	attr = Attribute{}
	attr.AttributeName = attributeName
	attr.AttributeVal = attributeVal
	return attr, nil
}

func (*CardioGraphChaincode) deleteTable(stub shim.ChaincodeStubInterface) (error){
	cg := createCardioGraph(stub)

	_ = cg.deleteTable()

	return nil
}

func (*CardioGraphChaincode) createTable(stub shim.ChaincodeStubInterface) (error){
	var err error

	cg := createCardioGraph(stub)

	err = cg.createTable()
	if err != nil{
		return err
	}

	return nil
}

/*
Init is called when chaincode is deployed.
*/
func (cc *CardioGraphChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	if function == "createTable" {

		err = cc.createTable(stub)

		return nil, err

	}

	return nil, errors.New("Unknown function " + function + ".")
}

func (cc *CardioGraphChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	var err error

	if function == "getAttribute" {
		var attr Attribute

		attr, err = cc.getAttribute(stub, args)

		if err != nil {
			attr.ErrMsg = err.Error()
			return formatOutput(attr)
		}

		return formatOutput(attr)
	}else
	
	if function == "getCardioGraph" { 
		var cg CardioGraph

		cg, err = getCardioGraph(stub, args)

		if err != nil {
			cg.ErrMsg = err.Error()
			return formatOutput(cg)
		}

		//jsonObj, err := json.Marshal(cg)

		//return []byte(string(jsonObj)), err
		return formatOutput(cg)
	}else
	if function == "getAllCardioGraphByAge" {
		var cgArr []CardioGraph

		cgArr, err = getAllCardioGraphByAge(stub, args)

		if err != nil {
			return nil, err
		}

		//jsonObj, err := json.Marshal(icArr)

		return formatOutput(cgArr)
	}else
	if function == "getAllCardioGraphByName" {
		var cgArr []CardioGraph

		cgArr, err = getAllCardioGraphByAge(stub, args)

		if err != nil {
			return nil, err
		}

		//jsonObj, err := json.Marshal(icArr)

		return formatOutput(cgArr)
	}
	
	return nil, errors.New("Unknown function " + function + ".")
}

func (cc *CardioGraphChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	var err error
	if function == "deleteTable" {
		err = cc.deleteTable(stub)

		return nil, err
	}else
	if function == "insertCardioGraph" {
		var cg CardioGraph

		cg, err = insertCardioGraph(stub, args)

		if err != nil {
			cg.ErrMsg = err.Error();
			stub.SetEvent("insertCardioGraph", formatPayload(cg))
			return nil, err
		}

		err = stub.SetEvent("insertCardioGraph", formatPayload(cg))
		if err != nil {
			return nil, err
		}
		return nil, nil

	}else
	if function == "deleteCardioGraph" {
		var cg CardioGraph

		cg, err = deleteCardioGraph(stub, args)

		if err != nil {
			cg.ErrMsg = err.Error();
			stub.SetEvent("deleteCardioGraph", formatPayload(cg))
			return nil, err
		}

		err = stub.SetEvent("deleteCardioGraph", formatPayload(cg))
		if err != nil {
			return nil, err
		}
		return nil, nil

	}

	return nil, nil
}

func main() {
	err := shim.Start(new(CardioGraphChaincode))
	if err != nil {
		fmt.Printf("Error creationing CardioGraphChaincode: %s", err)
	}
}