package compute

const (
	CommandSET        = "SET"
	CommandGET        = "GET"
	CommandDEL        = "DEL"
	CommandSETArgsNum = 2
	CommandGETArgsNum = 1
	CommandDELArgsNum = 1
)

func getCommandArgsNum(cmd string) int {
	switch cmd {
	case CommandSET:
		return CommandSETArgsNum
	case CommandGET:
		return CommandGETArgsNum
	case CommandDEL:
		return CommandDELArgsNum
	default:
		return -1
	}
}
