/*
query data from Blockchain since the data insert to chain is immutable
invoke: insert audit data to Blockchain and save all the aduit data to mongo
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

type AuditTrailChaincode struct {
}


type Audit struct {
	AuditHash string `json:"audit_hash"` //Audit identifier
	BusinessKey  string    `json:"business_key"`
	UpdatedBy string `json:"updated_by"`
	//return current time in milliseconds as a string
}

func (t *AuditTrailChaincode) createAudit(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	/*
	Json{
		"audit_hash"		: "SHA hash code",
		"business_key"		: "191566",
		"updated_by"		: "lw84456"
	}
	*/
	fmt.Println("Creating audit")

	if len(args) != 1 {
		fmt.Println("Error obtaining username")
		return nil, errors.New("createAudit accepts a single username argument")
	}

	//Build an audit object
	var audit Audit
	var err error

	fmt.Println("Unmarshalling Audit")
	err = json.Unmarshal([]byte(args[0]), &audit)
	if err != nil {
		fmt.Println("error invalid audit")
		return nil, errors.New("Invalid commercial audit")
	}

	auditBytes, err := json.Marshal(&audit)
	if err != nil {
		fmt.Println("Error marshalling audit")
		return nil, errors.New("Error adding new audit")
	}

	fmt.Println("Attempting to get state of any existing audit for " + audit.AuditHash)
	existingBytes, err := stub.GetState(audit.AuditHash)
	if err == nil {
		fmt.Println("Audit does not exist, creating it")
	} else {
		fmt.Println("Exist audit will be overrided for hash: " + audit.AuditHash)
	}

	err = stub.PutState(audit.AuditHash, auditBytes)
	if err != nil {
		fmt.Println("Error save audit to chain")
		return nil, errors.New("Error saving audit to chain")
	}

	fmt.Println("Creat audit %+v\n", audit)
	return nil, nil
}

func (t *AuditTrailChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Init firing. Function will be ignored: " + function)

	// Initialize the collection of AuditTrails
	fmt.Println("Initializing audit keys collection")
	var blank []string
	blankBytes, _ := json.Marshal(&blank)
	err := stub.PutState("AuditKeys", blankBytes)
	if err != nil {
		fmt.Println("Failed to initialize audit key collection")
	}

	fmt.Println("Initialization complete")
	return nil, nil
}

func IsValid(audit_hash string, stub shim.ChaincodeStubInterface) (bool, error) {
	auditBytes, err := stub.GetState(audit_hash)
	if err != nil {
		fmt.Println("Error retrieving audit " + audit_hash)
		return false, errors.New("Error retrieving cp " + audit_hash)
	}

	if auditBytes == nil {
		fmt.Println("Not exists data by searching chain with hash: " + audit_hash)
		return false, nil
	}

	return true, nil
}

func (t *AuditTrailChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query running. Function: " + function)

	if function == "IsValid" {
		fmt.Println("Validate hash whether existing")
		isHashValid, err1 :=IsValid(args[0], stub)

		if err1 != nil {
			fmt.Println("Error Validating")
			return nil, err1
		}
		return []byte(isHashValid), nil

	} else {
		fmt.Println("Generic Query call")
		bytes, err := stub.GetState(args[0]) //no function name provided, will search chain by key

		if err != nil {
			fmt.Println("Some error happenend: " + err.Error())
			return nil, err
		}

		fmt.Println("All success, returning from generic")
		return bytes, nil
	}
}

func (t *AuditTrailChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke running. Function: " + function)

	if function == "createAudit" {
		return t.createAudit(stub, args)
	}


	return nil, errors.New("Received unknown function invocation: " + function)
}

func main() {
	err := shim.Start(new(AuditTrailChaincode))
	if err != nil {
		fmt.Println("Error starting Simple chaincode: %s", err)
	}
}
