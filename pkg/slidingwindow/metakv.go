package slidingwindow

type kv struct {
	key          string
	defaultValue int

	next *kv
}

type Kv struct {
	head   *kv
	tail   *kv
	length int
}

func (list *Kv) setDefault(key string, value int) {
	if list.head == nil && list.tail == nil && list.length == 0 {
		list.head = &kv{key: key, defaultValue: value, next: nil}
		list.tail = list.head
		list.length = 1
		return
	}

	head := list.head
	for head != nil {
		if head.key == key {
			head.defaultValue = value
			return
		}
		head = head.next
	}

	list.tail.next = &kv{key: key, defaultValue: value, next: nil}
	list.tail = list.tail.next
	list.length++
}

func (list *Kv) remove(key string) int {
	if list.head.key == key {
		list.head = nil
		if list.length == 1 {
			list.tail = nil
			list.length = 0
		}
		return 0
	}

	var prev *kv = nil
	var head *kv = list.head

	for head != nil {
		if head.key == key {
			prev.next = head.next
			head = nil
			return 0
		}
		prev = head
		head = head.next
	}

	// return -1 means not found the key in key value list.
	return -1
}
