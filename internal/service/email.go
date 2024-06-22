package service

import (
	"strconv"
	"time"

	"github.com/Giafn/Depublic/pkg/background_job"
)

func ScheduleEmails(emailAddresses string, subject, body string) {
	background_job.ScheduleEmails(
		emailAddresses,
		subject,
		body,
	)
}

func CreateConfirmationAccountEmailHtml(confirmationUrl string, userName string) string {
	return `<!DOCTYPE html>
	<html>
	<head>
		<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title></title>
	
		<!--[if !mso]><!-->
		<style type="text/css">
			@import url('https://fonts.mailersend.com/css?family=Inter:400,600');
		</style>
		<!--<![endif]-->
	
		<style type="text/css" rel="stylesheet" media="all">
			@media only screen and (max-width: 640px) {
				.ms-header {
					display: none !important;
				}
				.ms-content {
					width: 100% !important;
					border-radius: 0;
				}
				.ms-content-body {
					padding: 30px !important;
				}
				.ms-footer {
					width: 100% !important;
				}
				.mobile-wide {
					width: 100% !important;
				}
				.info-lg {
					padding: 30px;
				}
			}
		</style>
		<!--[if mso]>
		<style type="text/css">
		body { font-family: Arial, Helvetica, sans-serif!important  !important; }
		td { font-family: Arial, Helvetica, sans-serif!important  !important; }
		td * { font-family: Arial, Helvetica, sans-serif!important  !important; }
		td p { font-family: Arial, Helvetica, sans-serif!important  !important; }
		td a { font-family: Arial, Helvetica, sans-serif!important  !important; }
		td span { font-family: Arial, Helvetica, sans-serif!important  !important; }
		td div { font-family: Arial, Helvetica, sans-serif!important  !important; }
		td ul li { font-family: Arial, Helvetica, sans-serif!important  !important; }
		td ol li { font-family: Arial, Helvetica, sans-serif!important  !important; }
		td blockquote { font-family: Arial, Helvetica, sans-serif!important  !important; }
		th * { font-family: Arial, Helvetica, sans-serif!important  !important; }
		</style>
		<![endif]-->
	</head>
	<body style="font-family:'Inter', Helvetica, Arial, sans-serif; width: 100% !important; height: 100%; margin: 0; padding: 0; -webkit-text-size-adjust: none; background-color: #f4f7fa; color: #4a5566;" >
	
	<div class="preheader" style="display:none !important;visibility:hidden;mso-hide:all;font-size:1px;line-height:1px;max-height:0;max-width:0;opacity:0;overflow:hidden;" ></div>
	
	<table class="ms-body" width="100%" cellpadding="0" cellspacing="0" role="presentation" style="border-collapse:collapse;background-color:#f4f7fa;width:100%;margin-top:0;margin-bottom:0;margin-right:0;margin-left:0;padding-top:0;padding-bottom:0;padding-right:0;padding-left:0;" >
		<tr>
			<td align="center" style="word-break:break-word;font-family:'Inter', Helvetica, Arial, sans-serif;font-size:16px;line-height:24px;" >
	
				<table class="ms-container" width="100%" cellpadding="0" cellspacing="0" style="border-collapse:collapse;width:100%;margin-top:0;margin-bottom:0;margin-right:0;margin-left:0;padding-top:0;padding-bottom:0;padding-right:0;padding-left:0;" >
					<tr>
						<td align="center" style="word-break:break-word;font-family:'Inter', Helvetica, Arial, sans-serif;font-size:16px;line-height:24px;" >
	
							<table class="ms-header" width="100%" cellpadding="0" cellspacing="0" style="border-collapse:collapse;" >
								<tr>
									<td height="40" style="font-size:0px;line-height:0px;word-break:break-word;font-family:'Inter', Helvetica, Arial, sans-serif;" >
										&nbsp;
									</td>
								</tr>
							</table>
	
						</td>
					</tr>
					<tr>
						<td align="center" style="word-break:break-word;font-family:'Inter', Helvetica, Arial, sans-serif;font-size:16px;line-height:24px;" >
	
							<table class="ms-content" width="640" cellpadding="0" cellspacing="0" role="presentation" style="border-collapse:collapse;width:640px;margin-top:0;margin-bottom:0;margin-right:auto;margin-left:auto;padding-top:0;padding-bottom:0;padding-right:0;padding-left:0;background-color:#FFFFFF;border-radius:6px;box-shadow:0 3px 6px 0 rgba(0,0,0,.05);" >
								<tr>
									<td class="ms-content-body" style="word-break:break-word;font-family:'Inter', Helvetica, Arial, sans-serif;font-size:16px;line-height:24px;padding-top:40px;padding-bottom:40px;padding-right:50px;padding-left:50px;" >
	
										<p class="logo" style="margin-right:0;margin-left:0;line-height:28px;font-weight:600;font-size:21px;color:#111111;text-align:center;margin-top:0;margin-bottom:40px;" ><span style="color:#0052e2;font-family:Arial, Helvetica, sans-serif;font-size:30px;vertical-align:bottom;" >❖&nbsp;</span>Depublic ticketing</p>
	
										<h1 style="margin-top:0;color:#111111;font-size:24px;line-height:36px;font-weight:600;margin-bottom:24px;" >Hi ` + userName + `,</h1>
	
										<p style="color:#4a5566;margin-top:20px;margin-bottom:20px;margin-right:0;margin-left:0;font-size:16px;line-height:28px;" >Thanks for registering in our system. Please confirm your email address with clicking this button.</p>
										<a href="` + confirmationUrl + `" target="_blank" style="background-color:#0052e2;padding-top:14px;padding-bottom:14px;padding-right:30px;padding-left:30px;display:inline-block;color:#FFF;text-decoration:none;border-radius:3px;-webkit-text-size-adjust:none;box-sizing:border-box;border-width:0px;border-style:solid;border-color:#0052e2;font-weight:600;font-size:15px;line-height:21px;letter-spacing:0.25px;" >Confirm my account</a>
									</td>
								</tr>
							</table>
	
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
	
	</body>
	</html>`
}

