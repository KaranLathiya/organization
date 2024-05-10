package dal

import (
	"database/sql"
	"organization/constant"
	error_handling "organization/error"
)

func StoreMicrosoftRefreshToken(db *sql.DB, refreshToken string) error {
	_, err := db.Query("UPSERT INTO public.token (id,token,event_type,token_type,updated_at) VALUES ($1, $2, $3, $4, current_timestamp())", "964917269383872513", refreshToken, constant.MICROSOFT_AUTH_EVENT_TYPE, constant.TOKEN_TYPE_REFRESH_TOKEN)
	if err != nil {
		return error_handling.InternalServerError
	}
	return nil
}

func FetchMicrosoftRefreshToken(db *sql.DB) (string, error) {
	var refreshToken string
	err := db.QueryRow("SELECT token FROM public.token WHERE event_type = $1 AND token_type = $2", constant.MICROSOFT_AUTH_EVENT_TYPE, constant.TOKEN_TYPE_REFRESH_TOKEN).Scan(&refreshToken)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return "", error_handling.NeedToLoginOnMicrosoft
		}
		return "", error_handling.InternalServerError
	}
	return refreshToken, nil
}
