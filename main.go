package main

import (
	"RTIW/RTIW"
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/engoengine/glm"
)

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

func ColorRGBAFromVec3(vec glm.Vec3) color.RGBA {
	return color.RGBA{R: uint8(255.99 * vec.X()), G: uint8(255.99 * vec.Y()), B: uint8(255.99 * vec.Z()), A: uint8(255)}
}

func ComputeColor(ray *RTIW.Ray, surfaces *RTIW.Surfaces, r *rand.Rand) glm.Vec3 {
	hitRecord := RTIW.HitRecord{}
	if surfaces.Hit(ray, 0.001, math.MaxFloat32, &hitRecord) {
		random := RandomInUnitSphere(r)
		target := hitRecord.P
		target.AddWith(&hitRecord.Normal)
		target.AddWith(&random)
		target.SubWith(&hitRecord.P)
		newRay := RTIW.Ray{Origin: hitRecord.P, Direction: target}
		n := ComputeColor(&newRay, surfaces, r)
		n = n.Mul(0.5)
		return n
	}
	unitDirection := ray.Direction.Normalized()
	t := 0.5 * (unitDirection.Y() + 1.0)
	interpA := glm.Vec3{1.0, 1.0, 1.0}
	interpB := glm.Vec3{0.5, 0.7, 1.0}
	computed := interpA.Mul(1.0 - t)
	computed.AddScaledVec(t, &interpB)
	return computed
}

func main() {

	file, err := os.Create("output.jpg")
	if err != nil {
		log.Fatal("error creating ouput file: ", err)
	}

	defer file.Close()

	camera := RTIW.Camera{
		Origin:          glm.Vec3{0.0, 0.0, 0.0},
		LowerLeftCorner: glm.Vec3{-2.0, -1.0, -1.0},
		Horizontal:      glm.Vec3{4.0, 0.0, 0.0},
		Vertical:        glm.Vec3{0.0, 2.0, 0.0},
	}

	nx := 200
	ny := 100
	ns := 100

	output := image.NewRGBA(image.Rect(0, 0, nx, ny))

	var surfaces RTIW.Surfaces
	sphere1 := RTIW.Sphere{Center: glm.Vec3{0.0, 0.0, -1}, Radius: 0.5}
	sphere2 := RTIW.Sphere{Center: glm.Vec3{0.0, -100.5, -1}, Radius: 100}
	surfaces.List = append(surfaces.List, &sphere1)
	surfaces.List = append(surfaces.List, &sphere2)

	r := rand.New(rand.NewSource(time.Now().Unix()))

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			acc := glm.Vec3{}
			for s := 0; s < ns; s++ {
				u := (float32(i) + r.Float32()) / float32(nx)
				v := (float32(j) + r.Float32()) / float32(ny)
				ray := camera.GetRay(u, v)
				color := ComputeColor(&ray, &surfaces, r)
				acc.AddWith(&color)
			}

			acc.MulWith(1 / float32(ns))
			acc = glm.Vec3{
				float32(math.Sqrt(float64(acc.X()))),
				float32(math.Sqrt(float64(acc.Y()))),
				float32(math.Sqrt(float64(acc.Z()))),
			}
			c := ColorRGBAFromVec3(acc)
			output.SetRGBA(i, ny-j, c)
		}
	}

	err = jpeg.Encode(file, output, nil)
	if err != nil {
		log.Fatal("error enconding jpg: ", err)
	}
}
