package document

import "syscall/js"

func DeleteCookie(name string) {
	js.Global().Get(`document`).Set(`cookie`, name+`=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`)
}
