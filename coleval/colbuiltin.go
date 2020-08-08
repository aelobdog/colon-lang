package coleval

import (
	obj "colon/colobj"
	"fmt"
)

var builtin = map[string]*obj.BuiltIn{
	"len": {
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
