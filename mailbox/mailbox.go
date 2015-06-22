package mailbox

import "github.com/PuerkitoBio/goquery"
import "fmt"
import "log"
import "regexp"
import "strings"

var mailboxBaseUrl = "http://www.yopmail.com/en/inbox.php?login=%v&p=%v&d=&ctrl=&scrl=&spam=true&v=2.6&r_c=&id="
var mailBaseUrl = "http://www.yopmail.com/mail.php?b=%v&id=%v"
var mailPerPage = 15

type Mailbox struct {
	mail string
}

func NewMailbox(mail string) *Mailbox {
	return &Mailbox{
		mail: mail,
	}
}

func (m *Mailbox) GetMails(limit int) []*Mail {
	var mails []*Mail

	for counter := 1; counter <= int(limit/mailPerPage)+1; counter++ {

		doc, err := goquery.NewDocument(fmt.Sprintf(mailboxBaseUrl, m.mail, counter))
		if err != nil {
			log.Fatal(err)
		}

		doc.Find("div.um").Each(func(i int, s *goquery.Selection) {

			id := func(s *goquery.Selection) string {
				re := regexp.MustCompile("mail.php.b=.*?id=(.*)")

				idUrl, _ := s.Find("a.lm").Attr("href")

				matches := re.FindStringSubmatch(idUrl)

				if len(matches) == 2 {
					return matches[1]
				}

				return ""
			}(s)

			if id != "" {
				mail := &Mail{
					Id:    id,
					Title: s.Find("span.lmf").Text(),
					SumUp: s.Find("span.lms").Text(),
				}

				mails = append(mails, mail)
			}
		})
	}

	if limit >= len(mails) {
		return mails
	}

	return mails[:limit]
}

func (m *Mailbox) GetMail(id string) *Mail {

	doc, err := goquery.NewDocument(fmt.Sprintf(mailBaseUrl, m.mail, id))
	if err != nil {
		log.Fatal(err)
	}

	var mail *Mail

	doc.Find("body").Each(func(i int, s *goquery.Selection) {

		fromString, fromMail := func(s *goquery.Selection) (string, string) {

			re := regexp.MustCompile(".*?: (.*?)<(.*?)>")

			matches := re.FindStringSubmatch(s.Find("div#mailhaut div:nth-child(2)").Text())

			if len(matches) == 3 {
				return matches[1], matches[2]
			}

			return "", ""
		}(s)

		mail = &Mail{
			Id:         id,
			FromString: fromString,
			FromMail:   fromMail,
			Body:       strings.TrimSpace(s.Find("div#mailmillieu").Text()),
			Title:      strings.TrimSpace(s.Find("div#mailhaut .f16").Text()),
		}
	})

	return mail
}
