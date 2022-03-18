package consts

var Values = map[string]string{}

type Const struct {
	Name  string
	Value string
}

func Init() {
	var e, pi Const
	e.Name = "e"
	pi.Name = "pi"

	e.Value = "2.71828182845904523536028747135266249775724709369995957496696763"
	pi.Value = "3.14159265358979323846264338327950288419716939937510582097494459"

	Values[e.Name] = e.Value
	Values[pi.Name] = pi.Value
}
