package forms

type errors map[string][]string

func (e errors) Len() int { return len(e) }

func (e errors) Add(field, message string) errors {
	if _, ok := e[field]; !ok {
		e[field] = make([]string, 0)

	}
	e[field] = append(e[field], message)
	return e
}

func (e errors) Get(field string) string {
	es := e[field]
	if len(es) != 0 {
		return es[0]
	}
	return ""
}
