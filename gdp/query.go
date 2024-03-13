package gdp

import (
	"bytes"
	"strings"
)

type queryAttr struct {
	query string
	i     int
	len   int
}

func splitQueries(mainQuery string) []string {
	var ret []string
	var buffer bytes.Buffer
	for _, char := range mainQuery {
		if char == ',' {
			ret = append(ret, strings.Trim(buffer.String(), " "))
			buffer = bytes.Buffer{}
			continue
		}
		buffer.WriteRune(char)
	}
	ret = append(ret, strings.Trim(buffer.String(), " "))
	return ret
}

func skipSpace(query *string, i *int, len int) {
	for *i < len {
		c1 := (*query)[*i]
		if c1 == ' ' {
			*i++
		} else {
			break
		}
	}
}

func splitQuery(query string) []string {
	var ret []string
	var buffer bytes.Buffer
	len := len(query)
	i := 0
	for i < len {
		c := query[i]
		i++

		if c == '>' || c == '+' || c == '~' || c == '|' {
			ret = append(ret, strings.Trim(buffer.String(), " "))
			ret = append(ret, string(c))

			skipSpace(&query, &i, len)
			buffer = bytes.Buffer{}
			continue
		}

		if c == ' ' {
			skipSpace(&query, &i, len)
			c = query[i]

			ret = append(ret, strings.Trim(buffer.String(), " "))
			if c == '>' || c == '+' || c == '~' || c == '|' {
				ret = append(ret, string(c))
				i++
			} else {
				ret = append(ret, " ")
			}

			skipSpace(&query, &i, len)

			buffer = bytes.Buffer{}
			continue
		}
		buffer.WriteByte(c)
	}
	ret = append(ret, strings.Trim(buffer.String(), " "))
	return ret
}

func (q *queryAttr) getAttr() string {
	var buffer bytes.Buffer
	for q.i < q.len {
		c := q.query[q.i]
		q.i++
		if c == '\'' || c == '"' {
			continue
		}

		if c == '#' || c == '.' || c == '[' || c == '=' || c == ']' || c == ':' || c == '(' {
			break
		}
		buffer.WriteByte(c)
	}

	return strings.Trim(buffer.String(), " ")
}

func (q *queryAttr) getParenthesis() string {
	var buffer bytes.Buffer
	for q.i < q.len {
		c := q.query[q.i]
		q.i++
		if c == '(' {
			continue
		}

		if c == ')' {
			break
		}

		buffer.WriteByte(c)
	}
	return strings.Trim(buffer.String(), " ")
}

func (q *queryAttr) parseAttr() map[string]string {
	ret := make(map[string]string)

	for q.i < q.len {
		c := q.query[q.i]

		q.i++

		if c == '.' {
			ret["class"] = q.getAttr()
		} else if c == '#' {
			ret["id"] = q.getAttr()

		} else if c == ']' {
			q.i++
		} else if c == ':' {
			key := q.getAttr()
			ret[key] = q.getParenthesis()
		} else if c == '[' {

			ret[q.getAttr()] = q.getAttr()

		} else {

			q.i--
			ret["tag"] = q.getAttr()
			if q.i >= q.len {
				break
			}
			q.i--
		}

	}

	return ret
}
