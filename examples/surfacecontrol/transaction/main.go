// Transaction lifecycle and chaining pattern example.
//
// Demonstrates:
//   - Creating and closing a Transaction
//   - Applying an empty transaction (a valid no-op in SurfaceFlinger)
//   - The chaining pattern: every setter returns *Transaction so calls compose
//
// A Transaction batches property changes (position, visibility, z-order, etc.)
// to one or more SurfaceControls and submits them atomically via Apply().
// Creating a SurfaceControl requires an ANativeWindow (obtained from an
// Activity), so this example focuses on the Transaction itself, which can be
// created standalone.
//
// This program must run on an Android device with API level 29+.
package main

import (
	"fmt"
	"log"

	"github.com/xaionaro-go/ndk/surfacecontrol"
)

func main() {
	// --- Create a transaction ---
	// NewTransaction allocates an ASurfaceTransaction on the native side.
	txn := surfacecontrol.NewTransaction()
	defer func() {
		if err := txn.Close(); err != nil {
			log.Printf("close transaction: %v", err)
		}
	}()
	fmt.Println("transaction created")

	// --- Apply an empty transaction ---
	// An empty transaction is valid. SurfaceFlinger accepts it as a no-op,
	// which is useful for flushing pending state or as a synchronization point.
	txn.Apply()
	fmt.Println("empty transaction applied")

	// --- Chaining pattern ---
	// Every setter (SetVisibility, SetPosition, SetScale, SetZOrder,
	// SetColor, SetBufferAlpha, SetBufferTransparency, SetDamageRegion)
	// returns *Transaction, so multiple property changes compose into a
	// single fluent expression:
	//
	//   txn.
	//       SetVisibility(sc, surfacecontrol.Show).
	//       SetPosition(sc, 100, 200).
	//       SetScale(sc, 2.0, 2.0).
	//       SetZOrder(sc, 10).
	//       SetBufferAlpha(sc, 0.8).
	//       SetBufferTransparency(sc, surfacecontrol.Translucent).
	//       SetColor(sc, 0.2, 0.4, 0.8, 1.0, 0).
	//       Apply()
	//
	// All setters require a *SurfaceControl target. Since SurfaceControl
	// creation needs an ANativeWindow (from an Activity), the chained call
	// above is shown in a comment. In a real app you would obtain the
	// SurfaceControl from your window and chain setters exactly like this.

	// --- Visibility and Transparency enums ---
	// The package exports two enums used with SetVisibility and
	// SetBufferTransparency:
	//
	//   Visibility:   Show (0), Hide (1)
	//   Transparency: Transparent (0), Translucent (1), Opaque (2)
	fmt.Printf("Visibility   constants: Show=%d, Hide=%d\n",
		surfacecontrol.Show, surfacecontrol.Hide)
	fmt.Printf("Transparency constants: Transparent=%d, Translucent=%d, Opaque=%d\n",
		surfacecontrol.Transparent, surfacecontrol.Translucent, surfacecontrol.Opaque)

	// --- Multiple transactions ---
	// You can create several transactions, configure them independently,
	// and apply them in any order. Each Apply() is an atomic
	// SurfaceFlinger commit.
	txn2 := surfacecontrol.NewTransaction()
	txn2.Apply()
	if err := txn2.Close(); err != nil {
		log.Fatalf("close txn2: %v", err)
	}
	fmt.Println("second transaction applied and closed")

	fmt.Println("done")
}
