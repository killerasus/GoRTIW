package GoRTIW

import (
	"math"
	"math/rand"
	"testing"

	"github.com/engoengine/glm"
)

// ---------------------------------------------------------------------------
// NewCamera tests
// ---------------------------------------------------------------------------

// TestNewCameraReturnsNonNil verifies that NewCamera always returns a non-nil Camera.
func TestNewCameraReturnsNonNil(t *testing.T) {
	cam := NewCamera(
		glm.Vec3{0, 0, 0},
		glm.Vec3{0, 0, -1},
		glm.Vec3{0, 1, 0},
		90, 1.0, 0, 1,
	)
	if cam == nil {
		t.Fatal("Expected non-nil Camera, got nil")
	}
}

// TestNewCameraOriginMatchesLookFrom verifies the camera origin equals lookFrom.
func TestNewCameraOriginMatchesLookFrom(t *testing.T) {
	lookFrom := glm.Vec3{1, 2, 3}
	cam := NewCamera(lookFrom, glm.Vec3{0, 0, -1}, glm.Vec3{0, 1, 0}, 60, 1.5, 0, 1)

	if cam.Origin != lookFrom {
		t.Errorf("Expected Origin %v, got %v", lookFrom, cam.Origin)
	}
}

// TestNewCameraLensRadius verifies the lens radius is aperture/2.
func TestNewCameraLensRadius(t *testing.T) {
	aperture := float32(1.4)
	cam := NewCamera(
		glm.Vec3{0, 0, 0}, glm.Vec3{0, 0, -1}, glm.Vec3{0, 1, 0},
		90, 1.0, aperture, 1,
	)
	expected := aperture / 2.0
	if cam.LensRadius != expected {
		t.Errorf("Expected LensRadius %v, got %v", expected, cam.LensRadius)
	}
}

// TestNewCameraZeroAperture verifies a pinhole camera (aperture == 0) has zero lens radius.
func TestNewCameraZeroAperture(t *testing.T) {
	cam := NewCamera(
		glm.Vec3{0, 0, 0}, glm.Vec3{0, 0, -1}, glm.Vec3{0, 1, 0},
		90, 1.0, 0, 1,
	)
	if cam.LensRadius != 0 {
		t.Errorf("Expected LensRadius 0 for zero aperture, got %v", cam.LensRadius)
	}
}

// TestNewCameraOrthonormalBasis verifies U, V, W form a right-handed orthonormal basis.
func TestNewCameraOrthonormalBasis(t *testing.T) {
	cam := NewCamera(
		glm.Vec3{0, 0, 0}, glm.Vec3{0, 0, -1}, glm.Vec3{0, 1, 0},
		90, 1.0, 0, 1,
	)

	const eps = 1e-5

	// |U| == 1
	lenU := float32(math.Sqrt(float64(cam.U[0]*cam.U[0] + cam.U[1]*cam.U[1] + cam.U[2]*cam.U[2])))
	if absF32(lenU-1) > eps {
		t.Errorf("Expected |U| == 1, got %v", lenU)
	}

	// |V| == 1
	lenV := float32(math.Sqrt(float64(cam.V[0]*cam.V[0] + cam.V[1]*cam.V[1] + cam.V[2]*cam.V[2])))
	if absF32(lenV-1) > eps {
		t.Errorf("Expected |V| == 1, got %v", lenV)
	}

	// |W| == 1
	lenW := float32(math.Sqrt(float64(cam.W[0]*cam.W[0] + cam.W[1]*cam.W[1] + cam.W[2]*cam.W[2])))
	if absF32(lenW-1) > eps {
		t.Errorf("Expected |W| == 1, got %v", lenW)
	}

	// U ⊥ W  (U · W == 0)
	dotUW := cam.U[0]*cam.W[0] + cam.U[1]*cam.W[1] + cam.U[2]*cam.W[2]
	if absF32(dotUW) > eps {
		t.Errorf("Expected U · W == 0 (orthogonal), got %v", dotUW)
	}

	// V ⊥ W  (V · W == 0)
	dotVW := cam.V[0]*cam.W[0] + cam.V[1]*cam.W[1] + cam.V[2]*cam.W[2]
	if absF32(dotVW) > eps {
		t.Errorf("Expected V · W == 0 (orthogonal), got %v", dotVW)
	}
}

