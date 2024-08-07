<div align="center">
<a href="https://github.com/dwiandhikaap/rawdog-md">
    <img src="https://github.com/dwiandhikaap/rawdog-md/blob/main/demo.gif" alt="demo">
</a>
<h1 align="center">
    rawdog-md
</h1>
<p align="center">
Caveman-level simple, Markdown-to-HTML static blog tooling with minimal configuration. 
</p>
<p align="center">
Inspired by <a href="https://motherfuckingwebsite.com">motherfuckingwebsite.com</a>, uses <a href="https://pkg.go.dev/text/template">Go Template</a> and HTML under the hood.
</p>
</div>



## 👍 You might want to use this if:
- You want to write your blog posts in markdown.
- You want to get started with writing your blog posts ASAP without caring much about config files.
- You don't wanna install things like Node.js, Ruby, or Python.
- You want to cook your own HTML and CSS and then forget about it.
- You want to host your blog on GitHub Pages or Cloudflare Pages.

## 👎 You might not want to use this if:
- You want a full-fledged static site generator.
- You want to use a CMS.
- You have dynamic contents.
- You need advanced UI.

## 🔑 Core Concepts
- **Minimal Config**: rawdog-md is designed to be used with minimal configuration. Choose a starter template and immediately start writing your blog posts in markdown.
- **Opiniated**: This tool is not meant to be a general purpose static site generator. It is meant to be used for a very specific use case, which is a static blog with a few simple templated pages. If you want something more general purpose and more advanced, check out Hugo, Frontmatter, or Jekyll.
- **Bring your own styling**: rawdog-md provides the minimal styling. However, it is very easy to adjust as it is just a plain CSS and HTML file.
- **Cross platform**: You can use it on Windows, Linux and Mac. (Need testing on Mac and Linux)


## 💻 Installation
rawdog-md is available on Windows, Linux and Mac. You can install it via package managers, manually, or build it from source.

### Windows
<details>
<summary>Install via Scoop</summary>

1. Install [Scoop](https://scoop.sh/)
2. Install rawdog-md
    ```shell
    scoop install https://raw.githubusercontent.com/dwiandhikaap/rawdog-md/main/package/rawdog-md.json
    ```
3. Done! You can now use it as `rawd` command in your terminal.
</details>
</br>
<details>
<summary>Install manually</summary>

1. Go to the latest [release page](https://github.com/dwiandhikaap/rawdog-md/releases/latest)
2. Download the one with `rawd-{version}-windows-amd64.zip` filename
3. Extract the zip file anywhere you like
4. Add the extracted folder to your PATH. If you don't know how, check out [this guide](https://www.architectryan.com/2018/03/17/add-to-the-path-on-windows-10/)
5. Done! You can now use it as `rawd` command in your terminal.

</details>

### Linux

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

### Mac
> 🚧 TODO 🚧

### Build from source

1. Install [Go](https://golang.org/doc/install)
2. Install rawdog-md
    ```shell 
    go install github.com/dwiandhikaap/rawdog-md
    ```
    ⚠ This will install the binary as `rawdog-md` instead of `rawd`
    if you want to change it to `rawd`, you can rename the binary file in your Go bin directory. 
    
    See this [reference](https://go.dev/ref/mod#go-install) for more information about `go install`.
3. Done!


## 🚀 Usage
> 🚧 TODO 🚧


## 📝 TODO:
- [ ] Add installation on readme
- [ ] Add usage on readme
- [ ] Make dev port configurable
- [ ] Add option to parse codeblocks deeper so we don't need JS to highlight code and use CSS only

## ⚠ Known Issues
- If you use VSCode as your editor and you're using auto save, set your auto save delay longer than rebuild duration. Otherwise, you might get an empty page because VSCode touches the file twice. See [this issue](https://github.com/microsoft/vscode/issues/9419).

## 🤝 Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. 
