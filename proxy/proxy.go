package proxy

import (
	"errors"
	"fmt"
	"reflect"
)

// 提供动态调用方法接口
type InvocationHandler interface {
	Invoke(proxy *Proxy, method *Method, args []interface{}) ([]interface{}, error)
}

// 代理，用来总管代理类的生成
type Proxy struct {
	target  interface{}        //目标类，后面的类型和java的Object一样
	methods map[string]*Method //map用来装载待增强的不同的方法
	handle  InvocationHandler  //用来暴露统一invoke接口，类似多态
}

// 创建新的代理
func NewProxy(target interface{}, h InvocationHandler) *Proxy {
	typ := reflect.TypeOf(target)          //用来显示目标类动态的真实类型
	value := reflect.ValueOf(target)       //获取目标类的值
	methods := make(map[string]*Method, 0) //初始化目标类的方法map
	//将目标类的方法逐个装载
	for i := 0; i < value.NumMethod(); i++ {
		method := value.Method(i)
		methods[typ.Method(i).Name] = &Method{value: method}
	}
	return &Proxy{target: target, methods: methods, handle: h}
}

// 代理调用代理方法
func (p *Proxy) InvokeMethod(name string, args ...interface{}) ([]interface{}, error) {
	return p.handle.Invoke(p, p.methods[name], args)
}

// 用来承载目标类的方法定位和调用
type Method struct {
	value reflect.Value //用来装载方法实例
}

// 这里相当于调用原方法，在该方法外可以做方法增强，需要调用者自己实现！！！
func (m *Method) Invoke(args ...interface{}) (res []interface{}, err error) {
	defer func() {
		//用来捕捉异常
		if p := recover(); p != nil {
			err = errors.New(fmt.Sprintf("%s", p))
		}
	}()

	//处理参数
	params := make([]reflect.Value, 0)
	if args != nil {
		for i := 0; i < len(args); i++ {
			params = append(params, reflect.ValueOf(args[i]))
		}
	}

	//调用方法
	call := m.value.Call(params)

	//接收返回值
	res = make([]interface{}, 0)
	if call != nil && len(call) > 0 {
		for i := 0; i < len(call); i++ {
			res = append(res, call[i].Interface())
		}
	}
	return
}