// TestNewCameraWPointsFromLookAtToLookFrom verifies W is in the direction from
// lookAt → lookFrom (i.e., opposite to the viewing direction), normalised.
func TestNewCameraWPointsFromLookAtToLookFrom(t *testing.T) {
	lookFrom := glm.Vec3{0, 0, 5}
	lookAt := glm.Vec3{0, 0, 0}
	cam := NewCamera(lookFrom, lookAt, glm.Vec3{0, 1, 0}, 60, 1.0, 0, 1)

	// Expected W == (0, 0, 1) — pointing away from the scene
	expected := glm.Vec3{0, 0, 1}
	if !vec3Near(cam.W, expected, 1e-5) {
		t.Errorf("Expected W %v, got %v", expected, cam.W)
	}
}

// TestNewCameraHorizontalVerticalScale verifies that a wider field-of-view
// produces a larger Horizontal extent than a narrower one at the same aspect.
func TestNewCameraHorizontalVerticalScale(t *testing.T) {
	wide := NewCamera(
		glm.Vec3{0, 0, 0}, glm.Vec3{0, 0, -1}, glm.Vec3{0, 1, 0},
		120, 1.0, 0, 1,
	)
	narrow := NewCamera(
		glm.Vec3{0, 0, 0}, glm.Vec3{0, 0, -1}, glm.Vec3{0, 1, 0},
		30, 1.0, 0, 1,
	)

	wideH := vec3Len(wide.Horizontal)
	narrowH := vec3Len(narrow.Horizontal)

	if wideH <= narrowH {
		t.Errorf("Expected wider FoV to have larger Horizontal (%v), but got narrow=%v wide=%v", narrowH, narrowH, wideH)
	}
}

// ---------------------------------------------------------------------------
// GetRay tests
// ---------------------------------------------------------------------------

// TestGetRayReturnsSameOriginForZeroAperture verifies that with aperture == 0
// the ray origin always equals the camera origin (pinhole model).
func TestGetRayReturnsSameOriginForZeroAperture(t *testing.T) {
	lookFrom := glm.Vec3{1, 2, 3}
	cam := NewCamera(lookFrom, glm.Vec3{0, 0, 0}, glm.Vec3{0, 1, 0}, 90, 1.0, 0, 1)
	r := rand.New(rand.NewSource(42))

	for i := 0; i < 10; i++ {
		ray := cam.GetRay(r.Float32(), r.Float32(), r)
		if !vec3Near(ray.Origin, lookFrom, 1e-5) {
			t.Errorf("Sample %d: Expected origin %v, got %v", i, lookFrom, ray.Origin)
		}
	}
}

// TestGetRayCenterRay verifies the direction of the centre ray (s=0.5, t=0.5)
// for a pinhole camera looking straight down -Z.
// The direction should point roughly toward -Z (W points +Z, so centre
// direction has a negative W component).
func TestGetRayCenterRay(t *testing.T) {
	cam := NewCamera(
		glm.Vec3{0, 0, 0}, glm.Vec3{0, 0, -1}, glm.Vec3{0, 1, 0},
		90, 1.0, 0, 1,
	)
	r := rand.New(rand.NewSource(0))
	ray := cam.GetRay(0.5, 0.5, r)

	// With aperture 0 the offset is zero; the direction should have a negative Z component.
	if ray.Direction[2] >= 0 {
		t.Errorf("Expected centre ray to have negative Z direction, got %v", ray.Direction)
	}
}

// TestGetRayNonZeroApertureVariesOrigin verifies that with a non-zero aperture
// the ray origin is jittered (i.e., not identical across calls).
func TestGetRayNonZeroApertureVariesOrigin(t *testing.T) {
	cam := NewCamera(
		glm.Vec3{0, 0, 5}, glm.Vec3{0, 0, 0}, glm.Vec3{0, 1, 0},
		90, 1.0, 2.0, 5,
	)
	r := rand.New(rand.NewSource(1234))

	first := cam.GetRay(0.5, 0.5, r)
	allSame := true
	for i := 0; i < 20; i++ {
		ray := cam.GetRay(0.5, 0.5, r)
		if ray.Origin != first.Origin {
			allSame = false
			break
		}
	}
	if allSame {
		t.Error("Expected ray origins to vary with non-zero aperture, but all were the same")
	}
}

// ---------------------------------------------------------------------------
// helpers (camera-local, avoids collision with color_test helpers)
// ---------------------------------------------------------------------------

func absF32(v float32) float32 {
	if v < 0 {
		return -v
	}
	return v
}

func vec3Len(v glm.Vec3) float32 {
	return float32(math.Sqrt(float64(v[0]*v[0] + v[1]*v[1] + v[2]*v[2])))
}

func vec3Near(a, b glm.Vec3, eps float32) bool {
	return absF32(a[0]-b[0]) < eps &&
		absF32(a[1]-b[1]) < eps &&
		absF32(a[2]-b[2]) < eps
}
