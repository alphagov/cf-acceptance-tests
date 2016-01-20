package backend_compatibility

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"github.com/cloudfoundry-incubator/cf-test-helpers/cf"
	"github.com/cloudfoundry-incubator/cf-test-helpers/generator"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/app_helpers"
	"github.com/cloudfoundry/cf-acceptance-tests/helpers/assets"
)

const binaryHi = "Hello from a binary"

var _ = Describe("Backend Compatibility", func() {
	var appName string

	BeforeEach(func() {
		appName = generator.RandomName()
		Eventually(cf.Cf(
			"push", appName,
			"-p", assets.NewAssets().Binary,
			"--no-start",
			"-m", DEFAULT_MEMORY_LIMIT,
			"-b", "binary_buildpack",
			"-c", "./app"),
			CF_PUSH_TIMEOUT).Should(Exit(0))
	})

	AfterEach(func() {
		app_helpers.AppReport(appName, DEFAULT_TIMEOUT)
		Eventually(cf.Cf("delete", appName, "-f"), DEFAULT_TIMEOUT).Should(Exit(0))
	})

	Describe("An app staged on Diego", func() {
		BeforeEach(func() {
			app_helpers.EnableDiego(appName)

			Eventually(cf.Cf("start", appName), CF_PUSH_TIMEOUT).Should(Exit(0))
			Eventually(helpers.CurlingAppRoot(appName), DEFAULT_TIMEOUT).Should(ContainSubstring(binaryHi))
		})

		It("runs on the DEAs", func() {
			app_helpers.DisableDiego(appName)
			Eventually(cf.Cf("restart", appName), CF_PUSH_TIMEOUT).Should(Exit(0))
			Eventually(helpers.CurlingAppRoot(appName), DEFAULT_TIMEOUT).Should(ContainSubstring(binaryHi))
		})
	})

	Describe("An app staged on the DEA", func() {
		BeforeEach(func() {
			app_helpers.DisableDiego(appName)
			Eventually(cf.Cf("start", appName), CF_PUSH_TIMEOUT).Should(Exit(0))
			Eventually(helpers.CurlingAppRoot(appName), DEFAULT_TIMEOUT).Should(ContainSubstring(binaryHi))
		})

		It("runs on Diego", func() {
			app_helpers.EnableDiego(appName)
			Eventually(cf.Cf("restart", appName), CF_PUSH_TIMEOUT).Should(Exit(0))
			Eventually(helpers.CurlingAppRoot(appName), DEFAULT_TIMEOUT).Should(ContainSubstring(binaryHi))
		})
	})
})
