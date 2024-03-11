package gdp

import "strings"

type QueryAttr struct {
	query string
	i     int
	len   int
}

func splitQueries(mainQuery string) []string {
	var ret []string
	query := ""
	for _, char := range mainQuery {
		if char == ',' {
			ret = append(ret, strings.Trim(query, " "))
			query = ""
			continue
		}
		query += string(char)
	}
	ret = append(ret, strings.Trim(query, " "))
	return ret
}

func splitQuery(query string) []string {
	var ret []string
	str := ""
	for _, c := range query {

		if c == ' ' || c == '>' {

			ret = append(ret, strings.Trim(str, " "))
			ret = append(ret, string(c))

			str = ""
			continue
		}
		str += string(c)
	}
	ret = append(ret, strings.Trim(str, " "))
	return ret
}

func (q *QueryAttr) getAttr() string {
	a := ""
	for q.i < q.len {
		c := q.query[q.i]
		q.i++
		if c == '\'' || c == '"' {
			q.i++
			break
		}

		if c == '#' || c == '.' || c == '[' || c == '=' || c == ']' || c == ':' || c == '(' {
			break
		}

		a += string(c)
	}

	return a
}

func (q *QueryAttr) getParenthesis() string {
	ret := ""
	for q.i < q.len {
		c := q.query[q.i]
		q.i++
		if c == '(' {
			continue
		}

		if c == ')' {
			break
		}

		ret += string(c)
	}
	return ret
}

func (q *QueryAttr) parseAttr() map[string]string {
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

			key := q.getAttr()
			q.i++

			ret[key] = q.getAttr()
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
