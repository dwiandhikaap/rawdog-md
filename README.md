<div align="center">
<a href="https://github.com/siveroo/rawdog-md">
    <img src="https://github.com/siveroo/rawdog-md/blob/main/logo.png" alt="Logo" height="128">
</a>

<p align="center">
"Less is more"
</p>
</div>


# rawdog-md

Caveman-level simple, Markdown-to-HTML static blog tooling with minimal configuration. 

Inspired by [motherfuckingwebsite.com](https://motherfuckingwebsite.com/), uses [Go Template](https://pkg.go.dev/text/template) and HTML under the hood.

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
- **Minimal Config**: rawdog-md is designed to be used with minimal configuration. Choose a template and immediately start writing your blog posts in markdown.
- **Opiniated**: This tool is not meant to be a general purpose static site generator. It is meant to be used for a very specific use case, which is a static blog with a few simple templated pages. If you want something more general purpose and more advanced, check out Hugo, Frontmatter, or Jekyll.
- **Bring your own styling**: rawdog-md provides the minimal styling. However, it is very easy to adjust as it is just a CSS file.
- **Cross platform**: You can use it on Windows, Linux and Mac. (Needs testing on Mac and Linux)
- **Text editor agnostic**: You can use any text editor you want, as long as it can save files to disk.


## 💻 Installation
rawdog-md is available on Windows, Linux and Mac. You can install it via package managers, manually, or build it from source.
> 🚧 TODO 🚧

### Windows

> 🚧 TODO 🚧

### Linux
> 🚧 TODO 🚧

### Mac
> 🚧 TODO 🚧

### Build from source
> 🚧 TODO 🚧

## 🚀 Usage
> 🚧 TODO 🚧


## 📝 TODO:
- [ ] Add installation instructions
- [ ] Add usage instructions
- [ ] Make dev port configurable
- [ ] Add option to parse codeblocks deeper so we don't need JS to highlight code and use CSS only

## ⚠ Known Issues
- If you use VSCode as your editor and you're using auto save, set your auto save delay longer than rebuild duration. Otherwise, you might get an empty page because VSCode touches the file twice. See [this issue](https://github.com/microsoft/vscode/issues/9419).

## 🤝 Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change. 
