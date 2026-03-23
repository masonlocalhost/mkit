package sqlutil

import (
	"fmt"
	"strconv"
	"strings"
)

func EqualClause(col string) string {
	return fmt.Sprintf("%s = ?", col)
}

func EqualClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s = ?", table, col)
}

func MultipleEqualClause(cols []string, op LogicOperator) string {
	var clauses []string
	for _, col := range cols {
		clauses = append(clauses, fmt.Sprintf("%s = ?", col))
	}

	return strings.Join(clauses, op.String())
}

func IsNotNullClause(col string) string {
	return fmt.Sprintf("%s IS NOT NULL", col)
}

func IsNotNullClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s IS NOT NULL", table, col)
}

func BoolClause(col string, boolV bool) string {
	return fmt.Sprintf("%s = %s", col, strconv.FormatBool(boolV))
}

func InClause(col string) string {
	return fmt.Sprintf("%s IN ?", col)
}

func NotInClause(col string) string {
	return fmt.Sprintf("%s NOT IN ?", col)
}

func NotInClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s NOT IN ?", table, col)
}

func BetweenClause(col string) string {
	return fmt.Sprintf("%s BETWEEN ? and ?", col)
}

func BetweenClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s BETWEEN ? and ?", table, col)
}

func NotEqualClause(col string) string {
	return fmt.Sprintf("%s <> ?", col)
}

func NotEqualClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s <> ?", table, col)
}

func ConcatClauses(clauses []*ConcatClause) string {
	var res string
	for _, c := range clauses {
		res += fmt.Sprintf("%s%s", c.Clause, c.Operator)
	}

	return res
}

func SumSelect(col string, as string) string {
	if as != "" {
		return fmt.Sprintf("sum(%s) as %s", col, as)
	}

	return fmt.Sprintf("sum(%s)", col)
}

func GroupConcatClause(col string, sep string, as string) string {
	if as != "" {
		return fmt.Sprintf("group_concat(%s SEPARATOR '%s') as %s", col, sep, as)
	}

	return fmt.Sprintf("group_concat(%s SEPARATOR %s)", col, sep)
}

func GroupConcatDistinctClause(col string, sep string, as string) string {
	if as != "" {
		return fmt.Sprintf("group_concat(DISTINCT %s SEPARATOR '%s') as %s", col, sep, as)
	}

	return fmt.Sprintf("group_concat(DISTINCT %s SEPARATOR %s)", col, sep)
}

func MaxClause(col string) string {
	return fmt.Sprintf("IFNULL(max(%s), 0)", col)
}

func MaxAndAddSelect(col string, n uint64) string {
	return fmt.Sprintf("IFNULL(max(%s), 0) + %v", col, n)
}

func IsNullClause(col string) string {
	return fmt.Sprintf("%s IS NULL", col)
}

func IsNullClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s IS NULL", table, col)
}

func LessThanOrEqualClause(col string) string {
	return fmt.Sprintf("%s <= ?", col)
}

func LessThanOrEqualClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s <= ?", table, col)
}

func CountSelect(col string, as string) string {
	if as != "" {
		return fmt.Sprintf("count(%s) as %s", col, as)
	}

	return fmt.Sprintf("count(%s)", col)
}

func JSONStringArrayContains(col string) string {
	return fmt.Sprintf("JSON_CONTAINS(%s, '\"%s\"', '$') = ?", col, col)
}

func GreaterThanOrEqualClause(col string) string {
	return fmt.Sprintf("%s >= ?", col)
}

func GreaterThanOrEqualClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s >= ?", table, col)
}

func LessThanClause(col string) string {
	return fmt.Sprintf("%s < ?", col)
}

func GreaterThanClause(col string) string {
	return fmt.Sprintf("%s > ?", col)
}

func LessThanClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s < ?", table, col)
}

func GreaterThanClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s > ?", table, col)
}

func LikeClause(col string) string {
	return fmt.Sprintf("%s LIKE ?", col)
}

func ILikeClause(col string) string {
	return fmt.Sprintf("%s ILIKE ?", col)
}

func InClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s IN ?", table, col)
}

func LikeClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s LIKE ?", table, col)
}

func ILikeClauseWithTable(table, col string) string {
	return fmt.Sprintf("%s.%s ILIKE ?", table, col)
}

func ExistClauseWithTable(subquery string) string {
	return fmt.Sprintf("EXISTS (%s)", subquery)
}

func ExistClause() string {
	return fmt.Sprintf("EXISTS (?)")
}

func NotExistClause() string {
	return fmt.Sprintf("NOT EXISTS (?)")
}
