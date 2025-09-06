/*
 * Copyright (C) distroy
 */

package ldhook

import (
	"fmt"
	"reflect"
)

func getMethod(obj any, methodName string) reflect.Method {
	tType, ok := obj.(reflect.Type)
	if !ok || tType == nil {
		tType = reflect.TypeOf(obj)
	}

	method, ok := tType.MethodByName(methodName)
	if !ok {
		panic(fmt.Sprintf("retrieve method by name fail. type:%s, method:%s",
			tType.String(), methodName))
	}

	return method
}

func GetMethod(obj any, methodName string) any {
	return getMethod(obj, methodName).Func.Interface()
}
