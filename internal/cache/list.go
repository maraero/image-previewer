package cache

type List interface {
	length() int
	front() *listItem
	back() *listItem
	pushFront(v any) *listItem
	pushBack(v any) *listItem
	remove(i *listItem)
	moveToFront(i *listItem)
}

type listItem struct {
	Value any
	Next  *listItem
	Prev  *listItem
}

type lst struct {
	f   *listItem
	b   *listItem
	len int
}

func newList() List {
	return new(lst)
}

func (l *lst) length() int {
	return l.len
}

func (l *lst) front() *listItem {
	return l.f
}

func (l *lst) back() *listItem {
	return l.b
}

func (l *lst) pushFront(v any) *listItem {
	elem := &listItem{Value: v}

	switch l.length() {
	case 0:
		l.f = elem
		l.b = elem
	case 1:
		l.f = elem
		l.f.Next = l.b
		l.b.Prev = l.f
	default:
		elem.Next = l.f
		l.f.Prev = elem
		l.f = elem
	}

	l.len++
	return elem
}

func (l *lst) pushBack(v any) *listItem {
	elem := &listItem{Value: v}

	switch l.length() {
	case 0:
		l.f = elem
		l.b = elem
	case 1:
		l.b = elem
		l.b.Prev = l.f
		l.f.Next = l.b
	default:
		elem.Prev = l.b
		l.b.Next = elem
		l.b = elem
	}

	l.len++
	return elem
}

func (l *lst) remove(i *listItem) {
	switch {
	case l.length() == 1:
		l.f = nil
		l.b = nil
	case i == l.f:
		l.f = l.f.Next
		l.f.Prev = nil
	case i == l.b:
		l.b = l.b.Prev
		l.b.Next = nil
	default:
		i.Prev.Next = i.Next
		i.Next.Prev = i.Prev
	}

	l.len--
}

func (l *lst) moveToFront(i *listItem) {
	l.remove(i)
	l.pushFront(i.Value)
}
