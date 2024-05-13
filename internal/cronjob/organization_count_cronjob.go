package cronjob

import (
	"fmt"
	error_handling "organization/error"
	microsoftauth "organization/internal/microsoft-auth"
	"organization/repository"
	"strconv"

	"github.com/robfig/cron"
)

func InitializeCronjob(repo repository.Repository) {
	c := cron.New()

	c.AddFunc("00 59 23 * * *", OrganizationCountCronJob(repo))

	c.Start()

}

func OrganizationCountCronJob(repo repository.Repository) func() {
	return func() {
		fmt.Println("This function runs daily!")
		numberOfOrganizationsCreatedToday, err := repo.FindNumberOfOrganizationsCreatedToday()
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
		refreshToken, err := repo.FetchMicrosoftRefreshToken()
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
		microsoftAuthToken, err := microsoftauth.GetAccessTokenUsingRefreshToken(refreshToken)
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
		go repo.StoreMicrosoftRefreshToken(microsoftAuthToken.RefreshToken)
		err = microsoftauth.SendMessageOnChannel("Total number of organization created today was "+strconv.Itoa(numberOfOrganizationsCreatedToday), microsoftAuthToken.AccessToken)
		if err != nil {
			error_handling.LogErrorMessage(err)
			return
		}
		fmt.Println(numberOfOrganizationsCreatedToday)
	}
}
