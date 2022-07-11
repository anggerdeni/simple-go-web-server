package helpers

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"

	"assignment-runner-base/entity"
)

func ReadAssignmentConfigFromDirectory(path string) (*entity.AssignmentConf, error) {
	jsonFile, err := os.Open(path + "/assignment-config.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	// we initialize our Users array
	var assignmentConf entity.AssignmentConf

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &assignmentConf)

	if assignmentConf.Command == "" {
		return nil, errors.New("invalid assignment config file, command is required")
	}

	return &assignmentConf, nil
}
