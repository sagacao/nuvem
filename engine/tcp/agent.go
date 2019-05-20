package tcp

type Agent interface {
	Run()
	OnClose()
	OnConnect()
	WriteMsg(string, string, []byte) error
}
