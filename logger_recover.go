package golog

type RecoverParams struct {
	ExitCode int
	Panic    bool
}

func (l *loggerImpl) Recover(params ...RecoverParams) {
	p := RecoverParams{
		Panic: false,
	}
	if len(params) > 0 {
		p = params[0]
	}
	if v := recover(); v != nil {
		var msg Message
		if p.Panic {
			if p.ExitCode == 0 {
				p.ExitCode = 1
			}
			msg = l.Fatal(p.ExitCode)
		} else {
			msg = l.Error()
		}
		if err, ok := v.(error); ok {
			msg.Throw(err)
		} else {
			msg.Send("%v", v)
		}
	}
}
