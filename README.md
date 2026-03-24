# 🎙️ EchoMind

EchoMind is a colorful and animated CLI voice recorder built in Go. It provides a visually engaging "mini-app" experience right inside your terminal, featuring real-time waveforms, smooth animations, and interactive configuration.

## ✨ Features

- **🚀 Animated Startup**: A lively initialization sequence with progress bars.
- **🎙️ Real-time Recording**: Capture audio with a live blinking indicator and timer.
- **📊 Dynamic Waveform**: Visual feedback based on actual microphone input levels.
- **⚙️ Interactive Config**: Easily set your default recording format and save directory.
- **🌈 Styled UI**: Beautiful terminal interface powered by Bubble Tea and Lipgloss.
- **💾 Auto-Naming**: Files are saved with precise timestamps for easy organization.

---

## 📥 Installation (For Users)

### Prerequisites
- **Windows Users**: You must have a C compiler installed (like [Mingw-w64](https://www.mingw-w64.org/) or via [MSYS2](https://www.msys2.org/)) because the audio engine uses CGO.

### Steps
1. **Clone the repository**:
   ```bash
   git clone https://github.com/your-repo/echomind.git
   cd echomind
   ```

2. **Build the executable**:
   ```bash
   go build -o echomind.exe main.go
   ```

3. **Run EchoMind**:
   ```bash
   ./echomind.exe record
   ```

4. **(Optional) Add to Path**: Move `echomind.exe` to a folder in your system's PATH to run it from anywhere.

---

## 🛠️ Publication & Development (For Devs)

### Project Architecture
- `cmd/`: Command definitions using **Cobra**.
- `internal/audio/`: Audio capture logic using **malgo** (miniaudio wrapper) and **go-audio**.
- `internal/ui/`: Terminal UI components using **Bubble Tea** and **Lipgloss**.
- `internal/config/`: JSON-based preference management.

### Development Setup
1. **Install Dependencies**:
   ```bash
   go mod download
   ```

2. **Running in Debug Mode**:
   Since Bubble Tea takes over the terminal, use `fmt.Fprintf(os.Stderr, ...)` or log to a file for debugging.

3. **CGO Requirements**:
   This project requires `CGO_ENABLED=1`. On Windows, ensure `gcc` is in your environment variables.

### Publication / Release Workflow
To prepare a release for multiple platforms:

1. **Windows (Primary Focus)**:
   ```bash
   go build -ldflags="-s -w" -o dist/echomind-windows-amd64.exe main.go
   ```

2. **Cross-Compilation**:
   Note that because this uses `malgo` (CGO), cross-compiling requires a cross-compiler for the target OS (e.g., `x86_64-linux-gnu-gcc` for Linux).

3. **Versioning**:
   Tag your releases in Git:
   ```bash
   git tag -a v1.0.0 -m "Initial release"
   git push origin v1.0.0
   ```

### Future Enhancements
- [x] Add MP3/FLAC encoding support (requires additional C libraries).
- [x] Implement a `history` command to browse past recordings.
- [ ] Add audio playback within the CLI.

---

## 🎮 Usage
- `echomind record` (or `em record` after install): Start a recording session.
- `echomind config`: Change your default save path, format, and quality.
- `echomind history`: View a log of your past recordings.
- `echomind install`: Automatically add EchoMind to your Windows PATH and create an `em` shorthand.
- `echomind --help`: See all available options.

### Keyboard Shortcuts (Recording)
- **Enter**: Stop recording / Confirm save.
- **Arrows**: Navigate fields in the save menu.
- **'o'**: Open the saved recording in your default app.
- **'r'**: Start a new recording session immediately after saving.
- **'q'**: Quit the application.

Enjoy your digital voice journaling! 🚀
