package reasoner_test

import (
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/beyondcivic/goreasoner/pkg/reasoner"
)

const (
	nsKG   = "http://bench.example.com/kg/"
	nsRel  = "http://bench.example.com/rel/"
	nsRDF  = "http://www.w3.org/1999/02/22-rdf-syntax-ns#"
	nsRDFS = "http://www.w3.org/2000/01/rdf-schema#"
	nsOWL  = "http://www.w3.org/2002/07/owl#"
)

// buildBenchReasoner creates a reasoner populated with n entities and
// the standard test schema (Person/Employee/Manager + friendOf/colleagueOf/knows).
func buildBenchReasoner(nEntities int) *reasoner.Reasoner {
	estimated := nEntities*5 + 16
	r := reasoner.NewReasonerWithCapacity(estimated * 2)

	// TBox: schema axioms
	r.AddTriple(nsKG+"Employee", nsRDFS+"subClassOf", nsKG+"Person")
	r.AddTriple(nsKG+"Manager", nsRDFS+"subClassOf", nsKG+"Employee")
	r.AddTriple(nsRel+"friendOf", nsRDFS+"subPropertyOf", nsRel+"knows")
	r.AddTriple(nsRel+"colleagueOf", nsRDFS+"subPropertyOf", nsRel+"knows")
	r.AddTriple(nsRel+"friendOf", nsRDF+"type", nsOWL+"SymmetricProperty")

	// ABox: instance data
	for i := 0; i < nEntities; i++ {
		p := fmt.Sprintf("%sperson_%d", nsKG, i)
		r.AddTriple(p, nsRDF+"type", nsKG+"Person")
		r.AddTriple(p, nsKG+"name", fmt.Sprintf("\"Person_%d\"", i))

		if i%3 == 0 {
			r.AddTriple(p, nsRDF+"type", nsKG+"Employee")
		}
		if i%9 == 0 {
			r.AddTriple(p, nsRDF+"type", nsKG+"Manager")
		}

		friend := fmt.Sprintf("%sperson_%d", nsKG, (i+1)%nEntities)
		r.AddTriple(p, nsRel+"friendOf", friend)

		colleague := fmt.Sprintf("%sperson_%d", nsKG, (i+2)%nEntities)
		r.AddTriple(p, nsRel+"colleagueOf", colleague)
	}
	return r
}

// TestBenchmarkReport runs the canonical performance report across multiple scales.
// Run with: go test -run TestBenchmarkReport -v -timeout=600s
func TestBenchmarkReport(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping benchmark report in -short mode")
	}

	scales := []int{100, 1000, 10000, 100000, 1000000}

	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                  goreasoner v0.4.0 Performance Report                        ║")
	fmt.Println("╠══════════╤═══════════╤══════════════════╤══════════╤════════════╤═══════════╣")
	fmt.Println("║ Entities │ Input     │ Output           │ Time     │ Triples/s  │ Heap (MB) ║")
	fmt.Println("╠══════════╪═══════════╪══════════════════╪══════════╪════════════╪═══════════╣")

	for _, n := range scales {
		runtime.GC()
		runtime.GC()
		var m1 runtime.MemStats
		runtime.ReadMemStats(&m1)

		r := buildBenchReasoner(n)
		inputTriples := r.GetStore().Size()

		start := time.Now()
		inferred := r.RunForwardReasoning()
		elapsed := time.Since(start)

		var m2 runtime.MemStats
		runtime.ReadMemStats(&m2)
		heapMB := float64(m2.Alloc-m1.Alloc) / 1024 / 1024

		outputTriples := r.GetStore().Size()
		tps := float64(inputTriples) / elapsed.Seconds()

		fmt.Printf("║ %8d │ %9d │ %7d(+%7d) │ %8s │ %10.0f │ %7.0f MB ║\n",
			n, inputTriples, outputTriples, inferred,
			elapsed.Round(time.Millisecond), tps, heapMB)
	}

	fmt.Println("╚══════════╧═══════════╧══════════════════╧══════════╧════════════╧═══════════╝")
}

// BenchmarkReason measures pure reasoning time (excluding setup).
func BenchmarkReason(b *testing.B) {
	for _, n := range []int{1000, 10000, 100000} {
		b.Run(fmt.Sprintf("entities=%d", n), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				b.StopTimer()
				r := buildBenchReasoner(n)
				b.StartTimer()
				r.RunForwardReasoning()
			}
		})
	}
}
