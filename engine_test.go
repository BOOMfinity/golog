package golog

func printSimple(log Logger) {
	log.Info().Send("test")
}

func printParams2(log Logger) {
	msg := log.Info()
	msg.Param("one", "czxc")
	msg.Param("two", "fsdf")
	msg.Param("three", "fsdf")
	msg.Send("test")
}

func printParams10(log Logger) {
	msg := log.Info()
	msg.Param("one", 1)
	msg.Param("two", "fsdf")
	msg.Param("three", "2514")
	msg.Param("four", .43)
	msg.Param("five", "sdfsdf")
	msg.Param("six", "sdasli3,mn")
	msg.Param("seven", "cvsdf")
	msg.Param("eight", 36251)
	msg.Param("nine", "sdfsdf")
	msg.Param("ten", 34123)
	msg.Send("test")
}

func printFmtStringArg(log Logger) {
	log.Info().Send("hello %s", "world")
}

func printFmtIntArg(log Logger) {
	log.Info().Send("hello %d", 1234)
}
