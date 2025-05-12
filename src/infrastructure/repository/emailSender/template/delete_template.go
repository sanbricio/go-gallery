package emailTemplate

import "fmt"

type DeleteAccountTemplate struct{}

func (t DeleteAccountTemplate) Subject() string {
	return "⚠️ Verification code to delete your go-gallery account"
}

func (t DeleteAccountTemplate) Body(code string, email string) string {
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<meta charset="UTF-8">
			<title>Account Deletion Confirmation</title>
		</head>
		<body style="font-family: Arial, sans-serif; background-color: #f4f4f4; padding: 20px;">
			<table align="center" width="100%%" bgcolor="#ffffff" style="max-width: 600px; padding: 20px; border-radius: 8px; box-shadow: 0px 0px 10px #cccccc;">
				<tr>
					<td align="center" style="padding-bottom: 20px;">
						<h2 style="color: #333;">🔐 Verification Code</h2>
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 16px; color: #555;">
						We received a request to delete your account. To confirm this action, please use the following verification code:
					</td>
				</tr>
				<tr>
					<td align="center" style="padding: 20px;">
						<div style="display: inline-block; background-color: #f8f8f8; padding: 15px 30px; border-radius: 5px; font-size: 24px; font-weight: bold; color: #333; border: 1px solid #ddd;">
							%s
						</div>
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 14px; color: #777; padding-top: 20px;">
						⚠️ This code is valid only for the next <strong>5 minutes</strong>. Do not share it with anyone.
					</td>
				</tr>
				<tr>
					<td align="center" style="font-size: 14px; color: #777; padding-top: 10px;">
						If you did not request the deletion of your account, please ignore this message and your account will remain active.
					</td>
				</tr>
				<tr>
					<td align="center" style="padding-top: 30px; font-size: 12px; color: #aaa;">
						Best regards, <br>
						<strong>Support Team</strong><br>
						<a href="mailto:gogalleryteam@gmail.com" style="color: #3498db; text-decoration: none;">gogalleryteam@gmail.com</a>
					</td>
				</tr>
			</table>
		</body>
		</html>`, code)
}
