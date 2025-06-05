package mail

// 辅助函数
func jsonEmails(emails []string) string {
	result := ""
	for i, email := range emails {
		if i > 0 {
			result += ", "
		}
		result += email
	}
	return result
}
