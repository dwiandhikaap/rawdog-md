<div align="center">
<a href="https://github.com/dwiandhikaap/rawdog-md">
    <img src="https://github.com/dwiandhikaap/rawdog-md/blob/main/.docs/demo.gif" alt="demo">
</a>
<h1 align="center">
    rawdog-md
</h1>
<p align="center">
Caveman-level simple, Markdown-to-HTML static site tooling with minimal configuration. 
</p>
<p align="center">
Inspired by <a href="https://motherfuckingwebsite.com">motherfuckingwebsite.com</a>, uses <a href="https://pkg.go.dev/text/template">Go Template</a> and HTML under the hood.
</p>
</div>


## ✨ Features
- Easy to setup, customize, and deploy quickly.
- No bloated dependencies/framework code.
- Minification for HTML, CSS, JS, SVG, JSON, and XML
- Supports syntax highlighting **WITHOUT** JavaScript with customizable CSS.
- Supports embedding YouTube, Twitter, and Bilibili videos out of the box.
- Live reload server for local development.
- Various starter template to get going immediately


## 👍 You might want to use this if:
- You want to focus on creating your site instead of learning the tool.
- You want to build websites with Markdown and custom HTML/CSS.
- You prefer a simple setup without unnecessary installations.


## 🔑 Core Concepts
- **Minimal Config**: rawdog-md is designed to be used with minimal configuration. Choose a starter template and immediately start writing your posts in markdown.
- **Opiniated**: This tool is not meant to be a general purpose static site generator. It is meant to be used for a very specific use case, which is a static site with a few simple templated pages. If you want something more general purpose and more advanced, check out [Hugo](https://gohugo.io/) or [Jekyll](https://jekyllrb.com/).
- **Bring your own styling**: rawdog-md provides the minimal styling. However, it is very easy to adjust as it is just a plain CSS and HTML file.
- **Cross platform**: You can use it on Windows, Linux and MacOS.


## 💻 Installation
rawdog-md is available on Windows, Linux and MacOS. You can install it via package managers, manually, or build it from source.

### Windows
<details>
<summary>Install via PowerShell (Click to expand)</summary>

1. Open PowerShell as Administrator
2. Run this command
    ```shell
    Set-ExecutionPolicy RemoteSigned -Scope CurrentUser
    iex (irm https://raw.githubusercontent.com/dwiandhikaap/rawdog-md/main/.installer/powershell/install.ps1)
    ```
3. Done! You can now use it as `rawd` command in your terminal.

</details>
<details>
<summary>Install via Scoop</summary>

1. Install [Scoop](https://scoop.sh/)
2. Install rawdog-md
    ```shell
    scoop install https://raw.githubusercontent.com/dwiandhikaap/rawdog-md/main/.installer/scoop/rawdog-md.json
    ```
3. Done! You can now use it as `rawd` command in your terminal.

How to uninstall:
```shell
scoop uninstall rawdog-md
```

</details>
<details>
<summary>Install manually</summary>

1. Go to the latest [release page](https://github.com/dwiandhikaap/rawdog-md/releases/latest)
2. Download the one with `rawd-{version}-windows-amd64.zip` filename
3. Extract the zip file anywhere you like
4. Add the extracted folder to your PATH. If you don't know how, check out [this guide](https://www.architectryan.com/2018/03/17/add-to-the-path-on-windows-10/)
5. Done! You can now use it as `rawd` command in your terminal.

</details>

### Linux

<details>
<summary>Install via Bash (Click to expand)</summary>

1. Run this command in your terminal
    ```shell
    curl -fsSL https://raw.githubusercontent.com/dwiandhikaap/rawdog-md/main/.installer/bash/install.sh | bash
    ```
2. Done! You can now use it as `rawd` command in your terminal.
</details>

<details>
<summary>Install manually</summary>

1. Go to the latest [release page](https://github.com/dwiandhikaap/rawdog-md/releases/latest)
2. Copy the URL of the one with `rawd-{version}-linux-{architecture}.tar.gz` filename
3. Run this command in your terminal
    ```shell
    wget {copied-url} -O rawd.tar.gz
    tar -xvf rawd.tar.gz
    sudo mv rawd /usr/local/bin
    rm rawd.tar.gz
    ```
4. Done! You can now use it as `rawd` command in your terminal.

</details>

### MacOS
<details>
<summary>Install via Bash (Click to expand)</summary>

1. Run this command in your terminal
    ```shell
    curl -fsSL https://raw.githubusercontent.com/dwiandhikaap/rawdog-md/main/.installer/bash/install.sh | bash
    ```
2. Done! You can now use it as `rawd` command in your terminal.

</details>

<details>
<summary>Install manually</summary>
    
1. Go to the latest [release page](https://github.com/dwiandhikaap/rawdog-md/releases/latest)
2. Copy the URL of the one with `rawd-{version}-darwin-{architecture}.tar.gz` filename
3. Run this command in your terminal
    ```shell
    wget {copied-url} -O rawd.tar.gz
    tar -xvf rawd.tar.gz
    sudo mv rawd /usr/local/bin
    rm rawd.tar.gz
    ```
4. Done! You can now use it as `rawd` command in your terminal.
</details>

### Build from source

<details>

<summary>Build using Go (Click to expand)</summary>

1. Install [Go](https://golang.org/doc/install)
2. Install rawdog-md
    ```shell 
    go install github.com/dwiandhikaap/rawdog-md
    ```
    ⚠ This will install the binary as `rawdog-md` instead of `rawd`
    if you want to change it to `rawd`, you can rename the binary file in your Go bin directory. 
    
    See this [reference](https://go.dev/ref/mod#go-install) for more information about `go install`.
3. Done!

</details>


## 🚀 Usage
1. Create a new project and choose a preset template
    ```shell
    rawd init
    ```
    It will ask you for the project name and the template you want to use. 

2. Then, go to the project directory
    ```shell
    cd <your-project-name>
    ```

3. Start the development server
    ```shell
    rawd watch
    ```

4. Open your browser and go to `http://localhost:3000`
5. To create a new post, create a new file in the `pages` directory. You can use Markdown, Go Template, or HTML.
6. To edit the template, go to the `template` directory.
7. Each time you save your changes, the server will rebuild the site and refresh the browser.
8. When you're done, build the site
    ```shell
    rawd build
    ```

## 🤝 Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. 

