package cats_test

import (
	"os/exec"
	"testing"

	. "github.com/cloudfoundry/cf-acceptance-tests/cats_suite_helpers"

	_ "github.com/cloudfoundry/cf-acceptance-tests/apps"
	_ "github.com/cloudfoundry/cf-acceptance-tests/backend_compatibility"
	_ "github.com/cloudfoundry/cf-acceptance-tests/detect"
	_ "github.com/cloudfoundry/cf-acceptance-tests/docker"
	_ "github.com/cloudfoundry/cf-acceptance-tests/internet_dependent"
	_ "github.com/cloudfoundry/cf-acceptance-tests/route_services"
	_ "github.com/cloudfoundry/cf-acceptance-tests/routing"
	_ "github.com/cloudfoundry/cf-acceptance-tests/security_groups"
	_ "github.com/cloudfoundry/cf-acceptance-tests/services"
	_ "github.com/cloudfoundry/cf-acceptance-tests/ssh"
	_ "github.com/cloudfoundry/cf-acceptance-tests/v3"

	"github.com/cloudfoundry-incubator/cf-test-helpers/config"
	"github.com/cloudfoundry-incubator/cf-test-helpers/helpers"
	"github.com/cloudfoundry-incubator/cf-test-helpers/workflowhelpers"
	. "github.com/cloudfoundry/cf-acceptance-tests/helpers/cli_version_check"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const minCliVersion = "6.16.1"

func TestCATS(t *testing.T) {
	RegisterFailHandler(Fail)

	Config = config.LoadConfig()

	DEFAULT_TIMEOUT = Config.DefaultTimeoutDuration()
	SLEEP_TIMEOUT = Config.SleepTimeoutDuration()
	CF_PUSH_TIMEOUT = Config.CfPushTimeoutDuration()
	LONG_CURL_TIMEOUT = Config.LongCurlTimeoutDuration()
	DETECT_TIMEOUT = Config.DetectTimeoutDuration()

	TestSetup = workflowhelpers.NewTestSuiteSetup(Config)

	var _ = BeforeSuite(func() {
		installedVersion, err := GetInstalledCliVersionString()
		Expect(err).ToNot(HaveOccurred(), "Error trying to determine CF CLI version")

		Expect(ParseRawCliVersionString(installedVersion).AtLeast(ParseRawCliVersionString(minCliVersion))).To(BeTrue(), "CLI version "+minCliVersion+" is required")
		if Config.IncludeSsh {
			ScpPath, err = exec.LookPath("scp")
			Expect(err).NotTo(HaveOccurred())

			SftpPath, err = exec.LookPath("sftp")
			Expect(err).NotTo(HaveOccurred())
		}
		TestSetup.Setup()
	})

	AfterSuite(func() {
		TestSetup.Teardown()
	})

	rs := []Reporter{}

	if Config.ArtifactsDirectory != "" {
		helpers.EnableCFTrace(Config, "CATS")
		rs = append(rs, helpers.NewJUnitReporter(Config, "CATS"))
	}

	RunSpecsWithDefaultAndCustomReporters(t, "CATS", rs)
}