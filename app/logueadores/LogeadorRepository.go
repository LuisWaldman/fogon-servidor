package logueadores

type Logeador interface {
	Login(par_1 string, par_2 string) bool
}

type LogeadorRepository struct {
	logeadores map[string]Logeador
}

func NewLogeadorRepository() *LogeadorRepository {
	return &LogeadorRepository{
		logeadores: make(map[string]Logeador),
	}
}

func (r *LogeadorRepository) Add(key string, logeador Logeador) {
	r.logeadores[key] = logeador
}

func (r *LogeadorRepository) Login(key, par_1, par_2 string) bool {
	logeador, exists := r.logeadores[key]
	if !exists {
		return false
	}
	return logeador.Login(par_1, par_2)
}
