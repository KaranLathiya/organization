package utils

import (
	"fmt"
	"strings"
)

func OrganizationCreatedMessageTemplate(ownerPhoneNumberOrEmail string, organizationName string, ownerName string) string {
	var message strings.Builder
	message.WriteString("Hello Team, We're excited to announce that a new organization has been added to our channel! Organization Name: ")
	message.WriteString(organizationName)
	message.WriteString(" Owner: ")
	message.WriteString(ownerName)
	message.WriteString(" Owner's Email/Phone Number: ")
	message.WriteString(ownerPhoneNumberOrEmail)
	message.WriteString(" Please join us in welcoming  ")
	message.WriteString(ownerName)
	message.WriteString(" and exploring the opportunities this new partnership brings to our team. If you have any questions or need further information, feel free to reach out to ")
	message.WriteString(ownerName)
	message.WriteString(" at ")
	message.WriteString(ownerPhoneNumberOrEmail)
	message.WriteString(".")
	fmt.Println(message.String())
	return message.String()
}
