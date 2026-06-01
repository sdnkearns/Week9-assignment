package main

import (
	"flag"
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	trimmedmean "github.com/sdnkearns/TrimmedMean"
	"github.com/seehuhn/mt19937"
)

const (
	targetMean = 100.0
	targetSD   = 10.0
)

func lognormParams() (muLog, sigmaLog float64) {
	sigmaLogSq := math.Log(1 + (targetSD/targetMean)*(targetSD/targetMean))
	sigmaLog = math.Sqrt(sigmaLogSq)
	muLog = math.Log(targetMean) - sigmaLogSq/2
	return
}

func main() {
	seed := flag.Int64("seed", 42, "MT19937 seed")
	n := flag.Int("n", 10000, "number of samples (minimum 100)")
	lo := flag.Float64("lo", 0.05, "trimming proportion for the low end")
	hiFlag := flag.Float64("hi", math.NaN(), "trimming proportion for the high end (default: same as -lo)")
	flag.Parse()

	hi := *hiFlag
	symmetric := math.IsNaN(hi)
	if symmetric {
		hi = *lo
	}

	rng := rand.New(mt19937.New())
	rng.Seed(*seed)

	muLog, sigmaLog := lognormParams()

	floats := make([]float64, *n)
	ints := make([]int64, *n)
	for i := range floats {
		// Box-Muller via rng.NormFloat64(), then exponentiate for log-normal.
		z := rng.NormFloat64()
		floats[i] = math.Exp(muLog + sigmaLog*z)
		ints[i] = int64(math.Round(floats[i]))
	}

	var floatMean, intMean float64
	var floatElapsed, intElapsed time.Duration
	var err error
	if symmetric {
		t0 := time.Now()
		floatMean, err = trimmedmean.TrimmedMean(floats, *lo)
		floatElapsed = time.Since(t0)
		checkErr("float trimmed mean", err)

		t0 = time.Now()
		intMean, err = trimmedmean.TrimmedMean(ints, *lo)
		intElapsed = time.Since(t0)
		checkErr("int trimmed mean", err)
	} else {
		t0 := time.Now()
		floatMean, err = trimmedmean.TrimmedMean(floats, *lo, hi)
		floatElapsed = time.Since(t0)
		checkErr("float trimmed mean", err)

		t0 = time.Now()
		intMean, err = trimmedmean.TrimmedMean(ints, *lo, hi)
		intElapsed = time.Since(t0)
		checkErr("int trimmed mean", err)
	}

	trimLabel := trimDescription(symmetric, *lo, hi)

	fmt.Println("=======================================================")
	fmt.Println(" Trimmed Mean — Go Results")
	fmt.Println("=======================================================")
	fmt.Printf(" Seed (MT19937)      : %d\n", *seed)
	fmt.Printf(" Sample size (n)     : %d\n", *n)
	fmt.Printf(" Distribution        : log-normal (mu=%.6f, sigma=%.6f)\n", muLog, sigmaLog)
	fmt.Printf(" Theoretical mean    : %.1f\n", targetMean)
	fmt.Printf(" Trimming            : %s\n", trimLabel)
	fmt.Println("-------------------------------------------------------")
	fmt.Printf(" Float trimmed mean  : %.10f\n", floatMean)
	fmt.Printf("   (dropped %d low, %d high of %d observations)\n",
		int(math.Floor(float64(*n)**lo)),
		int(math.Floor(float64(*n)*hi)),
		*n)
	fmt.Printf("   elapsed           : %s\n", floatElapsed)
	fmt.Printf(" Integer trimmed mean: %.10f\n", intMean)
	fmt.Printf("   (dropped %d low, %d high of %d observations)\n",
		int(math.Floor(float64(*n)**lo)),
		int(math.Floor(float64(*n)*hi)),
		*n)
	fmt.Printf("   elapsed           : %s\n", intElapsed)

}

func trimDescription(symmetric bool, lo, hi float64) string {
	if symmetric {
		return fmt.Sprintf("symmetric (lo = hi = %.4g)", lo)
	}
	return fmt.Sprintf("asymmetric (lo = %.4g, hi = %.4g)", lo, hi)
}

func checkErr(label string, err error) {
	if err != nil {
		fatalf("%s: %v", label, err)
	}
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", args...)
	os.Exit(1)
}
