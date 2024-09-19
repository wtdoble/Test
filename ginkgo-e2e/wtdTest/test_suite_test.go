package test_test

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"testing"

	"prometheus-collector/otelcollector/test/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

var (
	K8sClient             *kubernetes.Clientset
	Cfg                   *rest.Config
	PrometheusQueryClient v1.API
	parmRuleName          string
)

const namespace = "kube-system"
const containerName = "prometheus-collector"
const controllerLabelName = "rsName"
const controllerLabelValue = "ama-metrics"

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}

func init() {
	flag.StringVar(&parmRuleName, "parmRuleName", "", "Prometheus rule name to use in this test suite")
}

var _ = BeforeSuite(func() {
	var err error
	K8sClient, Cfg, err = utils.SetupKubernetesClient()
	Expect(err).NotTo(HaveOccurred())

	amwQueryEndpoint := os.Getenv("AMW_QUERY_ENDPOINT")
	fmt.Printf("env (AMW_QUERY_ENDPOINT): %s\r\n", amwQueryEndpoint)
	Expect(amwQueryEndpoint).NotTo(BeEmpty())

	PrometheusQueryClient, err = utils.CreatePrometheusAPIClient(amwQueryEndpoint)
	Expect(err).NotTo(HaveOccurred())
	Expect(PrometheusQueryClient).NotTo(BeNil())

	fmt.Printf("parmRuleName: %s", parmRuleName)
	Expect(parmRuleName).ToNot(BeEmpty())
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

var _ = Describe("Regions Suite", func() {

	const mdsdErrFileName = "/opt/microsoft/linuxmonagent/mdsd.err"
	const mdsdInfoFileName = "/opt/microsoft/linuxmonagent/mdsd.info"
	const mdsdWarnFileName = "/opt/microsoft/linuxmonagent/mdsd.warn"
	const metricsExtDebugLogFileName = "/MetricsExtensionConsoleDebugLog.log"
	const ERROR = "error"
	const WARN = "warn"

	var podName string = ""

	type metricExtConsoleLine struct {
		line   string
		dt     string
		status string
		data   string
	}

	BeforeEach(func() {
		v1Pod, err := utils.GetPodsWithLabel(K8sClient, namespace, controllerLabelName, controllerLabelValue)
		Expect(err).To(BeNil())
		Expect(len(v1Pod)).To(BeNumerically(">", 0))

		fmt.Printf("pod array length: %d\r\n", len(v1Pod))
		fmt.Printf("Available pods matching '%s'='%s'\r\n", controllerLabelName, controllerLabelValue)
		for _, p := range v1Pod {
			fmt.Println(p.Name)
		}

		if len(v1Pod) > 0 {
			podName = v1Pod[0].Name
			fmt.Printf("Choosing the pod: %s\r\n", podName)
		}

		Expect(podName).ToNot(BeEmpty())
	})

	Context("Examine selected files and directories", func() {

		It("Check that there are no errors in /opt/microsoft/linuxmonagent/mdsd.err", func() {

			numErrLines := writeLines(readFile(mdsdErrFileName, podName))
			if numErrLines > 0 {
				writeLines(readFile(mdsdInfoFileName, podName))
				writeLines(readFile(mdsdWarnFileName, podName))
			}
		})

		It("Enumerate all the 'error' or 'warning' records in /MetricsExtensionConsoleDebugLog.log", func() {

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

		It("Check that /etc/mdsd.d/config-cache/metricsextension exists", func() {

			var cmd []string = []string{"ls", "/etc/mdsd.d/config-cache/"}
			stdout, _, err := utils.ExecCmd(K8sClient, Cfg, podName, containerName, namespace, cmd)
			Expect(err).To(BeNil())

			metricsExtExists := false

			list := strings.Split(stdout, "\n")
			for i := 0; i < len(list) && !metricsExtExists; i++ {
				s := list[i]
				fmt.Println(s)
				metricsExtExists = (strings.Compare(s, "metricsextension") == 0)
			}

			Expect(metricsExtExists).To(BeTrue())
		})
	})

	Context("Examine Prometheus via the AMW", func() {
		It("Query for a metric", func() {
			query := "up"

			fmt.Printf("Examining metrics via the query: '%s'", query)

			warnings, result, err := utils.InstantQuery(PrometheusQueryClient, query)
			Expect(err).NotTo(HaveOccurred())
			Expect(warnings).To(BeEmpty())

			fmt.Println(result)
		})

		It("Check that the specified recording rule exists", func() {
			fmt.Printf("Examining the recording rule: %s", parmRuleName)

			warnings, result, err := utils.InstantQuery(PrometheusQueryClient, parmRuleName)

			fmt.Println(warnings)
			Expect(err).NotTo(HaveOccurred())

			fmt.Println(result)
		})

		It("Query Prometheus alerts", func() {
			warnings, result, err := utils.InstantQuery(PrometheusQueryClient, "alerts")

			fmt.Println(warnings)
			Expect(err).NotTo(HaveOccurred())

			fmt.Println(result)
		})

		It("", func() {

		})

	})
})
