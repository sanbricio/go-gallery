package emailTemplate

import "fmt"

type RecoveryTemplate struct{}

func (t RecoveryTemplate) Subject() string {
	return "🔑 Recovery code to reset your go-gallery password"
}

func (t RecoveryTemplate) Body(code string, email string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Password Recovery</title>
		</head>
		<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
			<table align="center" width="100%%" bgcolor="#ffffff" style="max-width: 600px; padding: 20px; border-radius: 8px; box-shadow: 0px 0px 10px #cccccc;">
				<tr>
					<td align="center" style="padding-bottom: 20px;">
						<h2 style="color: #333;">🔐 Password Recovery Code</h2>
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 16px; color: #555;">
						You have requested to reset your Go Gallery account password. Use the following code to proceed:
					</td>
				</tr>
				<tr>
					<td align="center" style="padding: 20px;">
						<input type="text" value="%s" id="recoveryCode" readonly style="text-align: center; width: 250px; font-size: 24px; padding: 12px; border: 1px solid #ddd; border-radius: 5px; background-color: #f8f8f8; color: #333;" />
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 14px; color: #777; padding-top: 20px;">
						⚠️ This code is valid for <strong>5 minutes</strong>. Do not share it with anyone.
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 14px; color: #777; padding-top: 10px;">
						If you didn't request a password reset, please ignore this message. Your account remains secure.
					</td>
				</tr>
				<tr>
					<td align="center" style="padding-top: 30px; font-size: 12px; color: #aaa;">
						Best regards, <br>
						<strong>Go Gallery Support Team</strong><br>
						<a href="mailto:gogalleryteam@gmail.com" style="color: #3498db; text-decoration: none;">gogalleryteam@gmail.com</a>
					</td>
				</tr>
			</table>
		</body>
		</html>`, code)
}
