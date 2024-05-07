package cronjob

import (
	"fmt"
	"organization/repository"

	"github.com/robfig/cron"
)

func OrganizationCountCronjob(repo repository.Repository) {
	c := cron.New()

	c.AddFunc("1 23 59 * *", func() {
		fmt.Println("This function runs daily!")
		numberOfOrganizationsCreatedToday,_ := repo.FindNumberOfOrganizationsCreatedToday()
		fmt.Println(numberOfOrganizationsCreatedToday)
	})

	c.Start()
}
