/*
Copyright IBM Corp 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"errors"
	"fmt"

	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

// An entry for a product
type DataEntry struct {
	PlaceId     string `json:"placeid"`
	Temperature string `json:"temperature"`
	Timestamp   string `json:"timestamp"`
}

// ============================================================================================================================
// Main
// ============================================================================================================================
func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}

// Init resets all the things
func (t *SimpleChaincode) Init(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	err := stub.PutState("hello_world", []byte(args[0]))
	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Invoke is our entry point to invoke a chaincode function
func (t *SimpleChaincode) Invoke(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("invoke is running " + function)

	// Handle different functions
	if function == "init" { //initialize the chaincode state, used as reset
		return t.Init(stub, "init", args)
	} else if function == "write" {
		return t.write(stub, args)
	}

	fmt.Println("invoke did not find func: " + function) //error

	return nil, errors.New("Received unknown function invocation")
}

// Query is our entry point for queries
func (t *SimpleChaincode) Query(stub *shim.ChaincodeStub, function string, args []string) ([]byte, error) {
	fmt.Println("query is running " + function)

	// Handle different functions
	if function == "read" { //read a variable
		return t.read(stub, args)
	}
	fmt.Println("query did not find func: " + function) //error

	return nil, errors.New("Received unknown function query")
}

// args has productid, placeid, temperature, and timestamp as strings.
func (t *SimpleChaincode) write(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var productid, placeid, temperature, timestamp string
	var entries []DataEntry
	var err error
	fmt.Println("running write()")

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Expecting 4. Product id, place id, temperature, and taime stamp.")
	}

	productid = args[0] //rename for fun
	placeid = args[1]
	temperature = args[2]
	timestamp = args[3]
	dataAsBytes, err := stub.GetState(productid)
	if err != nil {
		return nil, err
	}

	// the data may be empty, in that case just go straight to adding the new
	// information.
	if dataAsBytes != nil {
		// we have a json string, unmarshal it
		json.Unmarshal(dataAsBytes, entries)
	}
	// add the new data entry to the json object
	new_entry := DataEntry{
		PlaceId:     placeid,
		Temperature: temperature,
		Timestamp:   timestamp,
	}
	entries = append(entries, new_entry)
	// now marshal it back to json, and write it to the cblockchain
	new_json, _ := json.Marshal(entries)
	new_json_str := string(new_json)
	fmt.Println(new_json_str)

	err = stub.PutState(productid, []byte(new_json_str)) //write the variable into the chaincode state
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (t *SimpleChaincode) read(stub *shim.ChaincodeStub, args []string) ([]byte, error) {
	var productid, jsonResp string
	var err error

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting productid to query")
	}

	productid = args[0]
	valAsbytes, err := stub.GetState(productid)
	if err != nil {
		jsonResp = "{\"Error\":\"Failed to get state for " + productid + "\"}"
		return nil, errors.New(jsonResp)
	}

	return valAsbytes, nil
}
