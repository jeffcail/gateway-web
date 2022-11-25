package dao

type ServiceDetail struct {
	Info *GatewayServiceInfo `json:"info" description:"基本信息"`
	HTTPRule *GatewayServiceHttpRule `json:"http_rule" description:"http_rule"`
	TcpRule *TcpRule `json:"tcp_rule" description:"tcp_rule"`
	GrpcRule *GrpcRule `json:"grpc_rule" description:"grpc_rule"`
	LoadBalance *GatewayServiceLoadBalance `json:"load_balance" description:"load_balance"`
	AccessControl *GatewayServiceAccessControl `json:"access_control" description:"access_control"`
}