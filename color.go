package GoRTIW

import (
	"GoRTIW/Utils"
	"image"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/engoengine/glm"
)

func ComputeColor(ray *Ray, surfaces *Surfaces, depth int, r *rand.Rand) glm.Vec3 {
	hitRecord := HitRecord{}
	if surfaces.Hit(ray, 0.001, math.MaxFloat32, &hitRecord) {
		scattered := Ray{}
		attenuation := glm.Vec3{}
		if depth < 50 && hitRecord.Material.Scatter(ray, &hitRecord, &attenuation, &scattered, r) {
			c := ComputeColor(&scattered, surfaces, depth+1, r)
			return glm.Vec3{attenuation[0] * c[0], attenuation[1] * c[1], attenuation[2] * c[2]}
		}

		return glm.Vec3{}
	}

	unitDirection := ray.Direction.Normalized()
	t := 0.5 * (unitDirection.Y() + 1.0)
	interpA := glm.Vec3{1.0, 1.0, 1.0}
	interpB := glm.Vec3{0.5, 0.7, 1.0}
	computed := interpA.Mul(1.0 - t)
	computed.AddScaledVec(t, &interpB)
	return computed
}

func ComputePixel(i int, j int, nx int, ny int, ns int, camera *Camera, surfaces *Surfaces, output *image.RGBA, wg *sync.WaitGroup) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	acc := glm.Vec3{}
	for s := 0; s < ns; s++ {
		u := (float32(i) + r.Float32()) / float32(nx)
		v := (float32(j) + r.Float32()) / float32(ny)
		ray := camera.GetRay(u, v, r)
		color := ComputeColor(&ray, surfaces, 0, r)
		acc.AddWith(&color)
	}

	acc.MulWith(1 / float32(ns))
	acc = glm.Vec3{
		float32(math.Sqrt(float64(acc.X()))),
		float32(math.Sqrt(float64(acc.Y()))),
		float32(math.Sqrt(float64(acc.Z()))),
	}
	c := Utils.ColorRGBAFromVec3(acc)
	output.SetRGBA(i, ny-j, c)
	wg.Done()
}
