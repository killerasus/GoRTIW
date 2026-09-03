package Utils

import (
	"image/color"
	"math"
	"math/rand"

	"github.com/engoengine/glm"
)

func Reflect(v, n *glm.Vec3) glm.Vec3 {
	mult := 2 * v.Dot(n)
	multN := n.Mul(mult)
	return v.Sub(&multN)
}

func Refract(v, n *glm.Vec3, niOverNt float32) (refracted glm.Vec3, ok bool) {
	normV := v.Normalized()
	dt := normV.Dot(n)
	discriminant := 1.0 - niOverNt*niOverNt*(1.0-dt*dt)

	if discriminant > 0 {
		sqrtDn := n.Mul(float32(math.Sqrt(float64(discriminant))))
		dtn := n.Mul(dt)
		refracted = normV.Sub(&dtn)
		refracted.MulWith(niOverNt)
		refracted.SubWith(&sqrtDn)
		ok = true
	} else {
		ok = false
	}
	return
}

func RandomInUnitSphere(r *rand.Rand) glm.Vec3 {
	onesVec := glm.Vec3{1, 1, 1}
	p := glm.Vec3{r.Float32(), r.Float32(), r.Float32()}
	p.SubWith(&onesVec)
	p.MulWith(2)

	for p.Len2() >= 1.0 {
		p = glm.Vec3{r.Float32(), r.Float32(), r.Float32()}
		p.SubWith(&onesVec)
		p.MulWith(2)
	}

	return p
}

func Schlick(cosine, refidx float32) float32 {
	r0 := (1.0 - refidx) / (1.0 + refidx)
	r0 = r0 * r0
	return r0 + (1.0-r0)*float32(math.Pow(float64(1.0-cosine), 5))
}

func ColorRGBAFromVec3(vec glm.Vec3) color.RGBA {
	return color.RGBA{R: uint8(255.99 * vec.X()), G: uint8(255.99 * vec.Y()), B: uint8(255.99 * vec.Z()), A: uint8(255)}
}
