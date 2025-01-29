package configs

type EtcdGrpcCfg struct {
	EndPoints   []string
	ServiceName string
	PoolSize    int
}

func (ec *EtcdGrpcCfg) SetEndPoints(endPoints []string) *EtcdGrpcCfg {
	ec.EndPoints = endPoints
	return ec
}

func (ec *EtcdGrpcCfg) SetServiceName(serviceName string) *EtcdGrpcCfg {
	ec.ServiceName = serviceName
	return ec
}

func (ec *EtcdGrpcCfg) SetPoolSize(poolSize int) *EtcdGrpcCfg {
	ec.PoolSize = poolSize
	return ec
}
