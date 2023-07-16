package libraries

import (
	"backend/app/config"
	"backend/app/models"
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"gorm.io/gorm"
)

type validation struct {
	conn *gorm.DB
}

func NewValidation() *validation {
	conn := config.DBConn()
	db, err := gorm.Open(conn, &gorm.Config{})
	if err != nil {
		fmt.Println("error")
	}
	return &validation{
		conn: db,
	}
}

func (v *validation) Init() (*validator.Validate, ut.Translator) {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	validate.RegisterTagNameFunc(func(field reflect.StructField) string {
		labelName := field.Tag.Get("label")
		return labelName
	})

	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} tidak boleh kosong", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	validate.RegisterValidation("isunique", func(fl validator.FieldLevel) bool {
		params := fl.Param()
		split_param := strings.Split(params, "-")

		var tablename models.User
		fieldname := split_param[1]
		fieldvalue := fl.Field().String()

		return v.checkIsUnique(tablename, fieldname, fieldvalue)
	})

	validate.RegisterTranslation("isunique", trans, func(ut ut.Translator) error {
		return ut.Add("isunique", "{0} sudah digunakan", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("isunique", fe.Field())
		return t
	})

	return validate, trans

}

func (c *validation) checkIsUnique(tablename interface{}, fieldname, fieldvalue string) bool {
	var count int64

	err := c.conn.Debug().Model(&tablename).Select(fieldname).Where(fieldname+" = ?", fieldvalue).Count(&count).Error
	if err != nil {
		log.Fatal(err)
	}

	return count < 1
}

func (v *validation) Struct(s interface{}) interface{} {
	validate, trans := v.Init()
	vErrors := make(map[string]interface{})

	err := validate.Struct(s)

	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			vErrors[e.StructField()] = e.Translate(trans)
		}
	}

	if len(vErrors) > 0 {
		return vErrors
	}
	return nil
}
