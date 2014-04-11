gowiki
======

Wiki built in golang with a front page and directory. Uses Markdown. 


# What it is
I took the wiki example: [golang.org](http://golang.org/doc/articles/wiki/) and added markdown from this library: [github.com/knieriem/markdown](http://github.com/knieriem/markdown)

# Dependency
> github.com/knieriem/markdown

#Usage
* Download the repo
* Install the prerequisite
 
  > go get http://github.com/knieriem/markdown
* From the wiki directory, 
 
 > go run wiki.go

The default port is 8080


#TODO
* The search does not work at all now. It's just a pretty placeholder
* Make a master template to hold the header and CSS link
* Create reusable tools for the wiki page such as an automatic legend like wikimedia has. Basically, it takes the headers and constructs a hierarchical legend. It is kind of cool.
