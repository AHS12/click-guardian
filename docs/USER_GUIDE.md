# Welcome to Click Guardian!

This guide will help you get started and explain what each part of the app does.

---

## Main Screen Overview

<p align="center">
  <img src="https://raw.githubusercontent.com/AHS12/click-guardian/refs/heads/main/image.png" alt="Screenshot 1" width="320" style="display:inline-block; margin-right:10px;">
  <img src="https://raw.githubusercontent.com/AHS12/click-guardian/refs/heads/main/image1.png" alt="Screenshot 2" width="320" style="display:inline-block;">
</p>

---

### 1. Blocked Double Clicks Counter

- The large circle shows how many double clicks have been blocked since you started protection.
- **When protection is active, the circle is green. When protection is stopped, the circle turns red**

---

### 2. Start/Stop Protection Button

- **Start Protection**: Click to begin blocking accidental double clicks.
- **Stop Protection**: Click to stop protection and allow all clicks as normal.

---

### 3. Settings Configuration

Click the **Settings** (gear icon) button to open the configuration window. The settings are organized into three tabs:

#### ⚙️ General Tab

- **Double-Click Threshold:**
  Adjust how much time (in milliseconds) must pass between clicks to avoid blocking.
  - Lower values = faster double clicks allowed
  - Higher values = stricter blocking

- **Application Behavior:**
  - **Minimize to system tray when closing:** If checked, closing the window will minimize the app to the tray instead of exiting.
  - **Start with Windows:** If checked, Click Guardian will launch automatically when you start your computer.

#### 🖱️ Buttons Tab

Select which mouse buttons you want to protect from accidental double-clicks. You can toggle protection individually for:

- **Left Mouse Button**
- **Right Mouse Button**
- **Middle Mouse Button**
- **Mouse Button 4 (XBUTTON1)**
- **Mouse Button 5 (XBUTTON2)**

#### 🛠️ Advanced Tab

*These settings help fix faulty mouse hardware that drops items while dragging.*

**Drag Fix (Anti-Bounce):**
- **Enable Drag Fix (Experimental):** Turns on protection against accidental drops during drags.
- **Drag Fix Threshold:** Adjusts the sensitivity for detecting a "bouncing" switch release.
  - *Problem:* A faulty mouse switch might momentarily "release" while you are dragging a file.
  - *Solution:* This feature ignores that momentary release if you re-click immediately.

**High Privilege Pause:**
- **High Privilege Pause Duration:** Sets how long protection is paused after interacting with an Admin window.
  - *Problem:* Windows security can block the app on "Administrator" programs (e.g., Task Manager), causing stuck drags.
  - *Solution:* Pauses protection when interacting with these windows.

---

### 4. Activity Log

- Shows a real-time log of all allowed and blocked clicks, with timestamps and reasons.
- **Clear Log:** Removes all entries from the log.
- **About:** Shows information about the app.

---

## Quick Start

1. Set your preferred delay using the slider.
2. Click **Start Protection**.
3. Watch the counter and log to see how many double clicks are blocked.
4. To stop, click **Stop Protection**.

---

## ⚠️ Important Anti-Cheat Warning

**MULTIPLAYER GAMING COMPATIBILITY NOTICE**

This application uses low-level mouse hooks that may be detected by anti-cheat systems as potentially malicious software. **Use with caution when playing multiplayer games**.

### Anti-Cheat Systems That May Flag This Application:

- **BattlEye** (PUBG, Rainbow Six Siege, Fortnite, etc.)
- **EasyAntiCheat** (Apex Legends, Fall Guys, Rocket League, etc.)
- **Vanguard** (Valorant) - _Extremely strict detection_
- **VAC** (Steam games like CS2, Dota 2, etc.)
- **FairFight/PunkBuster** (Battlefield series, etc.)

### Recommendations for Gamers:

1. **Stop Protection Before Gaming**: Always disable Click Guardian before launching multiplayer games
2. **Exit Completely**: Use "Quit Application" from the system tray rather than just minimizing
3. **Test in Single Player**: If unsure, test with single-player games first
4. **Create Gaming Profile**: Consider using Windows Task Scheduler to automatically stop the service during gaming hours

### Why This Happens:

Anti-cheat systems flag applications that:

- Install system-wide mouse hooks
- Block or modify mouse input events
- Monitor global system activity

While Click Guardian is legitimate accessibility software, its technical methods are similar to those used by cheating software.

**I am not responsible for any account bans or penalties resulting from anti-cheat detection. Use at your own risk with online games.**

---

## Tips

- The default delay (50–100ms) works well for most users.
- For gaming, try a lower delay.
- For accessibility, use a higher delay.

---

## System Tray

- When minimized, Click Guardian stays in your system tray.
- Right-click the tray icon for quick access to start/stop protection or exit the app.

---

## Need Help?

If you have questions or issues, please open an issue on GitHub!
