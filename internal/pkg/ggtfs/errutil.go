package ggtfs

import (
	"errors"
	"fmt"
)

func createFileRowError(fileName string, row int, err string) error {
	return errors.New(fmt.Sprintf("%s:%v: %s", fileName, row, err))
}

func createFileRowRecommendation(fileName string, row int, err string) string {
	return fmt.Sprintf("%s:%v: %s", fileName, row, err)
}

func createFileError(fileName string, err string) error {
	return errors.New(fmt.Sprintf("%s: %s", fileName, err))
}

func createInvalidFieldString(fieldName string) string {
	return fmt.Sprintf("invalid field: %s", fieldName)
}

func createInvalidRequiredFieldString(fieldName string) string {
	return fmt.Sprintf("invalid mandatory field: %s", fieldName)
}
