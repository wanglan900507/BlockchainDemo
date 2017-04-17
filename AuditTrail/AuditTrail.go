/*
query data from Blockchain since the data insert to chain is immutable
invoke: insert audit data to Blockchain and save all the aduit data to mongo
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"strings"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"hash"
)

type AuditTrailChaincode struct {
}

const (
	millisPerSecond = int64(time.Second / time.Millisecond)
	nanosPerMillisecond = int64(time.Millisecond / time.Nanosecond)
)

func msToTime(ms string) (time.Time, error) {
	msInt, err := strconv.ParseInt(ms, 10, 64)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(msInt / millisPerSecond,
		(msInt % millisPerSecond) * nanosPerMillisecond), nil
}

type Audit struct {
	AuditHash hash.Hash64 `json:"audit_hash"` //Audit identifier
	BusinessKey  float64    `json:"business_key"`
	Version float64 `json:"audit_version"'`
	UpdatedBy string `json:"updated_by"`
	UpdatedTimestamp time.Time `json:"updated_timestamp"`
}

func (t *AuditTrailChaincode) createAudit(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	/*
	Json{
		"audit_hash"		: "SHA hash code",
		"business_key"		: 191566,
		"audit_version"		: 10,
		"updated_by"		: "lw84456"
		"updated_timestamp"	: "

	}
	*/


	fmt.Println("Creating audit")

	if len(args) != 1 {
		fmt.Println("Error obtaining username")
		return nil, errors.New("createAudit accepts a single username argument")
	}
}

