package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("parse", func() {
	Specify("empty page returns nil elements", func() {
		es, err := parse(nil)
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(BeEmpty())
		es, err = parse([]byte{})
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(BeEmpty())
	})
	Specify("page without tags returns nil elements", func() {
		es, err := parse([]byte(" as sadas dasd as "))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(BeEmpty())
	})
	Specify("comments are ignored", func() {
		es, err := parse([]byte("<!-- asdsd a -->"))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(BeEmpty())
	})
	Specify("doctype tag is counted", func() {
		es, err := parse([]byte("<!DOCTYPE html>\n"))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(Equal([]element{{TagName: "!doctype", Count: 1}}))
	})
	Specify("simple page elements", func() {
		es, err := parse([]byte(`<!DOCTYPE html>
<html>
<head>
<title>Page Title</title>
</head>
<BODY>
<h1>This is a Heading</h1>
<p>This is a paragraph.</p>
<P>This is a paragraph 2.</P>
asd<br/>
<ul>
	<li> 1
	<li> 2
	<li> 3 </li>
</ul>
<script>alert(1)</script>
</BODY>
</html>
		`))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(ConsistOf(
			element{TagName: "!doctype", Count: 1},
			element{TagName: "html", Count: 1},
			element{TagName: "head", Count: 1},
			element{TagName: "title", Count: 1},
			element{TagName: "body", Count: 1},
			element{TagName: "h1", Count: 1},
			element{TagName: "p", Count: 2},
			element{TagName: "br", Count: 1},
			element{TagName: "ul", Count: 1},
			element{TagName: "li", Count: 3},
			element{TagName: "script", Count: 1},
		))
	})
	Specify("invalid page results without errors", func() {
		es, err := parse([]byte("<>"))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(BeEmpty())
		es, err = parse([]byte("<<<"))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(BeEmpty())
		es, err = parse([]byte("<<< <html>"))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(ConsistOf(element{TagName: "html", Count: 1}))
		es, err = parse([]byte("<html>>"))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(ConsistOf(element{TagName: "html", Count: 1}))
		es, err = parse([]byte("<html>"))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(ConsistOf(element{TagName: "html", Count: 1}))
		es, err = parse([]byte("<html/>"))
		Expect(err).ToNot(HaveOccurred())
		Expect(es).To(ConsistOf(element{TagName: "html", Count: 1}))
	})
})
