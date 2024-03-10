package gdp

type Query struct {
	query string
	i     int
	len   int
}

type QueryAttr struct {
	query string
	i     int
	len   int
}

func (q *Query) parseQuery(splited *string) bool {
	if q.i == q.len {
		return false
	}

	c := q.query[q.i]

	if c == ' ' || c == ',' || c == '>' {
		q.i++
		*splited = string(c)
		return true
	}

	a := ""
	for {
		c := q.query[q.i]
		q.i++

		if c == ' ' || c == ',' || c == '>' {
			q.i--
			break
		}
		a += string(c)

		if q.i == q.len {
			break
		}
	}
	*splited = a
	return true
}

func (q *QueryAttr) getAttr() string {
	a := ""
	for {
		c := q.query[q.i]
		q.i++
		if c == '\'' || c == '"' {
			q.i++
			break
		}

		if c == '#' || c == '.' || c == '[' || c == '=' || c == ']' {
			break
		}

		a += string(c)

		if q.i >= q.len {
			break
		}
	}

	return a
}

func (q *QueryAttr) parseAttr() map[string]string {
	ret := make(map[string]string)

	for {
		c := q.query[q.i]

		q.i++

		if c == '.' {
			ret["class"] = q.getAttr()
		} else if c == '#' {
			ret["id"] = q.getAttr()

		} else if c == ']' {
			q.i++
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

		if q.i >= q.len {
			break
		}

	}

	return ret
}
