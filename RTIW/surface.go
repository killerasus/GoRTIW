package RTIW

type Surface interface {
	Hit(ray *Ray, tMin, tMax float32, hit *HitRecord) bool
	GetMaterial() Material
}

type Surfaces struct {
	List []Surface
}

func (s *Surfaces) Add(o Surface) {
	s.List = append(s.List, o)
}

func (s *Surfaces) AddList(oo ...Surface) {
	for _, o := range oo {
		s.Add(o)
	}
}

func (s *Surfaces) Hit(ray *Ray, tMin, tMax float32, hit *HitRecord) bool {
	hitAnything := false
	closest := tMax
	for i := 0; i < len(s.List); i++ {
		if s.List[i].Hit(ray, tMin, closest, hit) {
			hitAnything = true
			hit.Material = s.List[i].GetMaterial()
			closest = hit.T
		}
	}
	return hitAnything
}
