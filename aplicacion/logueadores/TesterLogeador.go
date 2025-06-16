package logueadores

type TesterLogeador struct {
	claves []string
}

func NewTesterLogeador(claves []string) *TesterLogeador {
	return &TesterLogeador{claves: claves}
}

func (l *TesterLogeador) Login(par_1 string, par_2 string) bool {
	for _, usuario := range l.claves {
		if par_2 == usuario {
			return true
		}
	}
	return false
}
