package RTIW

type Surface interface {
	Hit(ray *Ray, tMin, tMax float32, hit *HitRecord) bool
}

type Surfaces struct {
	List []Surface
}

func (s *Surfaces) Hit(ray *Ray, tMin, tMax float32, hit *HitRecord) bool {
	hitAnything := false
	closest := tMax
	for i := 0; i < len(s.List); i++ {
		if s.List[i].Hit(ray, tMin, closest, hit) {
			hitAnything = true
			closest = hit.T
		}
	}
	return hitAnything
}
