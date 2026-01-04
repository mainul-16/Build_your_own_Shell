[![progress-banner](https://backend.codecrafters.io/progress/shell/8eb383d3-818a-4a47-8947-f72b4766def1)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is a starting point for Go solutions to the
["Build Your Own Shell" Challenge](https://app.codecrafters.io/courses/shell/overview).

In this challenge, you'll build your own POSIX compliant shell that's capable of
interpreting shell commands, running external programs and builtin commands like
cd, pwd, echo and more. Along the way, you'll learn about shell command parsing,
REPLs, builtin commands, and more.

**Note**: If you're viewing this repo on GitHub, head over to
[codecrafters.io](https://codecrafters.io) to try the challenge.

# Passing the first stage

The entry point for your `shell` implementation is in `app/main.go`. Study and
uncomment the relevant code, and push your changes to pass the first stage:

```sh
git commit -am "pass 1st stage" # any msg
git push origin master
```

Time to move on to the next stage!

# Stage 2 & beyond

Note: This section is for stages 2 and beyond.
1. Ensure you have `go (1.25)` installed locally
1. Run `./your_program.sh` to run your program, which is implemented in
   `app/main.go`.
1. Commit your changes and run `git push origin master` to submit your solution
   to CodeCrafters. Test output will be streamed to your terminal.

<hr>
<h1 align="center">Build Your Own Shell (Go)</h1>

<p align="center">
  A custom Unix-like shell implemented in Go with support for built-in commands,
  pipelines, redirection, history, and autocompletion.
</p>

<hr/>

<h2>🚀 Features</h2>
<ul>
  <li>Interactive shell with custom prompt</li>
  <li>Built-in commands: <code>exit</code>, <code>cd</code>, <code>pwd</code>, <code>echo</code>, <code>type</code>, <code>history</code></li>
  <li>Pipeline support (<code>|</code>)</li>
  <li>Input/output redirection (<code>&gt;</code>, <code>&gt;&gt;</code>, <code>2&gt;</code>)</li>
  <li>Command history with file persistence</li>
  <li>Tab autocompletion</li>
</ul>

<hr/>

<h2>📦 Requirements</h2>
<ul>
  <li>Go <strong>1.25+</strong></li>
  <li>Linux / macOS / WSL (Windows Subsystem for Linux)</li>
</ul>

<hr/>

<h2>⚙️ How to Run (Local Execution)</h2>

<ol>
  <li>
    <strong>Clone the repository</strong>
    <pre><code>git clone https://github.com/mainul-16/Build_your_own_Shell.git
cd Build_your_own_Shell</code></pre>
  </li>

  <li>
    <strong>Navigate to the source directory</strong>
    <pre><code>cd app</code></pre>
  </li>

  <li>
    <strong>Build the shell</strong>
    <pre><code>go build -o myshell</code></pre>
  </li>

  <li>
    <strong>Launch the shell</strong>
    <pre><code>./myshell</code></pre>
  </li>
</ol>

<p>You should now see an interactive shell prompt:</p>
<pre><code>$ pwd
$ echo hello
hello</code></pre>

<hr/>

<h2>🐳 Docker (Local Testing)</h2>

<p>
This project was also tested using Docker for local, containerized execution
to ensure environment consistency and reproducible builds.
</p>

<ul>
  <li>A multi-stage Docker build was used locally</li>
  <li>The Dockerfile is intentionally not committed to this repository</li>
  <li>Docker was used strictly for development and testing purposes</li>
</ul>

<p>
This approach keeps the repository clean while still validating Docker-based
deployment.
</p>

<hr/>

<h2>📁 Project Structure</h2>
<pre><code>Build_your_own_Shell/
├── app/
│   └── main.go
├── go.mod
├── go.sum
├── README.md
└── your_program.sh
</code></pre>

<hr/>

<h2>🧠 Notes !!</h2>
<ul>
  <li>This project focuses on understanding shell internals and process handling</li>
  <li>It is designed for learning purposes and system-level programming practice</li>
</ul>

<hr/>

<h2>📜 License</h2>
<p>
This project is intended for educational use.
</p>
