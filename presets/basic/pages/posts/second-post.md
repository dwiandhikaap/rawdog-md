---
Title: Basic Markdown Test
Date: "2024-01-01 02:00:00"
Tags: ["markdown", "tutorial"]
Template: post
---

# Understanding the Basics of Markdown

Markdown is a lightweight markup language that allows you to format plain text. It is widely used for writing documentation, creating readme files, and blogging. In this post, we will explore the basics of Markdown, its syntax, and some advanced features.

## Table of Contents
1. [What is Markdown?](#what-is-markdown)
2. [Basic Syntax](#basic-syntax)
   - [Headings](#headings)
   - [Emphasis](#emphasis)
   - [Lists](#lists)
   - [Links](#links)
   - [Images](#images)
3. [Advanced Features](#advanced-features)
   - [Blockquotes](#blockquotes)
   - [Code Blocks](#code-blocks)
   - [Tables](#tables)
4. [Conclusion](#conclusion)

## What is Markdown?

Markdown is a text-based format that allows you to write using an easy-to-read and easy-to-write plain text format. It was created by John Gruber in 2004. The main idea behind Markdown is to make it simple to convert text into HTML.

## Basic Syntax

### Headings

You can create headings by using the `#` symbol. The number of `#` symbols indicates the level of the heading. For example:

```markdown
# This is a H1
## This is a H2
### This is a H3
#### This is a H4
##### This is a H5
###### This is a H6
```

### Emphasis

You can emphasize text using asterisks or underscores:

- *Italic* or _Italic_
- **Bold** or __Bold__
- ***Bold and Italic*** or ___Bold and Italic___

### Lists

Markdown supports both ordered and unordered lists.

#### Unordered List
Source:
```markdown
- Item 1
- Item 2
  - Subitem 1
  - Subitem 2
```
Output:
- Item 1
- Item 2
  - Subitem 1
  - Subitem 2

#### Ordered List
Source:
```markdown
1. First item
2. Second item
   1. Subitem
   2. Subitem
```
Output:
1. First item
2. Second item
   1. Subitem
   2. Subitem

### Links

You can create links using the following syntax:

Source:
```markdown
[Google](https://www.google.com)
```

Output:
[Google](https://www.google.com)

### Images

Images can be added with a similar syntax to links, but with an exclamation mark in front:

Source:
```markdown
![Alt text](/images/cat.jpg)
```

Output:
![Alt text](/images/cat.jpg)

## Advanced Features

### YouTube Embed
![](https://www.youtube.com/watch?v=dQw4w9WgXcQ)

### Blockquotes

Blockquotes are created using the `>` symbol:

Source:
```markdown
> This is a blockquote.
```

Output:
> This is a blockquote.

### Code Blocks

To create inline code, use backticks:

Source:
```markdown
`inline code`
```

Output:
`inline code`

For multi-line code blocks, use triple backticks:

Source:
~~~markdown
```
function test() {
  console.log("Hello, World!");
}
```
~~~

Output:
```js
function test() {
  console.log("Hello, World!");
}
```

### Tables

Tables can be created using pipes and dashes:

Source:
```markdown
| Header 1 | Header 2 |
|----------|----------|
| Row 1    | Row 2    |
| Row 3    | Row 4    |
```

Output:
| Header 1 | Header 2 |
|----------|----------|
| Row 1    | Row 2    |
| Row 3    | Row 4    |

## Conclusion

Markdown is a powerful tool for formatting text in a simple and efficient way. Whether you're writing documentation, creating a blog post, or just taking notes, Markdown can help you keep your text organized and readable. 

Feel free to experiment with the syntax and see what works best for you!

---

### References
- [Markdown Guide](https://www.markdownguide.org/)
- [GitHub Flavored Markdown](https://github.github.com/gfm/)
