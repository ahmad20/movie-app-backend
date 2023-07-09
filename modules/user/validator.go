package user

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func BlacklistValidation(fl validator.FieldLevel) bool {
	field := fl.Field().String()

	if field == "" {
		return true
	}

	match, _ := regexp.MatchString(`^[^'"\[\]<>\{\}]+$`, field)
	return match
}

// func BlacklistValidation(field string) validation.RuleFunc {
// 	return func(value interface{}) error {
// 		val, ok := value.(string)
// 		if !ok {
// 			return errors.New("The " + field + " is not a string")
// 		}

// 		if val == "" {
// 			return nil
// 		}

// 		match, _ := regexp.MatchString(`^[^'"\[\]<>\{\}]+$`, val)
// 		if !match {
// 			return errors.New("The " + field + " contains unsafe characters")
// 		}

// 		return nil
// 	}
// }

// func DatetimeValidation(field string) validation.RuleFunc {
// 	return func(value interface{}) error {
// 		val, ok := value.(string)
// 		if !ok {
// 			return errors.New("The " + field + " is not datetime format")
// 		}

// 		_, err := time.Parse("2006-01-02 15:04:05", val)
// 		if err != nil {
// 			return errors.New("the " + field + " is not a datetime format")
// 		}

// 		return nil
// 	}
// }

//	func (m Request) Validate() interface{} {
//		return govalidator.MapData{
//			"page":    []string{"numeric"},
//			"limit":   []string{"numeric"},
//			"orderBy": []string{"in:id,name,created_at,updated_at"},
//			"sortBy":  []string{"in:asc,desc"},
//			"search":  []string{"blacklist"},
//		}
//	}
// func (m Register) Validate() interface{} {
// 	return govalidator.MapData{
// 		"username": []string{"blacklist", "required"},
// 		"password": []string{"required"},
// 		"name":     []string{"blacklist"},
// 		"age":      []string{"blacklist, numeric", "required"},
// 	}
// }
