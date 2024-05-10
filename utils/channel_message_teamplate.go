package utils

func ChannelMessageTemplate(ownerPhoneNumberOrEmail string, organizationName string, ownerName string) string {
	message := `Hello Team,\n
We're excited to announce that a new organization has been added to our channel!\n
Organization Name: ` + organizationName + `\n
Owner: ` + ownerName + `\n
Owner's Email/Phone Number: ` + ownerPhoneNumberOrEmail + `\n
Please join us in welcoming  ` + ownerName + ` and exploring the opportunities this new partnership brings to our team.\n
If you have any questions or need further information, feel free to reach out to  ` + ownerName + ` at ` + ownerPhoneNumberOrEmail + `.`
	return message
}
