package input

import(
    "configcenter/src/framework/core/types"
)

// Inputer is the interface that must be implemented by every Inputer.
type Inputer interface {
    // Run the input main loop. This should block until singnalled to stop by invocation of the Stop() method.
    Run(fr * types.Framework) error

    // Stop is the invoked to signal that the Run() method should its execution.
    // It will be invoked at most once.
    Stop() 
}

