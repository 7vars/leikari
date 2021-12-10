package leikari

import (
	"reflect"
)

var (
	ErrorType   = reflect.TypeOf((*error)(nil)).Elem()
	ActorContextType = reflect.TypeOf((*ActorContext)(nil)).Elem()
	MessageType = reflect.TypeOf((*Message)(nil)).Elem()
)

func IsErrorType(gtype reflect.Type) bool {
	return gtype.Implements(ErrorType)
}

func IsActorContextType(gtype reflect.Type) bool {
	return gtype.Implements(ActorContextType)
}

func IsMessageType(gtype reflect.Type) bool {
	return gtype.Implements(MessageType)
}

func PtrValue(gvalue reflect.Value) reflect.Value {
	if gvalue.Kind() != reflect.Ptr {
		pto := reflect.PtrTo(gvalue.Type())
		ptr := reflect.New(pto.Elem())
		ptr.Elem().Set(gvalue)
		return ptr
	}
	return gvalue
}

func CheckImplements(atype reflect.Type, btype reflect.Type) bool {
	if btype.Kind() == reflect.Interface {
		return atype.Implements(btype)
	} else if atype.Kind() == reflect.Interface {
		return btype.Implements(atype)
	}
	return false
}

func CheckImplementsOneOf(atype reflect.Type, types ...reflect.Type) bool {
	for _, tx := range types {
		if CheckImplements(atype, tx) {
			return true
		}
	}
	return false
}

func CompareType(atype reflect.Type, btype reflect.Type) bool {
	return atype == btype || CheckImplements(atype, btype)
}

func CheckIn(ftype reflect.Type, params ...reflect.Type) bool {
	if ftype.NumIn() == len(params) {
		for i, ttype := range params {
			if !CompareType(ftype.In(i), ttype) {
				return false
			}
		}
		return true
	}
	return false
}

func CheckOut(ftype reflect.Type, params ...reflect.Type) bool {
	if ftype.NumOut() == len(params) {
		for i, ttype := range params {
			if !CompareType(ftype.Out(i), ttype) {
				return false
			}
		}
		return true
	}
	return false
}

func structPtr(v interface{}) (reflect.Value, bool) {
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Ptr:
		if val.Elem().Kind() == reflect.Struct {
			return val, true
		}
	case reflect.Struct:
		return PtrValue(val), true
	}
	return reflect.ValueOf(nil), false
}

func MethodByName(v interface{}, name string) (reflect.Value, bool) {
	if val, ok  := structPtr(v); ok {
		if mv := val.MethodByName(name); mv.IsValid() {
			return mv, true
		}
	}
	return reflect.ValueOf(nil), false
}

func PreStartMethod(v interface{}) (reflect.Value, bool) {
	if mv, ok := MethodByName(v, "PreStart"); ok {
		mt := mv.Type()
		if CheckIn(mt, ActorContextType) && CheckOut(mt, ErrorType) {
			return mv, true
		}
	}
	return reflect.ValueOf(nil), false
}

func PreStartFunc(v interface{}) func(ActorContext) error {
	if val, ok := PreStartMethod(v); ok {
		return func(ctx ActorContext) error {
			result := val.Call([]reflect.Value{reflect.ValueOf(ctx)})
			if len(result) > 0 {
				if err, ok := result[0].Interface().(error); ok {
					return err
				}
			}
			return nil
		}
	}
	return func(ac ActorContext) error { return nil }
}

func PostStopMethod(v interface{}) (reflect.Value, bool) {
	if mv, ok := MethodByName(v, "PostStop"); ok {
		mt := mv.Type()
		if CheckIn(mt, ActorContextType) && CheckOut(mt, ErrorType) {
			return mv, true
		}
	}
	return reflect.ValueOf(nil), false
}

func PostStopFunc(v interface{}) func(ActorContext) error {
	if val, ok := PostStopMethod(v); ok {
		return func(ctx ActorContext) error {
			result := val.Call([]reflect.Value{reflect.ValueOf(ctx)})
			if len(result) > 0 {
				if err, ok := result[0].Interface().(error); ok {
					return err
				}
			}
			return nil
		}
	}
	return func(ac ActorContext) error { return nil }
}

func ReceiveMethod(v interface{}) (reflect.Value, bool) {
	if mv, ok := MethodByName(v, "Receive"); ok {
		mt := mv.Type()
		if CheckIn(mt, ActorContextType, MessageType) && mt.NumOut() == 0 {
			return mv, true
		}
	}
	return reflect.ValueOf(nil), false
}

func ReceiveFunc(v interface{}) func(ActorContext, Message) {
	if val, ok := ReceiveMethod(v); ok {
		return func(ctx ActorContext, msg Message) {
			val.Call([]reflect.Value{ reflect.ValueOf(ctx), reflect.ValueOf(msg) })
		}
	}
	return func(ctx ActorContext, msg Message) { msg.Reply(ErrUnknownCommand) }  // TODO maybe reply with done?
}
