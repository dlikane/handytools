package pinterest

func getColumnsForLayout(layout string) int {
	switch layout {
	case "compact":
		return 6
	case "medium":
		return 3
	case "large":
		return 1
	case "fit":
		return 1 // auto-fit, one column at a time
	default:
		return -1
	}
}
