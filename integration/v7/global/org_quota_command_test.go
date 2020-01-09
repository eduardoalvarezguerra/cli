package global

import (
	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("org-quota command", func() {
	var (
		quotaName string
	)

	BeforeEach(func() {
		quotaName = helpers.QuotaName()
	})

	Describe("help", func() {
		When("--help flag is set", func() {
			It("Displays command usage to output", func() {
				session := helpers.CF("org-quota", "--help")
				Eventually(session).Should(Say("NAME:"))
				Eventually(session).Should(Say("quota - Show quota info"))
				Eventually(session).Should(Say("USAGE:"))
				Eventually(session).Should(Say("cf quota QUOTA"))
				Eventually(session).Should(Say("SEE ALSO:"))
				Eventually(session).Should(Say("org, quotas"))
				Eventually(session).Should(Exit(0))
			})
		})
	})

	When("the environment is not setup correctly", func() {
		It("fails with the appropriate errors", func() {
			helpers.CheckEnvironmentTargetedCorrectly(false, false, ReadOnlyOrg, "org", "org-name")
		})
	})

	When("the environment is set up correctly", func() {
		BeforeEach(func() {
			helpers.LoginCF()
		})

		When("the quota does not exist", func() {
			It("displays quota not found and exits 1", func() {
				session := helpers.CF("org-quota", quotaName)
				userName, _ := helpers.GetCredentials()
				Eventually(session).Should(Say(`Getting org quota %s as %s\.\.\.`, quotaName, userName))
				Eventually(session.Err).Should(Say("FAILED"))
				Eventually(session).Should(Say("Quota %s not found", quotaName))
				Eventually(session).Should(Exit(1))
			})
		})

		When("the quota exists", func() {
			When("no flags are used", func() {
				It("displays a table with quota names and their values and exits 0", func() {
					session := helpers.CF("org-quota", "default")
					userName, _ := helpers.GetCredentials()
					Eventually(session).Should(Say(`Getting org quota %s as %s\.\.\.`, "default", userName))
					Eventually(session).Should(Say("OK"))

					Eventually(session).Should(Say(`Total Memory\s+100G`))
					Eventually(session).Should(Say(`Instance Memory\s+unlimited`))
					Eventually(session).Should(Say(`Routes\s+1000`))
					Eventually(session).Should(Say(`Services\s+unlimited`))
					Eventually(session).Should(Say(`Paid service plans\s+allowed`))
					Eventually(session).Should(Say(`App instance limit\s+unlimited`))
					Eventually(session).Should(Say(`Reserved Route Ports\s+100`))

					Eventually(session).Should(Exit(0))
				})
			})
		})
	})
})