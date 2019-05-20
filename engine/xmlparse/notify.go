package xmlparse

type Notifyer interface {
	Callback(*XMLConfig)
}
