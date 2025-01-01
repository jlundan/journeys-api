package validation

type Result interface {
	Code() string
}

type InvalidCharactersResult struct{}

func (e InvalidCharactersResult) Code() string {
	return "invalid_characters"
}

type MissingRequiredFieldResult struct {
	FieldName string
}

func (e MissingRequiredFieldResult) Code() string {
	return "missing_required_field"
}
