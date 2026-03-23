package sqlutil

const (
	LogicOperator_NONE LogicOperator = ""
	LogicOperator_AND  LogicOperator = " AND "
	LogicOperator_OR   LogicOperator = " OR "
)

type LogicOperator string

func (l LogicOperator) String() string {
	return string(l)
}

type ConcatClause struct {
	Clause   string
	Operator LogicOperator
}

type SortValue string

const (
	SortValue_ASC  = "asc"
	SortValue_DESC = "desc"
)

type SortItem struct {
	Field     string
	SortValue SortValue
}
