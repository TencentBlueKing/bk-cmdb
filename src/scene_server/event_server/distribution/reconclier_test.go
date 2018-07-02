package distribution

import (
	"testing"
)

func TestReconciler(t *testing.T) {
	initTester()
	reconcil := newReconciler()
	reconcil.loadAll()
	reconcil.reconcile()
}
