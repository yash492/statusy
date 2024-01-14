package helpers

import "fmt"

func SlackHyperlinkFormat(link, name string) string {
	return fmt.Sprintf("<%v|%v>", link, name)
}

func MarkdownHyperLinkFormat(name, link string) string {
	return fmt.Sprintf("[%v](%v)", name, link)
}
