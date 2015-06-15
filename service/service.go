package service

type(
	Service struct {
		Name string `json:"name"`
		NbInstances int `json:"nb_instances"`
		MaxInstances int `json:"max_instances"`
		Manager func()
	}
	Services []Service
)