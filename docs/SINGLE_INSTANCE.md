# Single Instance Implementation

This document explains how Click Guardian prevents multiple instances from running simultaneously.

## Overview

Click Guardian implements a single instance mechanism to prevent users from accidentally launching multiple copies of the application. When a user attempts to launch a second instance, they will see a dialog informing them that the application is already running and directing them to the system tray icon.

## Implementation Details

### Package Structure

The single instance functionality is implemented using the `github.com/juju/mutex/v2` package.

### Mechanism

The implementation uses a named mutex provided by the juju/mutex package:

1. **Named Mutex**
   - Creates a named mutex using the juju/mutex package
   - Named as `click-guardian-single-instance`
   - If the mutex already exists, another instance is running
   - The package handles cross-platform mutex creation

### Code Flow

1. **Application Startup** (`cmd/click-guardian/main.go`)
   - Defines a mutex specification with a unique name
   - Attempts to acquire the mutex using `mutex.Acquire()`
   - If acquisition fails, shows the "Already Running" dialog and exits
   - If successful, continues with normal application startup

2. **Mutex Acquisition**
   - Uses `github.com/juju/mutex/v2` package to create and acquire the mutex
   - If another instance exists, `mutex.Acquire()` returns an error
   - If this is the first instance, `mutex.Acquire()` returns a releaser

3. **Mutex Release**
   - Automatically handled when the application exits
   - The returned `releaser` is deferred to ensure mutex release

### User Experience

When a user attempts to launch a second instance:

1. The second process starts but immediately detects the running instance
2. A dialog appears with the message:
   ```
   Click Guardian is already running.
   
   Look for the icon in your system tray or check your taskbar.
   ```
3. The second process exits after the user acknowledges the dialog

## Benefits

- **Prevents Confusion**: Users won't accidentally run multiple copies
- **Resource Efficiency**: Only one mouse hook is installed
- **Consistent State**: Configuration and settings remain consistent
- **User Guidance**: Clear instructions on finding the running instance

## Edge Cases Handled

- **Crash Recovery**: If the application crashes, the mutex is automatically released
- **Force Termination**: If the process is killed, the mutex is released

## Platform Considerations

The implementation uses the juju/mutex package which provides cross-platform mutex support.

## Testing

To test the single instance functionality:

1. Launch Click Guardian normally
2. Attempt to launch a second instance
3. Verify that the "Already Running" dialog appears
4. Close the first instance
5. Launch the application again to verify normal startup
