---
title: rawdog-md Readme
tags: ["rawdog-md", "readme"]
template: post
---

# rawdog-md

Caveman-level simple, minimal Markdown-to-HTML static blog tooling with minimal configuration.

## Core Concepts
- **Minimal Config**: rawdog-md is designed to be used with minimal configuration. Choose a template and immediately start writing your blog posts in markdown.
- **Opiniated**: This tool is not meant to be a general purpose static site generator. It is meant to be used for a very specific use case, which is a static blog with a few templated pages. If you want something more general purpose and more advanced, check out Hugo, Front Matter, or Jekyll.
- **Bring your own styling (If you want)**: rawdog-md provides the bare minimum styling. However, it is very easy to adjust for your liking because it's just plain HTML and CSS.
- 🚧 TODO 🚧 ~~**Cross platform**: You can use it on Windows, Linux and Mac.~~ 
- **Text editor agnostic**: You can use any text editor you want, as long as it can save files to disk.
- **LSP Support**: Provides auto-complete that provides frontmatter properties suggestions for your markdown and handlebars files.

## Installation
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


## TODO:
- [ ] Add installation instructions
- [ ] Add usage instructions
- [ ] Create skeleton preset
- [ ] Make dev port configurable
- [ ] Enforce `date` property in frontmatter to sort posts
- [ ] Add some handlebars helpers for common tasks like comparing value, sorting, etc.

## Known Issues
- If you use VSCode as your editor and you're using auto save, set you delay longer than rebuild time. Otherwise, you might get an empty page because VSCode touches the file twice. See [this issue](https://github.com/microsoft/vscode/issues/9419).