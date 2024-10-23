package ginValidator

import (
	"fmt"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type CustomValidator interface {
	TagName() string
	Validator(fl validator.FieldLevel) bool
	Messages() string
	RegisterTranslationsFunc(tag string, msg string) validator.RegisterTranslationsFunc
	TranslationFunc(ut ut.Translator, fe validator.FieldError) string
}

type BaseValidator struct {
}

func (d *BaseValidator) RegisterTranslationsFunc(tag string, msg string) validator.RegisterTranslationsFunc {
	return func(trans ut.Translator) error {
		if err := trans.Add(tag, msg, false); err != nil {
			return err
		}
		return nil
	}
}

func (d *BaseValidator) TranslationFunc(trans ut.Translator, fe validator.FieldError) string {
	msg, err := trans.T(fe.Tag(), fe.Field())
	if err != nil {
		panic(fe.(error).Error())
	}
	return msg
}

// CustomValidatorManager 自定义验证器管理器
// 验证管理器没里Map 没有锁，所以只能在项目启动时注册
type CustomValidatorManager struct {
	customValidators map[string]CustomValidator
	validate         *validator.Validate
	trans            ut.Translator
}

func NewCustomValidatorManager(customValidators map[string]CustomValidator, validate *validator.Validate, trans ut.Translator) *CustomValidatorManager {
	return &CustomValidatorManager{
		customValidators: customValidators,
		validate:         validate,
		trans:            trans,
	}
}

func (c *CustomValidatorManager) GetValidate() *validator.Validate {
	return c.validate
}

func (c *CustomValidatorManager) GetTrans() ut.Translator {
	return c.trans
}

// Adds 注册一个验证器到 manager 中
func (c *CustomValidatorManager) Adds(validators ...CustomValidator) {

	for _, customValidator := range validators {
		c.Add(customValidator)
	}

}

// Add 注册一个验证器到 manager 中
func (c *CustomValidatorManager) Add(validator CustomValidator) {
	if c.HasValidator(validator.TagName()) {
		panic(fmt.Sprintf("validator %s is has!", validator.TagName()))
	} else {
		c.customValidators[validator.TagName()] = validator
	}
}

// RegisterToValidate 把自定义验证器注册到系统中
func (c *CustomValidatorManager) RegisterToValidate() {
	for key, customValidator := range c.customValidators {
		_ = c.validate.RegisterValidation(key, customValidator.Validator)
		_ = c.validate.RegisterTranslation(key, c.trans, customValidator.RegisterTranslationsFunc(customValidator.TagName(), customValidator.Messages()), customValidator.TranslationFunc)
	}
}

// HasValidator 名字为key验证器是否存在
func (c *CustomValidatorManager) HasValidator(key string) bool {
	_, ok := c.customValidators[key]
	return ok
}

// GetCustomValidators  获取已注册的验证器
func (c *CustomValidatorManager) GetCustomValidators() map[string]CustomValidator {
	return c.customValidators
}
