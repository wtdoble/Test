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
	myArg                 string
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
	flag.StringVar(&myArg, "myArg", "", "Description of the usage for myArg")
}

var _ = BeforeSuite(func() {
	var err error
	K8sClient, Cfg, err = utils.SetupKubernetesClient()
	Expect(err).NotTo(HaveOccurred())

	//amwQueryEndpoint := "https://wtdaks9-amw-gjexfkctfvb6c5gr.westus2.prometheus.monitor.azure.com"
	amwQueryEndpoint := os.Getenv("AMW_QUERY_ENDPOINT")
	fmt.Printf("env: %s\r\n", amwQueryEndpoint)
	Expect(amwQueryEndpoint).NotTo(BeEmpty())

	PrometheusQueryClient, err = utils.CreatePrometheusAPIClient(amwQueryEndpoint)
	Expect(err).NotTo(HaveOccurred())
	Expect(PrometheusQueryClient).NotTo(BeNil())

	fmt.Printf("myArg: %s", myArg)

	fmt.Println("CHECKING ALERTS")
	//var a v1.AlertsResult
	warnings, result, err := utils.InstantQuery(PrometheusQueryClient, "alerts")
	//a, err = utils.InstantQuery(PrometheusQueryClient, "alerts") //Alerts(context.Background())
	//fmt.Println(a)
	fmt.Println(warnings)
	fmt.Println(result)
	fmt.Println(err)
	Expect(err).NotTo(HaveOccurred())

	fmt.Println("CHECKING RULES")
	//var a v1.AlertsResult
	warnings2, result2, err := utils.InstantQuery(PrometheusQueryClient, "node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate")
	//a, err = utils.InstantQuery(PrometheusQueryClient, "alerts") //Alerts(context.Background())
	//fmt.Println(a)
	fmt.Println(warnings2)
	fmt.Println(result2)
	fmt.Println(err)
	// var r v1.RulesResult
	// r, err = PrometheusQueryClient.Rules(context.Background())
	fmt.Println(err)
	Expect(err).NotTo(HaveOccurred())
	//fmt.Println(r)
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

	It("metrics", func() {
		var query string
		query = "up"

		warnings, result, err := utils.InstantQuery(PrometheusQueryClient, query)
		Expect(err).NotTo(HaveOccurred())
		Expect(warnings).To(BeEmpty())

		fmt.Println(result)
	})
})
