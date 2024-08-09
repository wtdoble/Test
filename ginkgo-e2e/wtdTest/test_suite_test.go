package test_test

import (
	"fmt"
	"strings"
	"testing"

	"prometheus-collector/otelcollector/test/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	K8sClient *kubernetes.Clientset
	Cfg       *rest.Config
	//PrometheusQueryClient v1.API
)

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var _ = BeforeSuite(func() {
	var err error
	K8sClient, Cfg, err = utils.SetupKubernetesClient()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})

var _ = Describe("Test", func() {

	var cmd []string
	var podName string
	var namespace string
	var containerName string
	// var controllerLabelName string
	// var controllerLabelValue string
	// var isLinux bool
	// var apiResponse utils.APIResponse

	BeforeEach(func() {
		cmd = []string{}
		podName = "ama-metrics-bdddc945f-92gtw"
		namespace = "kube-system"
		containerName = "prometheus-collector"
		// controllerLabelName = "rsName"
		// controllerLabelValue = "ama-metrics"
		// isLinux = true
	})

	type mdsdInfoConfigLine struct {
		dt      string
		message string
	}

	type metricExtConsole struct {
		dt, status, message string
	}

	It("test", func() {
		//err := utils.QueryPromUIFromPod(K8sClient, Cfg, namespace, controllerLabelName, controllerLabelValue, containerName, "/api/v1/scrape_pools", isLinux, &apiResponse)

		//cmd = []string{"ls", "/etc/mdsd.d/config-cache/metricsextension"}
		///opt/microsoft/linuxmonagent/mdsd.warn
		///opt/microsoft/linuxmonagent/mdsd.info
		///MetricsExtensionConsoleDebugLogs.log

		//cmd = []string{"cat", "/opt/microsoft/linuxmonagent/mdsd.info"}
		cmd = []string{"cat", "/MetricsExtensionConsoleDebugLog.log"}

		stdout, _, err := utils.ExecCmd(K8sClient, Cfg, podName, containerName, namespace, cmd)
		////fmt.Println(fmt.Sprintf("stdout: %s", stdout))

		var lines []string = strings.Split(stdout, "\n")
		fmt.Println(len(lines))
		for i, line := range lines {
			fmt.Println(fmt.Sprintf("#line: %d, %s ***\r\n\r\n", i, line))

			//var l []string = strings.Split(line, " \t")
			var l []string = strings.Fields(line)
			fmt.Println(len(l))
			if len(l) >= 2 {
				//abc := mdsdInfoConfigLine{dt: l[0], message: l[1]}
				//fmt.Println(abc.dt)

				fmt.Println(fmt.Sprintf("dt: %s, status: %s", l[0], l[1]))
				fmt.Println(fmt.Sprintf("the rest: %s", strings.Join(l[2:], "%")))

			}
		}
		//fmt.Println(fmt.Sprintf("stderr: %s", stderr))
		fmt.Println(err)
		// Expect(_ = err).NotTo(HaveOccurred())
		// Expect(apiResponse.Data).NotTo(BeNil())

		// var targetsResult v1.TargetsResult
		// json.Unmarshal([]byte(apiResponse.Data), &targetsResult)
		// //fmt.Println(apiResponse)
		// fmt.Println(targetsResult)
	})

})
