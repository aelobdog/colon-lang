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
		fmt.Scanf("%d", &val)
		env.Set(varname, &obj.Integer{Value: int64(val)})
	case "float":
		var val float64
		fmt.Scanf("%f", &val)
		env.Set(varname, &obj.Floating{Value: float64(val)})
	case "boolean":
		var val bool
		fmt.Scanf("%t", &val)
		env.Set(varname, &obj.Boolean{Value: bool(val)})
	case "string":
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		env.Set(varname, &obj.String{Value: text})
	}
	return EMPTY
}
