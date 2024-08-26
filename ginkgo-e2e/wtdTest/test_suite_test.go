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

var _ = Describe("Test", func() {

	var cmd []string
	var podName string = ""
	// var apiResponse utils.APIResponse

	BeforeEach(func() {
		cmd = []string{}

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

	type mdsdInfoConfigLine struct {
		line   string
		dt     string
		status string
		data   string
	}

	// type metricExtConsole struct {
	// 	dt, status, message string
	// }

	type LineProcessor func(string) (bool, string)

	// func mdsdInfoLineProcessor(line string) (bool, string) {
	// 	line = strings.Trim(line, " ")
	// 	return len(line) > 0, line
	// }

	DescribeTable("Show contents of specified /opt/microsoft/linuxmonagent/ files",
		func(fileName string, proc LineProcessor) {
			Expect(podName).NotTo(BeEmpty())
			Expect(fileName).NotTo(BeEmpty())

			fullFileName := fmt.Sprintf("/opt/microsoft/linuxmonagent/%s", fileName)
			fmt.Printf("Examining %s\r\n", fullFileName)
			cmd = []string{"cat", fullFileName}
			stdout, _, err := utils.ExecCmd(K8sClient, Cfg, podName, containerName, namespace, cmd)
			Expect(err).To(BeNil())

			var lines []string = strings.Split(stdout, "\n")

			for _, rawLine := range lines {
				//fmt.Printf("raw line #%d: %s\r\n", i, rawLine)
				nonEmpty, formattedLine := proc(rawLine)
				if nonEmpty {
					//fmt.Printf("line #%d: %s\r\n", i, formattedLine)
					fmt.Printf("%s\r\n", formattedLine)
				} else {
					fmt.Println("<empty line>")
				}
			}
		},
		// func(fileName string, proc LineProcessor) string {
		// 	return fmt.Sprintf("Examining /opt/microsoft/linuxmonagent/%s", fileName)
		// },

		// Entry("Examine the contents of mdsd.info", "mdsd.info"),
		// Entry("Examine the contents of mdsd.err", "mdsd.err"),
		Entry(nil, "mdsd.info", func(line string) (bool, string) {
			line = strings.Trim(line, " ")
			return len(line) > 0, line
		}),
		Entry(nil, "mdsd.err", func(line string) (bool, string) {
			line = strings.Trim(line, " ")
			return len(line) > 0, line
		}),
	)

	It("/MetricsExtensionConsoleDebugLog Test", func() {
		//err := utils.QueryPromUIFromPod(K8sClient, Cfg, namespace, controllerLabelName, controllerLabelValue, containerName, "/api/v1/scrape_pools", isLinux, &apiResponse)

		//cmd = []string{"ls", "/etc/mdsd.d/config-cache/metricsextension"}
		// /MetricsExtensionConsoleDebugLogs.log

		Expect(podName).NotTo(BeEmpty())
		cmd = []string{"cat", "/MetricsExtensionConsoleDebugLog.log"}

		stdout, _, err := utils.ExecCmd(K8sClient, Cfg, podName, containerName, namespace, cmd)
		Expect(err).To(BeNil())
		////fmt.Println(fmt.Sprintf("stdout: %s", stdout))

		var lines []string = strings.Split(stdout, "\n")
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
	})
})
