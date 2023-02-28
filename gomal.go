package gomal

type ValidationResult struct {
	Name     string
	Messages []string
}

func Validate(validators ...Validator) []ValidationResult {
	results := []ValidationResult{}
	for _, validator := range validators {
		if len(validator.errorMessages) > 0 {
			results = append(results, ValidationResult{
				Name:     validator.name,
				Messages: validator.errorMessages,
			})
		}
	}
	return results
}
