
Grace's "E2E Prometheus Region Tests" document
https://microsoft-my.sharepoint.com/:w:/r/personal/grwehner_microsoft_com/_layouts/15/Doc.aspx?sourcedoc=%7B477a7ef2-0c27-4f97-a4bb-252b2b104864%7D&action=edit&wdPid=57637e9e&share=IQHyfnpHJwyXT6S7JSsrEEhkAUxvS9IoEDdLcxY73TZFU6s

Go reference
https://pkg.go.dev/

Ginkgo reference and tutorial
https://onsi.github.io/ginkgo/

az aks get-credentials --resource-group wtd-test --name wtdaks9

\users\widoble\.kube\config
..\..\kubectl exec ama-metrics-bdddc945f-92gtw -it -c prometheus-collector -n kube-system -- bash
..\..\kubectl config get-contexts

helm install bitnami/prometheus --generate-name

go mod init prometheus-collector/otelcollector/test/regionTests
ginkgo bootstrap
go get github.com/onsi/ginkgo/v2    
go get prometheus-collector/otelcollector/test/utils
