package error

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/gookit/validate"
	"github.com/lib/pq"
)

func init() {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})
	// validate.AddGlobalMessages(map[string]string{
	// 	"minLength": "OO! {field} min length is %d",
	// 	"required": "oh! the {field} is required",
	// 	"email": "email is invalid",
	//   })
}

// func (f UserForm) Messages() map[string]string {
// 	return validate.MS{
// 		"required": "oh! the {field} is required",
// 		"email": "email is invalid",
// 		"Name.required": "message for special field",
// 		"Age.int": "age must int",
// 		"Age.min": "age min value is 1",
// 	}
// }

type CustomError struct {
	StatusCode   int           `json:"statusCode" validate:"required" `
	ErrorMessage string        `json:"errorMessage" validate:"required" `
	InvalidData  []InvalidData `json:"invalidData" validate:"omitempty" `
}

type InvalidData struct {
	Field string
	Error map[string]string
}

func (c CustomError) Error() string {
	return c.ErrorMessage
}

func ErrorMessageResponse(w http.ResponseWriter, err error) {
	if error, ok := err.(CustomError); ok {
		response, _ := json.MarshalIndent(error, "", "  ")
		w.WriteHeader(error.StatusCode)
		w.Write(response)
		return
	}
	w.WriteHeader(http.StatusInternalServerError)
	response, _ := json.MarshalIndent(InternalServerError, "", "  ")
	w.Write(response)
}

func CreateCustomError(errorMessage string, statusCode int, invalidData ...InvalidData) error {
	return CustomError{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
		InvalidData:  invalidData,
	}
}

var (
	MarshalError         = CreateCustomError("Error while marshling data.", http.StatusInternalServerError)
	UnmarshalError       = CreateCustomError("Error while unmarshling data.", http.StatusBadRequest)
	InternalServerError  = CreateCustomError("Internal Server Error.", http.StatusInternalServerError)
	OTPGenerateError     = CreateCustomError("Error at generating OTP.", http.StatusInternalServerError)
	BcryptError          = CreateCustomError("Error at bcypting.", http.StatusInternalServerError)
	SendEmailError       = CreateCustomError("Error at sending email.", http.StatusInternalServerError)
	SendMessageError     = CreateCustomError("Error at sending message.", http.StatusInternalServerError)
	ExpiredOTP           = CreateCustomError("OTP expired.", http.StatusGone)
	InvalidOTP           = CreateCustomError("OTP is invalid.", http.StatusBadRequest)
	UserAlreadyExist     = CreateCustomError("User already exist.", http.StatusBadRequest)
	UserDoesNotExist     = CreateCustomError("User does not exist.", http.StatusNotFound)
	HeaderDataMissing    = CreateCustomError("Required header not found.", http.StatusBadRequest)
	UsernameAlreadyTaken = CreateCustomError("Username already taken.", http.StatusBadRequest)
	InvalidDetails       = CreateCustomError("Invalid details provided.", http.StatusBadRequest)
	AlreadyBlocked       = CreateCustomError("Already blocked.", http.StatusBadRequest)

	AlreadyInvited             = CreateCustomError("Already invited to the organization.", http.StatusConflict)
	OrganizationDoesNotExist   = CreateCustomError("Organization does not exist.", http.StatusNotFound)
	NoAccessRights             = CreateCustomError("You don't have permission due to the role.", http.StatusForbidden)
	AlreadyMember              = CreateCustomError("Already member of the organization.", http.StatusConflict)
	OwnerRoleChangeRestriction = CreateCustomError("Can't change the organization's owner role.", http.StatusForbidden)
	OwnerLeaveRestriction      = CreateCustomError("Can't leave the organization without assigning owner role to other.", http.StatusForbidden)
	OwnRoleChangeRestriction   = CreateCustomError("Can't change own role.", http.StatusForbidden)
	OwnRemoveRestriction       = CreateCustomError("Can't remove yourself from organization.", http.StatusForbidden)
	OwnerRemoveRestriction     = CreateCustomError("Can't remove the owner.", http.StatusForbidden)
	OwnerAccessRights          = CreateCustomError("Only owner have permission.", http.StatusForbidden)
	NotMemberOfOrganization    = CreateCustomError("You aren't member of organization.", http.StatusNotFound)
	JWTErrSignatureInvalid     = CreateCustomError("Invalid signature on jwt token.", http.StatusUnauthorized)
	JWTTokenInvalid            = CreateCustomError("Invalid jwt token.", http.StatusBadRequest)
	JWTTokenInvalidDetails     = CreateCustomError("Invalid jwt token details.", http.StatusBadRequest)
	NeedToLoginOnMicrosoft     = CreateCustomError("Need to login on microsoft.", http.StatusUnauthorized)
	// InvitationAlreadyExist   = CreateCustomError("Already invited the member.", http.StatusOK)

	NotNullConstraintError    = CreateCustomError("Required field cannot be empty or null. Please provide a value for the field.", http.StatusBadRequest)
	ForeignKeyConstraintError = CreateCustomError("Data doesn't exist.", http.StatusConflict)
	UniqueKeyConstraintError  = CreateCustomError("Data already exists.", http.StatusConflict)
	CheckConstraintError      = CreateCustomError("Data doesn't meet the required criteria.", http.StatusBadRequest)
	NoRowsError               = CreateCustomError("Data doesn't exist.", http.StatusNotFound)
)

func LogErrorMessage(err error) {
	pc, file, line, _ := runtime.Caller(1)
	functionName := runtime.FuncForPC(pc).Name()
	log.Printf("Error %s in file %s, function %s, line %d", err.Error(), file, functionName, line)
}

func DatabaseErrorShow(err error) error {
	if err == sql.ErrNoRows {
		return NoRowsError
	}
	if dbErr, ok := err.(*pq.Error); ok {
		errCode := dbErr.Code
		switch errCode {
		case "23502":
			// not-null constraint violation
			return NotNullConstraintError

		case "23503":
			// foreign key violation
			return ForeignKeyConstraintError

		case "23505":
			// unique constraint violation
			return UniqueKeyConstraintError

		case "23514":
			// check constraint violation
			return CheckConstraintError
		}
	}
	return InternalServerError
}
