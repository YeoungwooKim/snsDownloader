package validation

import "github.com/go-playground/validator/v10"

var Validate = validator.New()
var VideoRules = map[string]interface{}{
	"platform": "required,min=5",
	"uri":      "required,uri",
}
