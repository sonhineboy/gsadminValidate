# gsadminValidate
part of gsadmin, Elegant handling of custom validation for gin!

## 安装
```go
go get -u github.com/sonhineboy/gsadminValidator
```
## 初始化
### 中文翻译
```go
// initTrans 初始化 中文翻译
func initTrans() *ginValidator.Trans {
	tran := ginValidator.NewDefaultTrans()
	err := tran.SetUp()
	if err != nil {
		panic(err)
	}
	return tran
}
```
### 初始化自定义

```go
// initCustomValidator 第一步
func initCustomValidator(tran *ginValidator.Trans) *ginValidator.CustomValidatorManager {
	customValidator := ginValidator.NewCustomValidatorManager(make(map[string]ginValidator.CustomValidator), tran.GetValidate(), tran.GetTrans())
	// 注册验证规则，内部使用map类型注册，为协程不安全的，所以需要项目出事化是注册
	customValidator.Adds(
		new(validators.DemoValidator),
	)
	customValidator.RegisterToValidate()
	return customValidator
}

//执行
initCustomValidator(initTrans())

```

## 自定义验证规则

```go

package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/sonhineboy/gsadminValidator/ginValidator"
)
// DemoValidator 命名规则 名字+Validator
type DemoValidator struct {
	ginValidator.BaseValidator 
}

//TagName 返回规则名字
func (d *DemoValidator) TagName() string {
	return "demo"
}

//Messages 规则错误提示信息
func (d *DemoValidator) Messages() string {
	return "{0}长度必须超过6个"
}
//Validator 规则验证逻辑
func (d *DemoValidator) Validator(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) > 6
}
```
## 注册验证规则

```go
customValidator.Adds(
new(validators.DemoValidator),...
)

```
> 声明的是：要在项目的初始化的时候注册进去

## 使用验证规则
```go
type Login struct {
	Name     string `json:"name" form:"name" binding:"required,demo=xxx" label:"名字"`
	Password string `json:"password" form:"password" binding:"required"`
}
```
> 其中 label:是对字段的翻译，如果不需要可以不用

## 执行效果
```json
//demo 规则不通过
{
    "code": 402,
    "msg": "名字长度必须超过6个",
    "data": ""
}

//无label 效果
{
  "code": 402,
  "msg": "Name 长度必须超过6个",
  "data": ""
}
```