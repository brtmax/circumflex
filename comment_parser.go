package main

import (
	text "github.com/MichaelMure/go-term-text"
	"regexp"
	"strconv"
	"strings"

	terminal "github.com/wayneashleyberry/terminal-dimensions"
)

// Comments represent the JSON structure as
// retrieved from cheeaun's unofficial HN API
type Comments struct {
	Author        string      `json:"user"`
	Title         string      `json:"title"`
	Comment       string      `json:"content"`
	CommentsCount int         `json:"comments_count"`
	Time          string      `json:"time_ago"`
	Points        int         `json:"points"`
	URL           string      `json:"url"`
	Domain        string      `json:"domain"`
	ID            int         `json:"id"`
	Replies       []*Comments `json:"comments"`
}

func printCommentTree(comments Comments, indentSize int, commentWith int) string {
	header := getHeader(comments, commentWith)
	originalPoster := comments.Author
	commentTree := ""
	for _, reply := range comments.Replies {
		commentTree += prettyPrintComments(*reply, 0, indentSize, commentWith, originalPoster, "")
	}
	return header + commentTree
}

func getHeader(c Comments, commentWidth int) string {
	headline := getHeadline(c.Title, c.Domain, c.URL, c.ID, commentWidth)
	infoLine := getInfoLine(c.Points, c.Author, c.Time, c.CommentsCount)
	submissionComment := parseRootComment(c.Comment, commentWidth)
	separator := getSeparator(commentWidth)
	return headline + infoLine + submissionComment + separator + DoubleNewLine
}

func getInfoLine(points int, author string, timeAgo string, numberOfComments int) string {
	p := strconv.Itoa(points)
	c := strconv.Itoa(numberOfComments)
	return dimmed(p+" points by "+author+" "+timeAgo+" • "+c+" comments") + NewLine
}

func getSeparator(commentWidth int) string {
	separator := ""
	for i := 0; i < commentWidth; i++ {
		separator += "-"
	}
	return separator
}

func getHeadline(title, domain, URL string, id, commentWidth int) string {
	headline := title + " " + paren(domain) + NewLine
	wrappedHeadline, _ := text.Wrap(headline, commentWidth)
	hyperlink := getHyperlink(domain, URL, id)

	wrappedHeadline = strings.ReplaceAll(wrappedHeadline, domain, hyperlink)

	return wrappedHeadline
}

func getHyperlink(domain string, URL string, id int) string {
	if domain != "" {
		return getHyperlinkText(URL, domain)
	}
	linkToComments := "https://news.ycombinator.com/item?id=" + strconv.Itoa(id)
	linkText := "item?id=" + strconv.Itoa(id)
	return getHyperlinkText(linkToComments, linkText)
}

func getHyperlinkText(URL string, text string) string {
	return Link1 + URL + Link2 + text + Link3
}

func parseRootComment(c string, commentWidth int) string {
	if c == "" {
		return ""
	}

	comment, URLs := parseComment(c)
	wrappedComment, _ := text.Wrap(comment, commentWidth)
	wrappedComment = applyURLs(wrappedComment, URLs)

	return wrappedComment + NewLine
}

func prettyPrintComments(c Comments, level int, indentSize int, commentWidth int, originalPoster string, parentPoster string) string {
	comment, URLs := parseComment(c.Comment)
	adjustedCommentWidth := getAdjustedCommentWidth(level, indentSize, commentWidth)

	indentBlock := getIndentBlock(level, indentSize)
	paddingWithBlock := text.WrapPad(indentBlock)
	wrappedAndPaddedComment, _ := text.Wrap(comment, adjustedCommentWidth, paddingWithBlock)

	paddingWithNoBlock := text.WrapPad(getIndentBlockWithoutBar(level, indentSize))

	author := getCommentHeading(c, level, commentWidth, originalPoster, parentPoster)
	paddedAuthor, _ := text.Wrap(author, adjustedCommentWidth, paddingWithNoBlock)
	fullComment := paddedAuthor + wrappedAndPaddedComment + DoubleNewLine
	fullComment = applyURLs(fullComment, URLs)

	if level == 0 {
		parentPoster = c.Author
	}

	for _, s := range c.Replies {
		fullComment += prettyPrintComments(*s, level+1, indentSize, commentWidth, originalPoster, parentPoster)
	}
	return fullComment
}

func getCommentHeading(c Comments, level int, commentWidth int, originalPoster string, parentPoster string) string {
	author := labelAuthor(c.Author, originalPoster, parentPoster) + " "
	numberOfReplies := 0
	headerLine := ""
	timeAgo := dimmed(c.Time)
	replies := getReplies(level, getNumberOfReplies(c, &numberOfReplies))
	if level == 0 {
		timeAgo = underline(timeAgo)
		replies = underline(replies)
		headerLine = getUnderlineString(author, timeAgo, replies, commentWidth)
	}

	return author + timeAgo + replies + headerLine + NewLine
}

func getUnderlineString(author string, timeAgo string, replies string, commentWidth int) string {
	lengthOfUnderline := commentWidth - text.Len(author) - text.Len(timeAgo) - text.Len(replies)

	headerLine := ""

	for i := 0; i < lengthOfUnderline; i++ {
		headerLine += " "
	}

	return dimmed(underline(headerLine))
}

