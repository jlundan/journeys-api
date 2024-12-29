package ggtfs

type FieldTobeValidated interface {
	IsValid() bool
	IsPresent() bool
	IsEmpty() bool
}

func validateRequiredFields(fields map[string]FieldTobeValidated, validationErrors *[]error, lineNumber int, fileName string) {
	for name, value := range fields {
		if !value.IsValid() {
			*validationErrors = append(*validationErrors, createFileRowError(fileName, lineNumber, createInvalidRequiredFieldString(name)))
		}
	}
}

func validateOptionalFields(fields map[string]FieldTobeValidated, validationErrors *[]error, lineNumber int, fileName string) {
	for name, value := range fields {
		if value.IsPresent() && !value.IsEmpty() && !value.IsValid() {
			*validationErrors = append(*validationErrors, createFileRowError(fileName, lineNumber, createInvalidFieldString(name)))
		}
	}
}
