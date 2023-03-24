package tokenizer

import (
	"fmt"
	"os"
	"testing"
)

// Check if very basic html example can be tokenized without errors
func TestTokenizeWithBasicHtmlExample(t *testing.T) {
	expectedAmountOfIssues := 0
	expectedAmountOfTags := 17

	tags, issues := Tokenize(`<!DOCTYPE html>
<html>
	<body>
		<h1>The Heading</h1>
		<p name=1 hidden>The paragraph.</p>
		<a href="https://motherfuckingwebsite.com">The link</a>
		<a hidden href='https://motherfuckingwebsite.com'>The link</a>
	</body>
</html>
`) // Very basic html example

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}

	if len(tags) != expectedAmountOfTags {
		t.Fatal("wrong amount of Tags returned from Tokenize: expected", expectedAmountOfTags, ", got", len(tags), ", contents:", tags) // Printing contents as well, so it's easier to see where the node might come from
	}
}

// Check  whether offline-backup of "motherfuckingwebsite.com" can be tokenized without errors
func TestTokenizeWithMotherfuckingwebsiteExample(t *testing.T) {
	expectedAmountOfIssues := 0

	htmlCode, err := os.ReadFile("htmlExamples/motherfuckingwebsite-com.html")
	if err != nil {
		fmt.Print(err)
	}

	_, issues := Tokenize(string(htmlCode))

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check  whether offline-backup of "bettermotherfuckingwebsite.com" can be tokenized without errors
func TestTokenizeWithBettermotherfuckingwebsiteExample(t *testing.T) {
	expectedAmountOfIssues := 0

	htmlCode, err := os.ReadFile("htmlExamples/bettermotherfuckingwebsite-com.html")
	if err != nil {
		fmt.Print(err)
	}

	_, issues := Tokenize(string(htmlCode))

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check whether offline-backup of "google.com" can be tokenized without errors
func TestTokenizeWithGoogleExample(t *testing.T) {
	expectedAmountOfIssues := 0

	htmlCode, err := os.ReadFile("htmlExamples/google-com.html")
	if err != nil {
		fmt.Print(err)
	}

	_, issues := Tokenize(string(htmlCode))

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if error of double close tag (`</h/2 tag>`) is detected
func TestTokenizeWithDoubleCloseTagType(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
	<body>
		<h1>The Heading</h1>
		<p>The paragraph.</p>
		</h/2 tag>
	</body>
</html>
`)

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if error when tag is unexpectetly closed (`<a>>`)
func TestTokenizeWithUnexpectedTagCloseInStringdataMode(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
	<body>
		<h1>The Heading</h1>
		<p>The paragraph.</p>
		<a>>
	</body>
</html>
`)

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if error when tag is unexpectetly closed (`<a href=>`)
func TestTokenizeWithUnexpectedTagCloseInAttributeValueMode(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
	<body>
		<h1>The Heading</h1>
		<p>The paragraph.</p>
		<a href=>
	</body>
</html>
`)

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if issue is detected if there is trailing stringdata
func TestTokenizeWithStillInStringdataMode(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
	<body>
		<h1>The Heading</h1>
		<p> The paragraph.</p>
	</body>
</html>
ERROR
`) // trailing stringdata
	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if issue is detected if any intermediate tag wasn't closed
func TestTokenizeWithIntermediateTagStillOpen(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
	<body>
		<h1>The Heading</h1>
		<p The paragraph.</p>
	</body>
</html>
`) // p tag not closed

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if error is detected if the last tag wasn't closed
func TestTokenizeWithLastTagStillOpenInTagtypeMode(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
	<body>
		<h1>The Heading</h1>
		<p>The paragraph.</p>
	</body>
</html`) // Last tag not closed

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if error is detected if the last tag wasn't closed
func TestTokenizeWithLastTagStillOpenInAttributeKeyMode1(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
	<body>
		<h1>The Heading</h1>
		<p>The paragraph.</p>
	</body>
</html
`) // Last tag not closed

	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if issue is detected if last attributeKey wasn't finished / tag wasn't closed
func TestTokenizeWithLastTagStillOpenInAttributeKeyMode2(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
	<body>
		<h1>The Heading</h1>
		<p> The paragraph.</p>
	</body>
</html l
`) // last tag not closed (attributeKey mode)
	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if issue is detected if style tag wasn't ever closed
func TestTokenizeWithStillInStyledataMode(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
  <style>
	<body>
		<h1>The Heading</h1>
		<p> The paragraph.</p>
	</body>
</html>
`) // style node opened but not closed
	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if issue is detected if script tag wasn't ever closed
func TestTokenizeWithStillInScriptdataMode(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
  <script>
	<body>
		<h1>The Heading</h1>
		<p> The paragraph.</p>
	</body>
</html>
`) // script node opened but not closed
	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}

// Check if issue is detected if script tag wasn't ever closed
func TestTokenizeWithStillInCommentdataMode(t *testing.T) {
	expectedAmountOfIssues := 1

	_, issues := Tokenize(`<!DOCTYPE html>
<html>
  <!--
	<body>
		<h1>The Heading</h1>
		<p> The paragraph.</p>
	</body>
</html>
`) // script node opened but not closed
	if len(issues) != expectedAmountOfIssues {
		t.Fatal("wrong amount of Issues returned from Tokenize: expected", expectedAmountOfIssues, ", got", len(issues), ", contents:", issues) // Printing contents as well, so it's ieasier to see where the error might come from
	}
}
