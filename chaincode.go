package main

import (
	"errors"
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// ============================================================================================================================
// Main
// ============================================================================================================================
/*
 *	--> Executes when each peer deploys their instance of the chaincode
 * 	--> shim.Start() sets up the communication between this chaincode and the peer that deployed it
 */
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
/*
 *	Init is called when you first deploy your chaincode.
 * 	--> Used to do any initialization your chaincode need 
 */
func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if len(args) != 3 {
		return nil, errors.New("Incorrect number of arguments. Expecting 3")
	}

	// Stores the elements in the args argument to the right key
	errId := stub.PutState("id", []byte(args[0]))
	if errId != nil {
		return nil, errId
	}

	errName := stub.PutState("name", []byte(args[1]))
	if errName != nil {
		return nil, errName
	}

	errCert := stub.PutState("certification", []byte(args[2]))
	if errCert != nil {
		return nil, errCert
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
/*
 *	Invoke is called when you want to call chaincode functions to do real work.
 *	--> Invocations will be captured as transactions, which get grouped into blocks on the chain
 *	--> When you need to update the ledger, you will do so by invoking your chaincode
 */
func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" {						// initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {				// calls the write function
		return t.write(stub, args)
	}
	fmt.Println("invoke did not find func: " + function)	// error

	return nil, errors.New("Received unknown function invocation: " + function)
}

// #### Write function
/*
 *	Write is called by the Invoke function.
 * 	--> Checks for a certain number of arguments, and then write a key/value pait to the ledger
 * 	--> Stores any key/value pair you want into the blockchain ledger
 */
 func (t *SimpleChaincode) write(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
 	var key, value string
 	var err error
 	fmt.Println("running write()")

 	if len(args) != 2 {
 		return nil, errors.New("Incorrect number of arguments. Expecting 2: name of the key and value to set")
 	}

 	key = args[0]
 	value = args[1]
 	err = stub.PutState(key, []byte(value))		// write the variable into the chaincode state
 	if err != nil {
 		return nil, err
 	}

 	return nil, nil
 }

// Query is our entry point for queries
/*
 *	Query is called whenever you query your chaincode's state.
 *	--> Do not add blocks to the chain
 *	--> Used to read the value of your chaincode state's key/value pairs 
 */
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" {						// calls the read function
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function)		//error

	return nil, errors.New("Received unknown function query: " + function)
}

// #### Read function
/*
 *	Read is called by the Query function.
 *	--> Reads the value of a previously written key
 */
func (t *SimpleChaincode) read(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	var key, jsonResp string
	var err error

	if len(args) != 1 {
 		return nil, errors.New("Incorrect number of arguments. Expecting 1: name of the key to query")
	}

	key = args[0]
	valAsbytes, err := stub.GetState(key)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + key + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
