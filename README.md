# catree

**catree** is a command-line tool that recursively prints the names and
contents of all text files in a directory and its subdirectories. It
automatically detects text files by analyzing their content and skips binary or
archive files. Files and directories ignored by `.gitignore` are also skipped.

---

## Features

- Recursively prints all text files in a directory tree
- Skips non-text files
- Respects `.gitignore` rules
- Clearly separates each file’s name and contents in the output

---

## Example

Given a directory structure like:

project/
├── .gitignore
├── README.md
├── main.go
├── config.yaml
├── docs/
│   └── usage.txt
├── scripts/
│   └── deploy.sh
├── images/
│   └── logo.png
└── archive/
└── data.zip


Running `catree project` will output:


```
==.gitignore==
{.gitignore contents}

==README.md==
{README.md contents}

==main.go==
{main.go contents}

==config.yaml==
{config.yaml contents}

==docs/usage.txt==
{usage.txt contents}

==scripts/deploy.sh==
{deploy.sh contents}
```


Files such as `images/logo.png` and `archive/data.zip` are not ignored as
unprintable. Files and directories listed in `.gitignore` are also skipped.

---

## Installation

You can install `catree` using Go or by building from source.

### Option 1: Install with `go install` (Recommended)

If you have Go 1.17+:

```sh
go install github.com/jmeyering/catree@latest
```

This will place the `catree` binary in your `$GOPATH/bin` or `$HOME/go/bin` directory.  
Make sure this directory is in your `PATH`:

```sh
export PATH="$PATH:$(go env GOPATH)/bin"
```

### Option 2: Build from Source

1. **Clone the repository:**

   ```sh
   git clone https://github.com/jmeyering/catree.git
   cd catree
   ```

2. **Install dependencies:**

   ```sh
   go get github.com/sabhiram/go-gitignore
   ```

3. **Build the binary:**

   ```sh
   go build -o catree
   ```

4. **(Optional) Move to your `bin` directory:**

   ```sh
   mv catree ~/bin/
   # or for system-wide install
   sudo mv catree /usr/local/bin/
   ```

---

Now you can run `catree` from anywhere in your terminal.

