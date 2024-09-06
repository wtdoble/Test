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

	// fmt.Println("BeforeSuite")
	// fmt.Println(Cfg)
	// fmt.Println(err)

	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
})

func readFile(fileName string, podName string) []string {
	fmt.Printf("Examining %s\r\n", fileName)
	var cmd []string = []string{"cat", fileName}
	stdout, _, err := utils.ExecCmd(K8sClient, Cfg, podName, containerName, namespace, cmd)
	Expect(err).To(BeNil())

	return strings.Split(stdout, "\n")
}

func writeLines(lines []string) int {
	count := 0
	for _, rawLine := range lines {
		//fmt.Printf("raw line #%d: %s\r\n", i, rawLine)
		line := strings.Trim(rawLine, " ")
		if len(line) > 0 {
			//fmt.Printf("line #%d: %s\r\n", i, line)
			fmt.Printf("%s\r\n", line)
			count++
		} else {
			fmt.Println("<empty line>")
		}
	}

	return count
}

var _ = Describe("Files Test", func() {

	const mdsdErrFileName = "/opt/microsoft/linuxmonagent/mdsd.err"
	const mdsdInfoFileName = "/opt/microsoft/linuxmonagent/mdsd.info"
	const mdsdWarnFileName = "/opt/microsoft/linuxmonagent/mdsd.warn"
	const metricsExtDebugLogFileName = "/MetricsExtensionConsoleDebugLog.log"
	const ERROR = "error"
	const WARN = "warn"

	var podName string = ""
	// var apiResponse utils.APIResponse

	BeforeEach(func() {
		// cmd = []string{}
		v1Pod, err := utils.GetPodsWithLabel(K8sClient, namespace, controllerLabelName, controllerLabelValue)
		Expect(err).To(BeNil())

		fmt.Printf("pod array length: %d\r\n", len(v1Pod))
		for _, p := range v1Pod {
			fmt.Println(p.Name)
		}

		if len(v1Pod) > 0 {
			podName = v1Pod[0].Name
		}
	})

	type metricExtConsoleLine struct {
		line   string
		dt     string
		status string
		data   string
	}

	It("/opt/microsoft/linuxmonagent/mdsd.err Test", func() {

		Expect(podName).NotTo(BeEmpty())

		numErrLines := writeLines(readFile(mdsdErrFileName, podName))
		if numErrLines > 0 {
			writeLines(readFile(mdsdInfoFileName, podName))
			writeLines(readFile(mdsdWarnFileName, podName))
		}
	})

	It("/MetricsExtensionConsoleDebugLog.log Test", func() {

		Expect(podName).NotTo(BeEmpty())

		var lines []string = readFile(metricsExtDebugLogFileName, podName)

		// for i := 0; i < 10; i++ {
		// 	line := lines[i]
		for _, line := range lines {
			//fmt.Printf("#line: %d, %s \r\n", i, line)

			var fields []string = strings.Fields(line)
			if len(fields) > 2 {
				metricExt := metricExtConsoleLine{line: line, dt: fields[0], status: fields[1], data: fields[2]}
				//fmt.Println(metricExt.status)
				status := strings.ToLower(metricExt.status)
				if strings.Contains(status, ERROR) || strings.Contains(status, WARN) {
					fmt.Println(line)
				}
			}
		}
	})
})