func (t *AuditTrailChaincode) createAccount(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating account")

	// Obtain the username to associate with the account
	if len(args) != 1 {
		fmt.Println("Error obtaining username")
		return nil, errors.New("createAccount accepts a single username argument")
	}
	username := args[0]

	// Build an account object for the user
	var assetIds []string
	suffix := "000A"
	prefix := username + suffix
	var account = Account{ID: username, Prefix: prefix, CashBalance: 10000000.0, AssetsIds: assetIds}
	accountBytes, err := json.Marshal(&account)
	if err != nil {
		fmt.Println("error creating account" + account.ID)
		return nil, errors.New("Error creating account " + account.ID)
	}

	fmt.Println("Attempting to get state of any existing account for " + account.ID)
	existingBytes, err := stub.GetState(accountPrefix + account.ID)
	if err == nil {

		var company Account
		err = json.Unmarshal(existingBytes, &company)
		if err != nil {
			fmt.Println("Error unmarshalling account " + account.ID + "\n--->: " + err.Error())

			if strings.Contains(err.Error(), "unexpected end") {
				fmt.Println("No data means existing account found for " + account.ID + ", initializing account.")
				err = stub.PutState(accountPrefix + account.ID, accountBytes)

				if err == nil {
					fmt.Println("created account" + accountPrefix + account.ID)
					return nil, nil
				} else {
					fmt.Println("failed to create initialize account for " + account.ID)
					return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
				}
			} else {
				return nil, errors.New("Error unmarshalling existing account " + account.ID)
			}
		} else {
			fmt.Println("Account already exists for " + account.ID + " " + company.ID)
			return nil, errors.New("Can't reinitialize existing user " + account.ID)
		}
	} else {

		fmt.Println("No existing account found for " + account.ID + ", initializing account.")
		err = stub.PutState(accountPrefix + account.ID, accountBytes)

		if err == nil {
			fmt.Println("created account" + accountPrefix + account.ID)
			return nil, nil
		} else {
			fmt.Println("failed to create initialize account for " + account.ID)
			return nil, errors.New("failed to initialize an account for " + account.ID + " => " + err.Error())
		}

	}

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

func (t *SimpleChaincode) issueCommercialPaper(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating commercial paper")

	/*		0
		json
	  	{
			"ticker":  "string",
			"par": 0.00,
			"qty": 10,
			"discount": 7.5,
			"maturity": 30,
			"owners": [ // This one is not required
				{
					"company": "company1",
					"quantity": 5
				},
				{
					"company": "company3",
					"quantity": 3
				},
				{
					"company": "company4",
					"quantity": 2
				}
			],
			"issuer":"company2",
			"issueDate":"1456161763790"  (current time in milliseconds as a string)

		}
	*/
	//need one arg
	if len(args) != 1 {
		fmt.Println("error invalid arguments")
		return nil, errors.New("Incorrect number of arguments. Expecting one Audit trail record")
	}

	var cp CP
	var err error
	var account Account

	fmt.Println("Unmarshalling CP")
	err = json.Unmarshal([]byte(args[0]), &cp)
	if err != nil {
		fmt.Println("error invalid paper issue")
		return nil, errors.New("Invalid commercial paper issue")
	}

	//generate the CUSIP
	//get account prefix
	fmt.Println("Getting state of - " + accountPrefix + cp.Issuer)
	accountBytes, err := stub.GetState(accountPrefix + cp.Issuer)
	if err != nil {
		fmt.Println("Error Getting state of - " + accountPrefix + cp.Issuer)
		return nil, errors.New("Error retrieving account " + cp.Issuer)
	}
	err = json.Unmarshal(accountBytes, &account)
	if err != nil {
		fmt.Println("Error Unmarshalling accountBytes")
		return nil, errors.New("Error retrieving account " + cp.Issuer)
	}

	account.AssetsIds = append(account.AssetsIds, cp.CUSIP)

	// Set the issuer to be the owner of all quantity
	var owner Owner
	owner.Company = cp.Issuer
	owner.Quantity = cp.Qty

	cp.Owners = append(cp.Owners, owner)

	suffix, err := generateCUSIPSuffix(cp.IssueDate, cp.Maturity)
	if err != nil {
		fmt.Println("Error generating cusip")
		return nil, errors.New("Error generating CUSIP")
	}

	fmt.Println("Marshalling CP bytes")
	cp.CUSIP = account.Prefix + suffix

	fmt.Println("Getting State on CP " + cp.CUSIP)
	cpRxBytes, err := stub.GetState(cpPrefix + cp.CUSIP)
	if cpRxBytes == nil {
		fmt.Println("CUSIP does not exist, creating it")
		cpBytes, err := json.Marshal(&cp)
		if err != nil {
			fmt.Println("Error marshalling cp")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(cpPrefix + cp.CUSIP, cpBytes)
		if err != nil {
			fmt.Println("Error issuing paper")
			return nil, errors.New("Error issuing commercial paper")
		}

		fmt.Println("Marshalling account bytes to write")
		accountBytesToWrite, err := json.Marshal(&account)
		if err != nil {
			fmt.Println("Error marshalling account")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(accountPrefix + cp.Issuer, accountBytesToWrite)
		if err != nil {
			fmt.Println("Error putting state on accountBytesToWrite")
			return nil, errors.New("Error issuing commercial paper")
		}


		// Update the paper keys by adding the new key
		fmt.Println("Getting Paper Keys")
		keysBytes, err := stub.GetState("PaperKeys")
		if err != nil {
			fmt.Println("Error retrieving paper keys")
			return nil, errors.New("Error retrieving paper keys")
		}
		var keys []string
		err = json.Unmarshal(keysBytes, &keys)
		if err != nil {
			fmt.Println("Error unmarshel keys")
			return nil, errors.New("Error unmarshalling paper keys ")
		}

		fmt.Println("Appending the new key to Paper Keys")
		foundKey := false
		for _, key := range keys {
			if key == cpPrefix + cp.CUSIP {
				foundKey = true
			}
		}
		if foundKey == false {
			keys = append(keys, cpPrefix + cp.CUSIP)
			keysBytesToWrite, err := json.Marshal(&keys)
			if err != nil {
				fmt.Println("Error marshalling keys")
				return nil, errors.New("Error marshalling the keys")
			}
			fmt.Println("Put state on PaperKeys")
			err = stub.PutState("PaperKeys", keysBytesToWrite)
			if err != nil {
				fmt.Println("Error writting keys back")
				return nil, errors.New("Error writing the keys back")
			}
		}

		fmt.Println("Issue commercial paper %+v\n", cp)
		return nil, nil
	} else {
		fmt.Println("CUSIP exists")

		var cprx CP
		fmt.Println("Unmarshalling CP " + cp.CUSIP)
		err = json.Unmarshal(cpRxBytes, &cprx)
		if err != nil {
			fmt.Println("Error unmarshalling cp " + cp.CUSIP)
			return nil, errors.New("Error unmarshalling cp " + cp.CUSIP)
		}

		cprx.Qty = cprx.Qty + cp.Qty

		for key, val := range cprx.Owners {
			if val.Company == cp.Issuer {
				cprx.Owners[key].Quantity += cp.Qty
				break
			}
		}

		cpWriteBytes, err := json.Marshal(&cprx)
		if err != nil {
			fmt.Println("Error marshalling cp")
			return nil, errors.New("Error issuing commercial paper")
		}
		err = stub.PutState(cpPrefix + cp.CUSIP, cpWriteBytes)
		if err != nil {
			fmt.Println("Error issuing paper")
			return nil, errors.New("Error issuing commercial paper")
		}

		fmt.Println("Updated commercial paper %+v\n", cprx)
		return nil, nil
	}
}

func GetAllAudits(stub shim.ChaincodeStubInterface) ([]Audit, error) {
	var allAudits []Audit

	//Get list of all the keys
	keysBytes, err := stub.GetState("AuditKeys")
	if err != nil {
		fmt.Println("Error retrieving audit keys")
		return nil, errors.New("Error retrieving audit keys")
	}
	var keys []string
	err = json.Unmarshal(keysBytes, &keys)
	if err != nil {
		fmt.Println("Error unmarshalling audit keys")
		return nil, errors.New("Error unmarshalling audit keys")
	}

	//Get all the Audits
	for _, value := range keys {
		auditBytes, err := stub.GetState(value)

		var audit Audit
		err = json.Unmarshal(auditBytes, &audit)
		if err != nil {
			fmt.Println("Error retrieving audit " + value)
			return nil, errors.New("Error retrieving audit " + value)
		}

		fmt.Println("Appending Audit" + value)
		allAudits = append(allAudits, audit)
	}

	return allAudits, nil
}

func GetAudit(audit_hash string, stub shim.ChaincodeStubInterface) (Audit, error) {
	var audit Audit

	auditBytes, err := stub.GetState(audit_hash)
	if err != nil {
		fmt.Println("Error retrieving audit " + audit_hash)
		return audit, errors.New("Error retrieving cp " + audit_hash)
	}

	err = json.Unmarshal(auditBytes, &audit)
	if err != nil {
		fmt.Println("Error unmarshalling audit " + audit_hash)
		return audit, errors.New("Error unmarshalling cp " + audit_hash)
	}

	return audit, nil
}

func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Query running. Function: " + function)

	if function == "GetAllCPs" {
		fmt.Println("Getting all CPs")
		allCPs, err := GetAllCPs(stub)
		if err != nil {
			fmt.Println("Error from getallcps")
			return nil, err
		} else {
			allCPsBytes, err1 := json.Marshal(&allCPs)
			if err1 != nil {
				fmt.Println("Error marshalling allcps")
				return nil, err1
			}
			fmt.Println("All success, returning allcps")
			return allCPsBytes, nil
		}
	} else if function == "GetCP" {
		fmt.Println("Getting particular cp")
		cp, err := GetCP(args[0], stub)
		if err != nil {
			fmt.Println("Error Getting particular cp")
			return nil, err
		} else {
			cpBytes, err1 := json.Marshal(&cp)
			if err1 != nil {
				fmt.Println("Error marshalling the cp")
				return nil, err1
			}
			fmt.Println("All success, returning the cp")
			return cpBytes, nil
		}
	} else if function == "GetCompany" {
		fmt.Println("Getting the company")
		company, err := GetCompany(args[0], stub)
		if err != nil {
			fmt.Println("Error from getCompany")
			return nil, err
		} else {
			companyBytes, err1 := json.Marshal(&company)
			if err1 != nil {
				fmt.Println("Error marshalling the company")
				return nil, err1
			}
			fmt.Println("All success, returning the company")
			return companyBytes, nil
		}
	} else {
		fmt.Println("Generic Query call")
		bytes, err := stub.GetState(args[0])

		if err != nil {
			fmt.Println("Some error happenend: " + err.Error())
			return nil, err
		}

		fmt.Println("All success, returning from generic")
		return bytes, nil
	}
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("Invoke running. Function: " + function)

	if function == "issueCommercialPaper" {
		return t.issueCommercialPaper(stub, args)
	} else if function == "transferPaper" {
		return t.transferPaper(stub, args)
	} else if function == "createAccounts" {
		return t.createAccounts(stub, args)
	} else if function == "createAccount" {
		return t.createAccount(stub, args)
	}

	return nil, errors.New("Received unknown function invocation: " + function)
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Println("Error starting Simple chaincode: %s", err)
	}
}
