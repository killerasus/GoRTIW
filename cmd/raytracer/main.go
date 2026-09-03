package main

import (
	"GoRTIW"
	"GoRTIW/Materials"
	"GoRTIW/Shapes"
	"flag"
	"image"
	"image/png"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"sync"
	"time"

	"github.com/engoengine/glm"
)

func RandomScene(r *rand.Rand) *GoRTIW.Surfaces {
	scene := GoRTIW.Surfaces{}
	scene.Add(Shapes.NewSphere(glm.Vec3{0, -1000, 0}, 1000, Materials.NewLambertian(glm.Vec3{0.5, 0.5, 0.5})))

	limit := glm.Vec3{4, 0.2, 0}

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMaterial := r.Float32()
			center := glm.Vec3{float32(a) + 0.9*r.Float32(), 0.2, float32(b) + 0.9*r.Float32()}
			dist := center.Sub(&limit)
			if dist.Len() > 0.9 { //Diffuse
				if chooseMaterial < 0.8 {
					lamb := glm.Vec3{r.Float32() * r.Float32(), r.Float32() * r.Float32(), r.Float32() * r.Float32()}
					scene.Add(Shapes.NewSphere(center, 0.2, Materials.NewLambertian(lamb)))
				} else if chooseMaterial < 0.95 { //Metal
					metal := glm.Vec3{0.5 * (1 + r.Float32()), 0.5 * (1 + r.Float32()), 0.5 * (1 + r.Float32())}
					scene.Add(Shapes.NewSphere(center, 0.2, Materials.NewMetal(metal, 0.5*r.Float32())))
				} else { //Glass
					scene.Add(Shapes.NewSphere(center, 0.2, Materials.NewDielectric(1.5)))
				}
			}
		}
	}

	scene.Add(Shapes.NewSphere(glm.Vec3{0, 1, 0}, 1.0, Materials.NewDielectric(1.5)))
	scene.Add(Shapes.NewSphere(glm.Vec3{-4, 1, 0}, 1.0, Materials.NewLambertian(glm.Vec3{0.4, 0.2, 0.1})))
	scene.Add(Shapes.NewSphere(glm.Vec3{4, 1, 0}, 1.0, Materials.NewMetal(glm.Vec3{0.7, 0.6, 0.5}, 0.0)))

	return &scene
}

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var memoryprofile = flag.String("memoryprofile", "", "write memory profile to file")

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

	nx := 1200
	ny := 800
	ns := 100

	//Camera setup
	origin := glm.Vec3{13, 2, 3}
	lookAt := glm.Vec3{0, 0, 0}
	distToFocus := float32(10.0)
	aperture := float32(0.1)

	camera := GoRTIW.NewCamera(
		origin,                  //Origin
		lookAt,                  //LookAt
		glm.Vec3{0, 1, 0},       //Up
		20,                      //FOV
		float32(nx)/float32(ny), //Aspect
		aperture,                //Aperture
		distToFocus,             //Distance to focus
	)

	output := image.NewRGBA(image.Rect(0, 0, nx, ny))

	r := rand.New(rand.NewSource(time.Now().Unix()))
	surfaces := RandomScene(r)

	var wg sync.WaitGroup
	wg.Add(nx * ny)
	for j := 0; j < ny; j++ {
		for i := 0; i < nx; i++ {
			go GoRTIW.ComputePixel(i, j, nx, ny, ns, camera, surfaces, output, &wg)
		}
	}
	wg.Wait()

	err = png.Encode(file, output)
	if err != nil {
		log.Fatal("error enconding png: ", err)
	}

	if *memoryprofile != "" {
		f, err := os.Create(*memoryprofile)
		if err != nil {
			log.Fatal(err)
		}

		pprof.WriteHeapProfile(f)
		f.Close()
	}
}
