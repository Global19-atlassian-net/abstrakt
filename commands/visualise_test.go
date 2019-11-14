package commands

import (
	// "bufio"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestFileExists(t *testing.T) {

	//Setup variables and content for test
	testValidFilename := "testVisualise.out"
	testInvalidFilename := "nonexistant.out"
	testData := []byte("A file to test with")

	//Create a file to test against
	var err error = ioutil.WriteFile(testValidFilename, testData, 0644)
	if err != nil {
		fmt.Println("Could not create output testing file, cannot proceed")
		t.Error(err)
	}

	//Test that a valid file (created above) can be seen
	var result bool = fileExists(testValidFilename) //Expecting true - file does exists
	if result == false {
		t.Errorf("Test file does exist but testFile returns that it does not")
	}

	//Test that an invalid file (does not exist) is not seen
	result = fileExists(testInvalidFilename) //Expecting false - file does not exist
	if result != false {
		t.Errorf("Test file does not exist but testFile says it does")
	}

	err = os.Remove(testValidFilename)
	if err != nil {
		panic(err)
	}

	result = fileExists(testValidFilename) //Expecting false - file has been removed
	if result == true {
		t.Errorf("Test file has been removed but fileExists is finding it")
	}

}

func TestParseYaml(t *testing.T) {

	// var testInvalidYAMLString string = `Name: "Azure Event Hubs Sample"
	// Id: "d6e4a5e9-696a-4626-ba7a-534d6ff450a5"
	// Services:
	//   `

	testValidYAMLString := `
Description: "Event Generator to Event Hub connection"
From: 9e1bcb3d-ff58-41d4-8779-f71e7b8800f8
Id: 211a55bd-5d92-446c-8be8-190f8f0e623e
Name: "Azure Event Hubs Sample"
Properties: {}
Relationships: 
  - 
    Name: "Generator to Event Hubs Link"
Services: 
  - 
    Name: "Event Generator"
To: 3aa1e546-1ed5-4d67-a59c-be0d5905b490
Type: EventGenerator
`

	var retConfig Config = parseYaml(testValidYAMLString)

	if retConfig.Name != "Azure Event Hubs Sample" &&
		retConfig.ID != "d6e4a5e9-696a-4626-ba7a-534d6ff450a5" &&
		len(retConfig.Services) != 1 &&
		len(retConfig.Relationships) != 1 {
		t.Errorf("YAML did not parse correctly and it should have")
	}

	//Cannot test invalid YAML without a custom log function or rewrite of the parseYAML function
	// defer func() {
	// 	if r := recover(); r == nil {
	// 		t.Errorf("The code did not panic")
	// 	}
	// }()

	// parseYaml(testInvalidYAMLString)

}
