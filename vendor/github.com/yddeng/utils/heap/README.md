# heap
heap 

## Usage

```
h := New()
h.Push(&HElement{name: "e1", value: 5})
h.Push(&HElement{name: "e2", value: 12})
h.Push(&HElement{name: "e3", value: 8})
e4 := &HElement{name: "e4", value: 7}
h.Push(e4)
fmt.Println(h.Top())

e4.value = 2
h.Fix(e4)
fmt.Println(h.Top(), h.IsExist(e4))

h.Remove(e4)
fmt.Println(h.Top(), h.IsExist(e4))

fmt.Println(h.Pop())
```