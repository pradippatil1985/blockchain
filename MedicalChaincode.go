package main

import (
	"errors"
	"fmt"
	//"strconv"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	//"github.com/golang/protobuf/ptypes/timestamp"
)

// Region Chaincode implementation
type MedicalChaincode struct {
}

var medicalIndexTxStr = "_medicalIndexTxStr"

type MedicalData struct{
	PATIENT_ID string `json:"PATIENT_ID"`
	PATIENT_NAME string `json:"PATIENT_NAME"`
	DOC string `json:"DOC"`
  EXPIRY_DATE string `json:"EXPIRY_DATE"`

}

func (t *MedicalChaincode) Init(stub shim.ChaincodeStubInterface) ([]byte, error) {
	var err error
	// Initialize the chaincode

	fmt.Printf("Patient's Health Records Tracking started\n")

	var patientPolicyTxs []MedicalData
	jsonAsBytes, _ := json.Marshal(patientPolicyTxs)
	err = stub.PutState(medicalIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}


	return nil, nil
}

// Add Patient data
func (t *MedicalChaincode) Invoke(stub shim.ChaincodeStubInterface) ([]byte, error) {
	if function == medicalIndexTxStr {
		return t.AddPatientInfo(stub, args)
	}
	return nil, nil
}

func (t *MedicalChaincode) AddPatientInfo(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var MedicalDataObj MedicalData
	var MedicalDataList []MedicalData
	var err error

	if len(args) != 4 {
		return nil, errors.New("Incorrect number of arguments. Need 3 arguments")
	}

	// Initialize the chaincode
	MedicalDataObj.PATIENT_ID = args[0]
	MedicalDataObj.PATIENT_NAME = args[1]
	MedicalDataObj.DOC = args[2]
	MedicalDataObj.EXPIRY_DATE = args[3]

        fmt.Printf("Input from Hospital:%s\n", MedicalDataObj)

	medicalTxsAsBytes, err := stub.GetState(medicalIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get Patient Details")
	}
	json.Unmarshal(medicalTxsAsBytes, &MedicalDataList)

	MedicalDataList = append(MedicalDataList, MedicalDataObj)
	jsonAsBytes, _ := json.Marshal(MedicalDataList)

	err = stub.PutState(medicalIndexTxStr, jsonAsBytes)
	if err != nil {
		return nil, err
	}
	return nil, nil

}

// Query callback representing the query of a chaincode
func (t *MedicalChaincode) Query(stub shim.ChaincodeStubInterface,function string, args []string) ([]byte, error) {

	var PATIENT_ID string // Entities
	var err error
	var resAsBytes []byte

	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting Patient ID to query")
	}

	PATIENT_ID = args[0]

	resAsBytes, err = t.GetPatientDetails(stub, PATIENT_ID)

	fmt.Printf("Query Response:%s\n", resAsBytes)

	if err != nil {
		return nil, err
	}

	return resAsBytes, nil
}

func (t *MedicalChaincode)  GetPatientDetails(stub shim.ChaincodeStubInterface, PATIENT_ID string) ([]byte, error) {

	//var requiredObj MedicalData
	var objFound bool
	PatientTxsAsBytes, err := stub.GetState(medicalIndexTxStr)
	if err != nil {
		return nil, errors.New("Failed to get Patient Details")
	}
	var PatientTxObjects []MedicalData
	var PatientTxObjects1 []MedicalData
	json.Unmarshal(PatientTxsAsBytes, &PatientTxObjects)
	length := len(PatientTxObjects)
	fmt.Printf("Output from chaincode: %s\n", PatientTxsAsBytes)

	if PATIENT_ID == "" {
		res, err := json.Marshal(PatientTxObjects)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}

	objFound = false
	// iterate
	for i := 0; i < length; i++ {
		obj := PatientTxObjects[i]
		if PATIENT_ID == obj.PATIENT_ID {
			PatientTxObjects1 = append(PatientTxObjects1,obj)
			//requiredObj = obj
			objFound = true
		}
	}

	if objFound {
		res, err := json.Marshal(PatientTxObjects1)
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	} else {
		res, err := json.Marshal("No Data found")
		if err != nil {
		return nil, errors.New("Failed to Marshal the required Obj")
		}
		return res, nil
	}


}

func main() {
	err := shim.Start(new(MedicalChaincode))
	if err != nil {
		fmt.Printf("Error starting Medical chaincode: %s", err)
	}
}
