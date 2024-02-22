package rpcserver

type RPCServer interface {
	StartServer(port string, rcvr ...any) error
}
