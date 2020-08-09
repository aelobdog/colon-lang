package coleval

import (
	"bufio"
	obj "colon/colobj"
	"fmt"
	"os"
)

var builtin = map[string]*obj.BuiltIn{
	"len": {
		/*
			use: len(array)
		*/
		Bfunct: func(args ...obj.Object) obj.Object {
			if len(args) != 1 {
				reportRuntimeError(fmt.Sprintf("len takes only 1 argument, got %v", len(args)))
			}
			switch arg := args[0].(type) {
			case *obj.String:
				return &obj.Integer{
					Value: int64(len(arg.Value)),
				}
			case *obj.List:
				return &obj.Integer{
					Value: int64(len(arg.Elements)),
				}
			default:
				reportRuntimeError(fmt.Sprintf("len cannot operate of type %q.", args[0].ObType()))
			}
			return nil
		},
	},

	"print": {
		/*
			use: print(stuff_to_print)
		*/
		Bfunct: func(args ...obj.Object) obj.Object {
			if len(args) == 0 {
				fmt.Println()
			}
			for _, v := range args {
				fmt.Println(v.ObValue())
			}
			return EMPTY
		},
	},

	"head": {
		Bfunct: func(args ...obj.Object) obj.Object {
			if len(args) != 1 {
				reportRuntimeError(fmt.Sprintf("head takes only 1 argument, got %v", len(args)))
			}
			switch arg := args[0].(type) {
			case *obj.String:
				if len(arg.Value) < 1 {
					reportRuntimeError("cannot take the head of an empty string")
				}
				return &obj.String{
					Value: string(arg.Value[0]),
				}
			case *obj.List:
				if len(arg.Elements) < 1 {
					reportRuntimeError("cannot take the head of an empty list")
				}
				return arg.Elements[0]
			default:
				reportRuntimeError(fmt.Sprintf("head cannot operate of type %q.", args[0].ObType()))
			}
			return nil
		},
	},

	"last": {
		Bfunct: func(args ...obj.Object) obj.Object {
			if len(args) != 1 {
				reportRuntimeError(fmt.Sprintf("last takes only 1 argument, got %v", len(args)))
			}
			switch arg := args[0].(type) {
			case *obj.String:
				if len(arg.Value) < 1 {
					reportRuntimeError("cannot take the last of an empty string")
				}
				return &obj.String{
					Value: string(arg.Value[len(arg.Value)-1]),
				}
			case *obj.List:
				if len(arg.Elements) < 1 {
					reportRuntimeError("cannot take the last of an empty list")
				}
				return arg.Elements[len(arg.Elements)-1]
			default:
				reportRuntimeError(fmt.Sprintf("last cannot operate of type %q.", args[0].ObType()))
			}
			return nil
		},
	},

	"tail": {
		Bfunct: func(args ...obj.Object) obj.Object {
			if len(args) != 1 {
				reportRuntimeError(fmt.Sprintf("tail takes only 1 argument, got %v", len(args)))
			}
			switch arg := args[0].(type) {
			case *obj.List:
				if len(arg.Elements) < 1 {
					reportRuntimeError("cannot take the tail of an empty list")
				}
				tail := obj.List{
					Elements: []obj.Object{},
				}
				if len(arg.Elements) == 1 {
					return &tail
				}
				for k := 1; k < len((arg.Elements)); k++ {
					tail.Elements = append(tail.Elements, arg.Elements[k])
				}
				return &tail
			default:
				reportRuntimeError(fmt.Sprintf("tail cannot operate of type %q.", args[0].ObType()))
			}
			return nil
		},
	},

	"init": {
		Bfunct: func(args ...obj.Object) obj.Object {
			if len(args) != 1 {
				reportRuntimeError(fmt.Sprintf("init takes only 1 argument, got %v", len(args)))
			}
			switch arg := args[0].(type) {
			case *obj.List:
				if len(arg.Elements) < 1 {
					reportRuntimeError("cannot take the init of an empty list")
				}
				init := obj.List{
					Elements: []obj.Object{},
				}
				if len(arg.Elements) == 1 {
					return &init
				}
				for k := 0; k < len((arg.Elements))-1; k++ {
					init.Elements = append(init.Elements, arg.Elements[k])
				}
				return &init
			default:
				reportRuntimeError(fmt.Sprintf("init cannot operate of type %q.", args[0].ObType()))
			}
			return nil
		},
	},

	"isNull": {
		Bfunct: func(args ...obj.Object) obj.Object {
			if len(args) != 1 {
				reportRuntimeError(fmt.Sprintf("isNull takes only 1 argument, got %v", len(args)))
			}
			switch arg := args[0].(type) {
			case *obj.List:
				if len(arg.Elements) == 0 {
					return &obj.Boolean{Value: true}
				}
				return &obj.Boolean{Value: false}
			default:
				reportRuntimeError(fmt.Sprintf("isNull cannot operate of type %q.", args[0].ObType()))
			}
			return nil
		},
	},

	"push": {
		Bfunct: func(args ...obj.Object) obj.Object {
			if len(args) == 1 {
				reportRuntimeError("push takes a minimum of 3 arguments: a list and an element")
			}
			switch arg := args[0].(type) {
			case *obj.List:
				for k := 1; k < len(args); k++ {
					if _, ok := args[k].(*obj.Integer); ok {
						arg.Elements = append(arg.Elements, args[k].(*obj.Integer))
					} else if _, ok := args[k].(*obj.Floating); ok {
						arg.Elements = append(arg.Elements, args[k].(*obj.Floating))
					} else if _, ok := args[k].(*obj.Boolean); ok {
						arg.Elements = append(arg.Elements, args[k].(*obj.Boolean))
					} else if _, ok := args[k].(*obj.String); ok {
						arg.Elements = append(arg.Elements, args[k].(*obj.String))
					} else if _, ok := args[k].(*obj.List); ok {
						arg.Elements = append(arg.Elements, args[k].(*obj.List))
					} else {
						reportRuntimeError(fmt.Sprintf("cannot push element of type %q into a list", args[k].ObType()))
					}
					// check if adding arrays is possible
				}
				return EMPTY
			default:
				reportRuntimeError(fmt.Sprintf("push cannot operate of type %q.", args[0].ObType()))
			}
			return nil
		},
	},
}

var builtinTypeAssociations = map[string]obj.Object{
	"int":  &obj.DataType{Dtype: "integer"},
	"flt":  &obj.DataType{Dtype: "float"},
	"bool": &obj.DataType{Dtype: "boolean"},
	"str":  &obj.DataType{Dtype: "string"},
}

// GetInput : function that gets input and binds it to an name
func GetInput(env *obj.Env, varname string, dtype obj.DataType) obj.Object {
	switch dtype.Dtype {
	case "integer":
		var val int64
		fmt.Scanf("%d\n", &val)
		env.Set(varname, &obj.Integer{Value: int64(val)})
	case "float":
		var val float64
		fmt.Scanf("%f\n", &val)
		env.Set(varname, &obj.Floating{Value: float64(val)})
	case "boolean":
		var val bool
		fmt.Scanf("%t\n", &val)
		env.Set(varname, &obj.Boolean{Value: bool(val)})
	case "string":
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		env.Set(varname, &obj.String{Value: text})
	}
	return EMPTY
}
