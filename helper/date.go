package helper

func IntMountToText(mount int) string {
	var text string
	switch mount {
	case 1:
		text = "Januari"
	case 2:
		text = "Februari"
	case 3:
		text = "Maret"
	case 4:
		text = "April"
	case 5:
		text = "Mei"
	case 6:
		text = "Juni"
	case 7:
		text = "Juli"
	case 8:
		text = "Agustus"
	case 9:
		text = "September"
	case 10:
		text = "Oktober"
	case 11:
		text = "November"
	case 12:
		text = "Desember"
	}
	return text
}
