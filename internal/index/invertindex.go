package index

import "container/list"

type InvertIndex map[string]*list.List

func Cup(l1 *list.List, l2 *list.List) *list.List {
	result := list.New()
	used := make(map[DocID]bool)

	e1 := l1.Front()
	e2 := l2.Front()

	for e1 != nil && e2 != nil {
		if e1.Value.(DocID) < e2.Value.(DocID) {
			if _, ok := used[e1.Value.(DocID)]; !ok {
				result.PushBack(e1.Value)

				used[e1.Value.(DocID)] = true
			}

			e1 = e1.Next()
		} else {
			if _, ok := used[e2.Value.(DocID)]; !ok {
				result.PushBack(e2.Value)

				used[e2.Value.(DocID)] = true
			}

			e2 = e2.Next()
		}
	}

	for e1 != nil {
		if _, ok := used[e1.Value.(DocID)]; !ok {
			result.PushBack(e1.Value)

			used[e1.Value.(DocID)] = true
		}

		e1 = e1.Next()
	}

	for e2 != nil {
		if _, ok := used[e2.Value.(DocID)]; !ok {
			result.PushBack(e2.Value)

			used[e2.Value.(DocID)] = true
		}

		e2 = e2.Next()
	}

	return result
}

func Cap(l1 *list.List, l2 *list.List) *list.List {
	result := list.New()

	e1 := l1.Front()
	e2 := l2.Front()

	for e1 != nil && e2 != nil {
		switch {
		case e1.Value.(DocID) == e2.Value.(DocID):
			result.PushBack(e1.Value)

			e1 = e1.Next()
			e2 = e2.Next()
		case e1.Value.(DocID) < e2.Value.(DocID):
			e1 = e1.Next()
		default:
			e2 = e2.Next()
		}
	}

	return result
}

func Invert(l *list.List, all *list.List) *list.List {
	result := list.New()

	e := l.Front()
	a := all.Front()

	for a != nil && e != nil {
		if a.Value.(DocID) != e.Value.(DocID) {
			result.PushBack(a.Value)

			a = a.Next()
		} else {
			e = e.Next()
			a = a.Next()
		}
	}

	for a != nil {
		result.PushBack(a.Value)

		a = a.Next()
	}

	return result
}
