package constant

const (
	OTP_SENT                     = "Successfully OTP sent."
	USER_DETAILS_UPDATED         = "User details successfully updated."
	ALREADY_INVITED              = "Already invited to the organization."
	INVITED                      = "Successfully invited to the organization."
	ORGANIZATION_DETAILS_UPDATED = "Organization details successfully updated."
	INVITATION_REJECTED          = "Invitation rejected successfully."
	INVITATION_ACCEPTED          = "Invitation accepted successfully and joined the organization."
	ROLE_UPDATED                 = "Role updated successfully."
	LEAVE_ORGANIZATION           = "Successfully left the organization."
	REMOVED_FROM_ORGANIZATION    = "Successfully removed from organization."
	OWNERSHIP_TRANSFERRED        = "Ownership of organization transferred successfully."
	ORGANIZATION_DELETED         = "Organization deleted successfully."

	ORGANIZATION_ROLE_OWNER  = "owner"
	ORGANIZATION_ROLE_ADMIN  = "admin"
	ORGANIZATION_ROLE_VIEWER = "viewer"
	ORGANIZATION_ROLE_EDITOR = "editor"

	INVITATION_ACCEPT = "accept"
	INVITATION_REJECT = "reject"

	USER_SERVICE_BASE_URL = "http://localhost:8000/"

	MICROSOFT_AUTH_URL           = "https://login.microsoftonline.com/common/oauth2/v2.0/token"
	MICROSOFT_GRAPH_API_BASE_URL = "https://graph.microsoft.com"

	MICROSOFT_AUTH_EVENT_TYPE = "microsoft_auth"

	TOKEN_TYPE_REFRESH_TOKEN = "refresh_token"
	TOKEN_TYPE_ACCESS_TOKEN  = "access_token"
)
