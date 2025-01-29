package utils

func HasPermission(recordUserID string, userID string) bool {
	return recordUserID == userID
}
