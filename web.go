package main

import (
	"fmt"
	"io"
)

type PrintFunc func(format string, a ...interface{}) (n int, err error)

func makePrintFunc(w io.Writer) func(format string, a ...interface{}) (n int, err error) {
	// Return closure enclosing io.Writer.
	return func(format string, a ...interface{}) (n int, err error) {
		return fmt.Fprintf(w, format, a...)
	}
}

//*** HTML template functions ***
func printHtmlOpen(P PrintFunc, title string, jsurls []string) {
	P("<!DOCTYPE html>\n")
	P("<html>\n")
	P("<head>\n")
	P("<meta charset=\"utf-8\">\n")
	P("<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	P("<title>%s</title>\n", title)
	P("<link rel=\"stylesheet\" type=\"text/css\" href=\"/static/style.css\">\n")
	for _, jsurl := range jsurls {
		P("<script defer src=\"%s\"></script>\n", jsurl)
	}
	P("<style>\n")
	P(".myfont {font-family: Helvetica Neue,Helvetica,Arial,sans-serif;}\n")
	P("</style>\n")
	P("</head>\n")
	P("<body class=\"bg-page text-xs p-2\">\n")
}
func printHtmlClose(P PrintFunc) {
	P("</body>\n")
	P("</html>\n")
}
func printContainerOpen(P PrintFunc) {
	P("<div id=\"container\" class=\"\">\n")
}
func printContainerClose(P PrintFunc) {
	P("</div>\n")
}
