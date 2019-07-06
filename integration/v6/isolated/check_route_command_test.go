package isolated

import (
	"code.cloudfoundry.org/cli/integration/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gbytes"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("check-route command", func() {
	var (
		orgName   string
		spaceName string
		route     helpers.Route
	)
	BeforeEach(func() {
		helpers.SkipIfClientCredentialsTestMode()
		orgName = helpers.NewOrgName()
		spaceName = helpers.NewSpaceName()

		helpers.SetupCF(orgName, spaceName)
		route = helpers.NewRoute(spaceName, helpers.DefaultSharedDomain(), "integration", "")
	})

	AfterEach(func() {
		helpers.QuickDeleteOrg(orgName)
	})

	It("checks routes", func() {
		session := helpers.CF("check-route", route.Host, route.Domain)
		Eventually(session).Should(Say("Route %s.%s does not exist", route.Host, route.Domain))
		Eventually(session).Should(Exit(0))

		route.Create()

		session = helpers.CF("check-route", route.Host, route.Domain)
		Eventually(session).Should(Say("Route %s.%s does exist", route.Host, route.Domain))
		Eventually(session).Should(Exit(0))
	})
})
