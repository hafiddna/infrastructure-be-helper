package helper

import (
	"fmt"
	goValidator "github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"gorm.io/gorm"
	"reflect"
	"strconv"
	"strings"
)

var tagTranslations = map[string]string{
	"gte": "min",
	"lte": "max",
}

func Validate(body reflect.Type, errors error) interface{} {
	errorMessages := make(map[string]interface{})

	for _, err := range errors.(goValidator.ValidationErrors) {
		ns := err.Namespace() // e.g. "SaveLecturerGradingRequest.Grades[0].StudentID"
		parts := strings.Split(ns, ".")

		// 1. Check if it's a nested slice (e.g., Grades[0].StudentID)
		if len(parts) >= 2 && strings.Contains(parts[1], "[") {
			sliceFieldWithIndex := parts[1] // Grades[0]
			sliceFieldName := strings.Split(sliceFieldWithIndex, "[")[0]

			// Get json/form/query tag for the slice field
			sliceStructField, _ := body.FieldByName(sliceFieldName)
			jsonFieldTag := getFieldTag(sliceStructField, sliceFieldName)

			indexStr := strings.Split(strings.Split(parts[1], "[")[1], "]")[0]
			index, _ := strconv.Atoi(indexStr)

			// To get the item field tag, we need the slice element type
			sliceFieldType := sliceStructField.Type // should be []SomeStruct
			elemType := sliceFieldType.Elem()
			if elemType.Kind() == reflect.Ptr {
				elemType = elemType.Elem()
			}

			// Now get the field tag of the nested field
			itemStructField, ok := elemType.FieldByName(err.StructField())
			itemFieldTag := err.StructField()
			if ok {
				itemFieldTag = getFieldTag(itemStructField, err.StructField())
			}

			// Initialize slice of maps if needed
			if _, ok := errorMessages[jsonFieldTag]; !ok {
				errorMessages[jsonFieldTag] = []map[string][]string{}
			}
			slice := errorMessages[jsonFieldTag].([]map[string][]string)

			for len(slice) <= index {
				slice = append(slice, map[string][]string{})
			}

			errTag := translateTag(err.Tag())
			slice[index][itemFieldTag] = append(slice[index][itemFieldTag], errTag)
			errorMessages[jsonFieldTag] = slice

		} else {
			// 2. Top-level field
			field, _ := body.FieldByName(err.StructField())
			fieldTag := getFieldTag(field, err.StructField())

			errTag := translateTag(err.Tag())

			if existing, ok := errorMessages[fieldTag]; ok {
				errorMessages[fieldTag] = append(existing.([]string), errTag)
			} else {
				errorMessages[fieldTag] = []string{errTag}
			}
		}
	}

	return errorMessages
}

func getFieldTag(field reflect.StructField, fallback string) string {
	tag := field.Tag.Get("json")
	if tag == "" {
		tag = field.Tag.Get("form")
	}
	if tag == "" {
		tag = field.Tag.Get("query")
	}
	if tag == "" {
		return strings.ToLower(fallback)
	}
	if idx := strings.Index(tag, ","); idx != -1 {
		tag = tag[:idx]
	}
	return strings.ToLower(tag)
}

func translateTag(tag string) string {
	if translated, ok := tagTranslations[tag]; ok {
		return translated
	}
	return tag
}

type CustomValidator struct {
	postgresql *gorm.DB
	mongodb    *mongo.Database
	minio      *minio.Client
}

func (v *CustomValidator) exist(table string, field string, val interface{}) bool {
	var c int64

	if err := v.postgresql.Table(table).Where(fmt.Sprintf("%v = ?", field), val).Limit(1).Count(&c).Error; err != nil {
		return false
	}

	return c == 1
}

// a custom validator with format : validate:"exist={table_name}.{field_name}",
func (v *CustomValidator) dbExist(fl goValidator.FieldLevel) bool {
	param := fl.Param()
	parts := strings.Split(param, ".")

	table := parts[0]
	field := parts[1]
	fmt.Printf("table: %v, field: %v\n", table, field)

	value := fl.Field().Interface()

	if len(parts) == 2 {
		return v.exist(table, field, value)
	}

	return false
}

func (v *CustomValidator) unique(table string, field string, val interface{}) bool {
	var c int64

	if err := v.postgresql.Table(table).Unscoped().Where(fmt.Sprintf("%v = ?", field), val).Limit(1).Count(&c).Error; err != nil {
		return false
	}

	return c == 0
}

// a custom validator with format : validate:"unique={table_name}.{field_name}",
func (v *CustomValidator) dbUnique(fl goValidator.FieldLevel) bool {
	param := fl.Param()
	parts := strings.Split(param, ".")
	if len(parts) < 2 {
		return false
	}

	table := parts[0]
	field := parts[1]

	if len(parts) == 4 {
		ignoreField := parts[2]
		ignoreValue := fl.Parent().FieldByName(ignoreField).Interface()
		column := parts[3]

		var c int64

		if err := v.postgresql.
			Table(table).
			Where(fmt.Sprintf("%v = ? AND %v != ?", field, column), fl.Field().Interface(), ignoreValue).
			Limit(1).
			Count(&c).
			Error; err != nil {
			return false
		}

		return c == 0
	} else {
		fmt.Printf("table: %v, field: %v\n", table, field)
		value := fl.Field().Interface()
		if !v.unique(table, field, value) {
			return false
		}
		return true
	}
}

func (v *CustomValidator) Register(validate *goValidator.Validate) {
	validate.RegisterValidation("exist", v.dbExist)
	validate.RegisterValidation("unique", v.dbUnique)
}

func NewCustomValidator(postgresql *gorm.DB, mongodb *mongo.Database, minio *minio.Client) *CustomValidator {
	val := &CustomValidator{
		postgresql: postgresql,
		mongodb:    mongodb,
		minio:      minio,
	}
	return val
}