func CreateSuccessPaymentEmailHtml(userName string, status string, amont int, transactionId string) string {
	icon := "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRaJ4ujIv5bf--FOG7O-6gBYHgTbOprfyyOyg&s"
	if status == "accepted" {
		icon = "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRRqyzZL1qSjV5EIZ3upakU1NV3_SswUdtAdg&s"
	}
	// format amount to readable string
	amontFormatted := strconv.Itoa(amont)
	if amont > 1000 {
		amontFormatted = strconv.Itoa(amont/1000) + "K"
	}

	html := `<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
	<html style="-webkit-font-smoothing: antialiased; background-color: #f5f6fa; margin: 0; padding: 0;">
		<head>
			<meta name="viewport" content="width=device-width, initial-scale=1">
			<meta http-equiv="Content-Type" content="text/html; charset=UTF-8">
			<meta name="format-detection" content="telephone=no">
			<title>GO Email Templates: Generic Template</title>
			
		</head>
		<body  class="generic-template" style="-moz-osx-font-smoothing: grayscale; -webkit-font-smoothing: antialiased; background-color: #f5f6fa; margin: 0; padding: 0;">
	
			<table cellpadding="0" cellspacing="0" cols="1" align="center" style="max-width: 600px;">
				<tr bgcolor="#f5f6fa">
					<td height="50" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
				</tr>
	
				<tr bgcolor="#f5f6fa">
					<td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;">
						<!-- Seperator Start -->
						<table cellpadding="0" cellspacing="0" cols="1" bgcolor="#f5f6fa" align="center" style="max-width: 600px; width: 100%;">
							<tr bgcolor="#f5f6fa">
								<td height="30" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
							</tr>
						</table>
						<!-- Seperator End -->
	
	 <!-- Generic Pod Left Aligned with Price breakdown Start -->
						<table align="center" cellpadding="0" cellspacing="0" cols="3" bgcolor="white" class="bordered-left-right" style="border-left: 10px solid #f5f6fa; border-right: 10px solid #f5f6fa; max-width: 600px; width: 100%;">
							<tr height="50"><td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td></tr>
							<tr align="center">
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
								<td class="text-primary" style="color: #F16522; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;">
									<img src="` + icon + `" alt="GO" width="100" style="border: 0; font-size: 0; margin: 0; max-width: 100%; padding: 0;">
								</td>
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
							</tr>
							<tr height="17"><td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td></tr>
							<tr align="center">
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
								<td class="text-primary" style="color: #F16522; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;">
									<h1 style="color: #0052e2; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 30px; font-weight: 700; line-height: 34px; margin-bottom: 0; margin-top: 0;">Payment ` + status + `</h1>
								</td>
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
							</tr>
							<tr height="30"><td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td></tr>
							<tr align="left">
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
								<td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;">
									<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0;">
										Hi ` + userName + `, 
									</p>
								</td>
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
							</tr>
							<tr height="10"><td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td></tr>
							<tr align="left">
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
								<td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;">
									<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0;">Your transaction was ` + status + `!</p>
									<br>
									<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0; "><strong>Payment Details:</strong><br/>
	
	Amount: ` + amontFormatted + `<p/>
										<br>
									<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0;">We advise to keep this email for future reference.&nbsp;&nbsp;&nbsp;&nbsp;<br/></p>
								</td>
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
							</tr>
							<tr height="30">
								<td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
								<td style="border-bottom: 1px solid #D3D1D1; color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
								<td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
							</tr>
							<tr height="30"><td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td></tr>
							<tr align="center">
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
								<td style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;">
									<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0;"><strong>Transaction reference: ` + transactionId + `</strong></p>
									<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0;">Order date: ` + time.Now().Format("2006-01-02 15:04:05") + `</p>
									<p style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 22px; margin: 0;"></p>
								</td>
								<td width="36" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
							</tr>
	
							<tr height="50">
								<td colspan="3" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
							</tr>
	
						</table>
						<table cellpadding="0" cellspacing="0" cols="1" bgcolor="#f5f6fa" align="center" style="max-width: 600px; width: 100%;">
							<tr bgcolor="#f5f6fa">
								<td height="50" style="color: #464646; font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif; font-size: 14px; line-height: 16px; vertical-align: top;"></td>
							</tr>
						</table>
					</td>
				</tr>
			</table>
		</body>
	</html>`

	return html
}
