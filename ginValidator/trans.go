package ginValidator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
)

type Option func(t *Trans)

// getZhTans 中文转换器
func getZhTans() ut.Translator {
	uni := ut.New(zh.New())
	trans, ok := uni.GetTranslator("zh")
	if ok {
		panic("zh trans error")
	}
	return trans
}

// getDefaultValidate 获取默认gin 验证引擎
func getDefaultValidate() *validator.Validate {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("Validate init err")
	}
	return v
}

type Trans struct {
	validate *validator.Validate
	tans     ut.Translator
}

func NewTrans(option ...Option) *Trans {
	ts := new(Trans)
	for _, v := range option {
		v(ts)
	}
	return ts
}

func NewDefaultTrans() *Trans {
	return &Trans{
		validate: getDefaultValidate(),
		tans:     getZhTans(),
	}
}

func (t *Trans) GetValidate() *validator.Validate {
	return t.validate
}

func (t *Trans) GetTrans() ut.Translator {
	return t.tans
}

// SetUp 设置验证器
func (t *Trans) SetUp() error {
	t.registerTagFunc()
	return zhtranslations.RegisterDefaultTranslations(t.GetValidate(), t.GetTrans())
}

func (t *Trans) registerTagFunc() {
	t.GetValidate().RegisterTagNameFunc(func(field reflect.StructField) string {
		label := field.Tag.Get("label")
		if len(label) > 0 {
			return label
		} else {
			return field.Name
		}
	})
}
