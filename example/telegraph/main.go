package main

import (
	"log"

	"gitlab.com/toby3d/telegraph"
)

// Content in a string format (for this example).
// Be sure to wrap every media in a <figure> tag, okay? Be easy.
const data = `
    <figure>
        <img src="/file/6a5b15e7eb4d7329ca7af.jpg"/>
    </figure>
    <p><i>Hello</i>, my name is <b>Page</b>, <u>look at me</u>!</p>
    <figure>
        <iframe src="https://youtu.be/fzQ6gRAEoy0"></iframe>
        <figcaption>
            Yes, you can embed youtube, vimeo and twitter widgets too!
        </figcaption>
    </figure>
`

var (
	account *telegraph.Account
	page    *telegraph.Page
	content []telegraph.Node
)

func errCheck(err error) {
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {
	var err error
	// Create new Telegraph account.
	requisites := telegraph.Account{
		ShortName: "techcats", // required

		// Author name/link can be epmty. So secure. Much anonymously. Wow.
		AuthorName: "dimension",                       // optional
		AuthorURL:  "https://www.yuque.com/dimension", // optional
	}
	account, err = telegraph.CreateAccount(requisites)
	errCheck(err)

	// Make sure that you have saved acc.AuthToken for create new pages or make
	// other actions by this account in next time!

	// Format content to []telegraph.Node array. Input data can be string, []byte
	// or io.Reader.
	content, err = telegraph.ContentFormat(data)
	errCheck(err)

	// Boom!.. And your text will be understandable for Telegraph. MAGIC.

	// Create new Telegraph page
	pageData := telegraph.Page{
		Title:   "My super-awesome page", // required
		Content: content,                 // required

		// Not necessarily, but, hey, it's just an example.
		AuthorName: account.AuthorName, // optional
		AuthorURL:  account.AuthorURL,  // optional
	}
	page, err = account.CreatePage(pageData, false)
	errCheck(err)

	// Show link from response on created page.
	log.Println("Kaboom! Page created, look what happened:", page.URL)
}
