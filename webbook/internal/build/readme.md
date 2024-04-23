在internal目录下面执行
rm webook
go mod tidy
GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
docker rmi -f learn_go/webook:v0.0.1
docker build -t learn_go/webook:v0.0.1 .