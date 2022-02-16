package beercss

type IconType string

const (
	IconLocalBar  IconType = `local_bar`
	IconHistory   IconType = `history`
	IconAdd       IconType = `add`
	IconAddCircle IconType = `add_circle`
	IconList      IconType = `list`
	IconPerson    IconType = `person`
	IconSwapHoriz IconType = `swap_horiz`
	IconPayments  IconType = `payments`
	IconDelete    IconType = `delete`
	IconRemove    IconType = `remove`
	IconFastfood  IconType = `fastfood`
)

type Icon struct {
	Icon IconType
}
