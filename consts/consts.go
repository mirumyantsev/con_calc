package consts

var Values = map[string]string{}

var constants = map[string]string {
	"e": "2.71828182845904523536028747135266249775724709369995957496696763",
	"pi": "3.14159265358979323846264338327950288419716939937510582097494459",
	"phi": "1.61803398874989484820458683436563811772030917980576286213544862",
}

type Const struct {
	Name  string
	Value string
}

func Init() {
	for val, key := range constants {
		newConst := Const{
			Name: val,
			Value: key,
		}
		Values[newConst.Name] = newConst.Value
	}
}
