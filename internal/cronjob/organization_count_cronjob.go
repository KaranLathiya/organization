package cronjob

import (
	"fmt"
	microsoftauth "organization/internal/microsoft-auth"
	"organization/repository"
	"strconv"

	"github.com/robfig/cron"
)

func InitializeCronjob(repo repository.Repository) {
	c := cron.New()

	c.AddFunc("1 23 59 * *", OrganizationCountCronJob(repo))

	c.Start()

}

func OrganizationCountCronJob(repo repository.Repository) func() {
	return func() {
	fmt.Println("This function runs daily!")
	numberOfOrganizationsCreatedToday, err := repo.FindNumberOfOrganizationsCreatedToday()
	if err != nil {
		return 
	}
	refreshToken, err := repo.FetchMicrosoftRefreshToken()
	if err != nil {
		return 
	}
	microsoftAuthToken, err := microsoftauth.GetAccessTokenUsingRefreshToken(refreshToken)
	if err != nil {
		return 
	}
	go repo.StoreMicrosoftRefreshToken(microsoftAuthToken.RefreshToken)
	err = microsoftauth.SendMessageOnChannel("Total number of organization created today was "+strconv.Itoa(numberOfOrganizationsCreatedToday), microsoftAuthToken.AccessToken)
	if err != nil {
		return 
	}
	fmt.Println(numberOfOrganizationsCreatedToday)
}
}
