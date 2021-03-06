/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2014-02-05 21:53
 * description :
 * history :
 */
package mvc

import (
	"fmt"
	"github.com/atnet/gof/web"
	"net/http"
	"reflect"
	"regexp"
	"strings"
)

var (
	actionRegexp = regexp.MustCompile("/([^/]+)$")
)

func CustomHandle(controller Controller, ctx *web.Context, action string, args ...interface{}) {
	//	defer func() {
	//		if err := recover(); err != nil {
	//			if _, ok := err.(error); ok {
	//
	//				w.Write([]byte("Inject error, plese check parameter type or index!"))
	//			}
	//		}
	//	}()

	w := ctx.ResponseWriter

	// 拦截器
	filter, isFilter := controller.(Filter)
	if isFilter {
		if !filter.Requesting(ctx) { //如果返回false，终止执行
			return
		}
	}

	t := reflect.ValueOf(controller)
	method := t.MethodByName(action)

	if !method.IsValid() {
		errMsg := "No action named \"" + strings.Replace(action, "_post", "", 1) +
			"\" in " + reflect.TypeOf(controller).String() + "."
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	} else {
		//包含基础的ResponseWriter和http.Request 2个参数
		argsLen := len(args)
		numIn := method.Type().NumIn()

		if argsLen < numIn-1 {
			errMsg := fmt.Sprintf("Can't inject to method,it's possible missing parameter!"+
				"\r\ncontroller: %s , action: %s",
				reflect.TypeOf(controller).String(), action)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		} else {
			params := make([]reflect.Value, numIn)
			params[0] = reflect.ValueOf(ctx)
			for i := 1; i < numIn; i++ {
				params[i] = reflect.ValueOf(args[i-1])
			}

			method.Call(params)
		}
	}

	if isFilter {
		filter.RequestEnd(ctx)
	}
}

//控制器处理
//@controller ： 包含多种动作，URL中的文件名自动映射到控制器的函数
//				 注意，是区分大小写的,默认映射到index函数
//				 如果是POST请求将映射到控制器“函数名+_post”的函数执行
// @re_post : 是否为post请求额外加上_post来区分Post和Get请求
func Handle(controller Controller, ctx *web.Context, rePost bool, args ...interface{}) {
	r := ctx.Request
	// 处理末尾的/
	var path = r.URL.Path
	if strings.HasSuffix(path, "/") {
		path = path[:len(path)-1]
	}

	var action string
	groups := actionRegexp.FindAllStringSubmatch(path, 1)
	if len(groups) == 0 || len(groups[0]) == 0 {
		action = "Index"
	} else {
		action = groups[0][1]

		//去扩展名
		extIndex := strings.Index(action, ".")
		if extIndex != -1 {
			action = action[0:extIndex]
		}

		//将第一个字符转为大写,这样才可以
		upperFirstLetter := strings.ToUpper(action[0:1])
		if upperFirstLetter != action[0:1] {
			action = upperFirstLetter + action[1:]
		}
	}

	if rePost && r.Method == "POST" {
		action += "_post"
	}

	CustomHandle(controller, ctx, action, args...)
}
