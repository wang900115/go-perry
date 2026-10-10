package coroutine

import (
	"fmt"
	"testing"
)

func TestRawCoroutineSwitch(t *testing.T) {
	// Native pointer to record the main thread's execution progress snapshot
	var mainCoroRaw *coro

	// 1. Define the native coroutine task
	task := func(self *Coro, parent *coro) {
		fmt.Println("[Coroutine] 1. Task started! Collaborating seamlessly on the exact same OS Thread.")

		// Capture the authentic native context descriptor (bookmark) passed from the main thread
		mainCoroRaw = parent

		fmt.Println("[Coroutine] 2. Capture successful! Actively returning CPU control manually to the main thread...")
		SwitchRaw(mainCoroRaw) // Pass the ball back to the main thread

		fmt.Println("[Coroutine] 4. Main thread manually woke me up again! Executing phase 2 of the task...")
		fmt.Println("[Coroutine] 5. Manually yielding control back to the main thread once more...")

		// Upon each switch, the runtime automatically updates the current coroutine context state.
		// We directly reuse the captured mainCoroRaw descriptor to navigate back.
		SwitchRaw(mainCoroRaw)

		fmt.Println("[Coroutine] 7. Coroutine is about to terminate, performing final exit handover.")
		SwitchRaw(mainCoroRaw)
	}

	fmt.Println("[Main Thread] A. Initializing a clean native coroutine using the newcoro linkname black magic...")
	taskCoro := Create(task)

	fmt.Println("[Main Thread] B. Invoking SwitchTo() for the first time to spin up the coroutine...")
	taskCoro.SwitchTo() // Pass the context to the coroutine -> runs up to step 2 and yields back

	fmt.Println("[Main Thread] C. Main thread regained active CPU control.")
	fmt.Println("[Main Thread] D. Bypassing standard scheduling constraints; directly wake up the coroutine via SwitchTo...")
	taskCoro.SwitchTo() // Pass the context to the coroutine -> runs up to step 5 and yields back

	fmt.Println("[Main Thread] E. Main thread regained active CPU control again.")
	fmt.Println("[Main Thread] F. Final manual context switch to let the coroutine run to its safe completion...")
	taskCoro.SwitchTo() // Pass the context to the coroutine -> runs to the end of the func closure

	fmt.Println("[Main Thread] G. Experiment complete! 0 locks, 0 crashes, pure native runtime context switches succeeded!")
}
