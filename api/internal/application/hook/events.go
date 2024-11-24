package hook

type Event string

const (
	EventOnDBInit Event = "on_db_init"
	EventOnStart  Event = "on_start"
	EventOnStop   Event = "on_stop"
	EventOnConfig Event = "on_config"
)