func getNumberOfReplies(comments Comments, repliesSoFar *int) int {
	for _, reply := range comments.Replies {
		*repliesSoFar++
		getNumberOfReplies(*reply, repliesSoFar)
	}
	return *repliesSoFar
}

func applyURLs(comment string, URLs []string) string {
	for _, URL := range URLs {
		truncatedURL := truncateURL(URL)
		URLWithHyperlinkCode := getHyperlinkText(URL, truncatedURL)
		comment = strings.ReplaceAll(comment, truncatedURL, URLWithHyperlinkCode)
	}
	return comment
}

func truncateURL(URL string) string {
	if len(URL) < 60 {
		return URL
	}

	truncatedURL := ""
	for i, c := range URL {
		if i == 60 {
			truncatedURL += "..."
			break
		}
		truncatedURL += string(c)
	}
	return truncatedURL
}

func getReplies(level int, replies int) string {
	numberOfReplies := ""

	if level == 0 {
		if replies > 1 {
			r := strconv.Itoa(replies)
			numberOfReplies = " " + r + " ⤶"
		}
		return underline(dimmed(" ::" + numberOfReplies))
	}
	return ""
}

// Adjusted comment width shortens the commentWidth if the available screen size
// is smaller than the size of the commentWidth
func getAdjustedCommentWidth(level int, indentSize int, commentWidth int) int {
	x, _ := terminal.Width()
	screenWidth := int(x)

	currentIndentSize := indentSize * level
	usableScreenSize := screenWidth - currentIndentSize

	if commentWidth == 0 {
		return max(usableScreenSize, 40)
	}
	if usableScreenSize < commentWidth {
		return usableScreenSize
	}

	return commentWidth + indentSize*level
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}
func labelAuthor(author, originalPoster, parentPoster string) string {
	authorInBold := bold(author)

	switch author {
	case "dang":
		return authorInBold + green(" mod")
	case "sctb":
		return authorInBold + green(" mod")
	case originalPoster:
		return authorInBold + red(" OP")
	case parentPoster:
		return authorInBold + purple(" PP")
	default:
		return authorInBold
	}
}

func getIndentBlockWithoutBar(level int, indentSize int) string {
	if level == 0 {
		return ""
	}
	indentation := " "
	for i := 0; i < indentSize*level; i++ {
		indentation += " "
	}
	return indentation
}

func getIndentBlock(level int, indentSize int) string {
	if level == 0 {
		return ""
	}
	indentation := Normal + getColoredIndentBlock(level) + "▎" + Normal
	for i := 0; i < indentSize*level; i++ {
		indentation = " " + indentation
	}
	return indentation
}

func parseComment(comment string) (string, []string) {
	comment = replaceCharacters(comment)
	comment = replaceHTML(comment)
	comment = colorizeLinkNumbers(comment)
	URLs := extractURLs(comment)
	comment = trimURLs(comment)
	return comment, URLs
}

func replaceCharacters(input string) string {
	input = strings.ReplaceAll(input, "&#x27;", "'")
	input = strings.ReplaceAll(input, "&gt;", ">")
	input = strings.ReplaceAll(input, "&lt;", "<")
	input = strings.ReplaceAll(input, "&#x2F;", "/")
	input = strings.ReplaceAll(input, "&quot;", "\"")
	input = strings.ReplaceAll(input, "&amp;", "&")
	return input
}

func replaceHTML(input string) string {
	input = strings.Replace(input, "<p>", "", 1)

	input = strings.ReplaceAll(input, "<p>", DoubleNewLine)
	input = strings.ReplaceAll(input, "<i>", Italic)
	input = strings.ReplaceAll(input, "</i>", Normal)
	input = strings.ReplaceAll(input, "</a>", "")
	input = strings.ReplaceAll(input, "<pre><code>", Dimmed)
	input = strings.ReplaceAll(input, "</code></pre>", Normal)
	return input
}

func colorizeLinkNumbers(input string) string {
	input = strings.ReplaceAll(input, "[0]", "["+white("0")+"]")
	input = strings.ReplaceAll(input, "[1]", "["+red("1")+"]")
	input = strings.ReplaceAll(input, "[2]", "["+yellow("2")+"]")
	input = strings.ReplaceAll(input, "[3]", "["+green("3")+"]")
	input = strings.ReplaceAll(input, "[4]", "["+blue("4")+"]")
	input = strings.ReplaceAll(input, "[5]", "["+teal("5")+"]")
	input = strings.ReplaceAll(input, "[6]", "["+purple("6")+"]")
	input = strings.ReplaceAll(input, "[7]", "["+white("7")+"]")
	input = strings.ReplaceAll(input, "[8]", "["+red("8")+"]")
	input = strings.ReplaceAll(input, "[9]", "["+yellow("9")+"]")
	input = strings.ReplaceAll(input, "[10]", "["+green("10")+"]")
	return input
}

func extractURLs(input string) []string {
	expForFirstTag := regexp.MustCompile(`<a href=".*?" rel="nofollow">`)
	URLs := expForFirstTag.FindAllString(input, 10)

	for i := range URLs {
		URLs[i] = strings.ReplaceAll(URLs[i], `<a href="`, "")
		URLs[i] = strings.ReplaceAll(URLs[i], `" rel="nofollow">`, "")
	}

	return URLs
}

func trimURLs(comment string) string {
	expression := regexp.MustCompile(`<a href=".*?" rel="nofollow">`)
	return expression.ReplaceAllString(comment, "")
}
