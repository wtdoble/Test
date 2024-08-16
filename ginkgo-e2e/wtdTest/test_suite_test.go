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

const namespace = "kube-system"
const containerName = "prometheus-collector"
const controllerLabelName = "rsName"
const controllerLabelValue = "ama-metrics"

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

var _ = BeforeSuite(func() {
	var err error
	K8sClient, Cfg, err = utils.SetupKubernetesClient()

	fmt.Println("BeforeSuite")
	fmt.Println(Cfg)
	fmt.Println(err)

	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})

var _ = Describe("Test", func() {

	var cmd []string
	var podName string = ""
	// var isLinux bool
	// var apiResponse utils.APIResponse

	BeforeEach(func() {
		cmd = []string{}
		////podName = "ama-metrics-57c4f5c898-twwn7"
		// isLinux = true

		v1Pod, err := utils.GetPodsWithLabel(K8sClient, namespace, controllerLabelName, controllerLabelValue)
		Expect(err).To(BeNil())
		//fmt.Println(err)

		fmt.Printf("pod array length: %d\r\n", len(v1Pod))
		for _, p := range v1Pod {
			fmt.Println(p.Name)
		}

		if len(v1Pod) > 0 {
			podName = v1Pod[0].Name
		}
	})

	type mdsdInfoConfigLine struct {
		line   string
		dt     string
		status string
		data   string
	}

	type metricExtConsole struct {
		dt, status, message string
	}

	It("/opt/microsoft/linuxmonagent/mdsd.info Test", func() {
		// /opt/microsoft/linuxmonagent/mdsd.warn
		// /opt/microsoft/linuxmonagent/mdsd.info

		Expect(podName).NotTo(BeEmpty())
		cmd = []string{"cat", "/opt/microsoft/linuxmonagent/mdsd.info"}

		stdout, _, err := utils.ExecCmd(K8sClient, Cfg, podName, containerName, namespace, cmd)
		Expect(err).To(BeNil())

		var lines []string = strings.Split(stdout, "\n")
		fmt.Println(len(lines))

		for i, line := range lines {
			fmt.Printf("line #%d: %s\r\n", i, line)
		}
	})

	It("/MetricsExtensionConsoleDebugLog Test", func() {
		//err := utils.QueryPromUIFromPod(K8sClient, Cfg, namespace, controllerLabelName, controllerLabelValue, containerName, "/api/v1/scrape_pools", isLinux, &apiResponse)

		//cmd = []string{"ls", "/etc/mdsd.d/config-cache/metricsextension"}
		// /MetricsExtensionConsoleDebugLogs.log

		Expect(podName).NotTo(BeEmpty())
		cmd = []string{"cat", "/MetricsExtensionConsoleDebugLog.log"}

		stdout, stderr, err := utils.ExecCmd(K8sClient, Cfg, podName, containerName, namespace, cmd)
		Expect(err).To(BeNil())
		////fmt.Println(fmt.Sprintf("stdout: %s", stdout))

		fmt.Printf("stderr: %s", stderr)
		//fmt.Println(err)

		var lines []string = strings.Split(stdout, "\n")
		fmt.Println(len(lines))
		//for line = lines[0, 10] {
		for i := 0; i < 10; i++ {
			line := lines[i]
			fmt.Printf("#line: %d, %s ***\r\n", i, line)

			//var l []string = strings.Split(line, " \t")
			var l []string = strings.Fields(line)
			fmt.Println(len(l))
			if len(l) >= 2 {
				abc := mdsdInfoConfigLine{line: line, dt: l[0], status: l[1], data: l[2]}
				fmt.Println(abc.status)

				// fmt.Println(fmt.Sprintf("dt: %s, status: %s", l[0], l[1]))
				// fmt.Println(fmt.Sprintf("the rest: %s", strings.Join(l[2:], "%")))
			}
		}

		fmt.Println("Error test")
		fmt.Println("Err test")
		fmt.Println("Warning test")
		fmt.Println("Warn test")

		// Expect(_ = err).NotTo(HaveOccurred())
		// Expect(apiResponse.Data).NotTo(BeNil())

		// var targetsResult v1.TargetsResult
		// json.Unmarshal([]byte(apiResponse.Data), &targetsResult)
		// //fmt.Println(apiResponse)
		// fmt.Println(targetsResult)
	})
})
