
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

github.com/onsi/ginkgo/v2
go install github.com/onsi/ginkgo/v2/ginkgo
go get github.com/onsi/gomega

go mod init prometheus-collector/otelcollector/test/regionTests
ginkgo bootstrap
go get github.com/onsi/ginkgo/v2    
go get prometheus-collector/otelcollector/test/utils
//go get github.com/prometheus/prometheus/web/api/v1
go get github.com/prometheus/client_golang/api/prometheus/v1
go get k8s.io/client-go/tools/clientcmd@v0.30.3

go get github.com/Azure/azure-sdk-for-go/sdk/monitor/azquery
go get github.com/Azure/azure-sdk-for-go/sdk/azcore
go get github.com/Azure/azure-sdk-for-go/sdk/azidentity
go get github.com/Azure/azure-sdk-for-go/sdk/internal


ginkgo -- -parmRuleName=<recording rule name>


az provider list --query "[?namespace=='<resource-provider-namespace>']" --output table

az pipelines runs show --organization https://dev.azure.com/widoble0322/ --project TestPipeline  --id 388

--------------
az pipelines run --organization https://dev.azure.com/widoble0322/ --project TestPipeline --name wtdoble.Test

az pipelines build definition show --organization https://dev.azure.com/widoble0322/ --project TestPipeline --name "wtdoble.Test"

az pipelines runs list --organization https://dev.azure.com/widoble0322/ --project TestPipeline --pipeline-ids 1



{
    "id": "/subscriptions/B9842C7C-1A38-4385-8F39-A51314758BCF/resourceGroups/wtd-test/providers/microsoft.insights/actionGroups/wtdTestAg",
    "type": "Microsoft.Insights/ActionGroups",
    "name": "wtdTestAg",
    "location": "Global",
    "kind": null,
    "tags": null,
    "properties": {
        "groupShortName": "wtdTestAg",
        "enabled": true,
        "emailReceivers": [],
        "smsReceivers": [],
        "webhookReceivers": [],
        "eventHubReceivers": [],
        "itsmReceivers": [],
        "azureAppPushReceivers": [],
        "automationRunbookReceivers": [],
        "voiceReceivers": [],
        "logicAppReceivers": [],
        "azureFunctionReceivers": [],
        "armRoleReceivers": []
    }
}