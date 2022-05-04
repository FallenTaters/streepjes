package document

import "syscall/js"

func DeleteCookie(name string) {
	js.Global().Get(`document`).Set(`cookie`, name+`=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`)
}

func GetElementById(id string) (js.Value, bool) {
	elem := js.Global().Get(`document`).Call(`getElementById`, id)

	return elem, elem.Type() == js.TypeObject
}
