package utils

func IsValidExtension(extension string) bool {
	validExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	for _, validExt := range validExtensions {
		if extension == validExt {
			return true
		}
	}
	return false
}
