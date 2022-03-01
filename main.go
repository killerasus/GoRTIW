package main

import (
	"RTIW/RTIW"
	"RTIW/RTIW/Materials"
	"RTIW/RTIW/Shapes"
	"RTIW/RTIW/Utils"
	"flag"
	"image"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/engoengine/glm"
)

func ComputeColor(ray *RTIW.Ray, surfaces *RTIW.Surfaces, depth int, r *rand.Rand) glm.Vec3 {
	hitRecord := RTIW.HitRecord{}
	if surfaces.Hit(ray, 0.001, math.MaxFloat32, &hitRecord) {
		scattered := RTIW.Ray{}
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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	file, err := os.Create("output.png")
	if err != nil {
		log.Fatal("error creating ouput file: ", err)
	}

	defer file.Close()

	nx := 200
	ny := 100
	ns := 100

	//Camera setup
	origin := glm.Vec3{3, 3, 2}
	lookAt := glm.Vec3{0, 0, -1}
	focusVector := origin.Sub(&lookAt)
	distToFocus := focusVector.Len()
	aperture := float32(2.0)

	camera := RTIW.NewCamera(
		origin,                  //Origin
		lookAt,                  //LookAt
		glm.Vec3{0, 1, 0},       //Up
		20,                      //FOV
		float32(nx)/float32(ny), //Aspect
		aperture,                //Aperture
		distToFocus,             //Distance to focus
	)

	output := image.NewRGBA(image.Rect(0, 0, nx, ny))

	surfaces := RTIW.Surfaces{
		List: []RTIW.Surface{
			Shapes.NewSphere(glm.Vec3{0, 0, -1}, 0.5, Materials.NewLambertian(glm.Vec3{0.1, 0.2, 0.5})),
			Shapes.NewSphere(glm.Vec3{0, -100.5, -1}, 100, Materials.NewLambertian(glm.Vec3{0.8, 0.8, 0.0})),
			Shapes.NewSphere(glm.Vec3{1, 0, -1}, 0.5, Materials.NewMetal(glm.Vec3{0.8, 0.6, 0.2}, 1.0)),
			Shapes.NewSphere(glm.Vec3{-1, 0, -1}, 0.5, Materials.NewDieletric(1.5)),
		},
	}

	r := rand.New(rand.NewSource(time.Now().Unix()))

	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			acc := glm.Vec3{}
			for s := 0; s < ns; s++ {
				u := (float32(i) + r.Float32()) / float32(nx)
				v := (float32(j) + r.Float32()) / float32(ny)
				ray := camera.GetRay(u, v, r)
				color := ComputeColor(&ray, &surfaces, 0, r)
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
		}
	}

	err = png.Encode(file, output)
	if err != nil {
		log.Fatal("error enconding jpg: ", err)
	}
}
