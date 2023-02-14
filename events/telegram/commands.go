package telegram

import (
	"errors"
	"log"
	"net/url"
	"read_adviser_bot/lib/e"
	"read_adviser_bot/storage"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p Processor) DoCmd(text string, chatId int, userName string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command '%s' from '%s", text, userName)

	if isAddCmd(text) {
		return p.savePage(chatId, text, userName)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatId, userName)
	case HelpCmd:
		return p.sendHelp(chatId)
	case StartCmd:
		return p.sendHello(chatId)
	default:
		return p.newMessageSender(chatId, msgUnknownCommand)
	}
}

func (p Processor) newMessageSender(chatID int, msg string) error {
	return p.tg.SendMessage(chatID, msg)
}

func (p Processor) savePage(chatID int, pageURL string, userName string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: save page", err) }()

	page := &storage.Page{
		URL:      pageURL,
		UserName: userName,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.newMessageSender(chatID, msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := p.newMessageSender(chatID, msgSaved); err != nil {
		return err
	}

	return nil
}

func (p Processor) sendRandom(chatID int, userName string) (err error) {
	defer func() { err = e.WrapIfErr("can't do command: can't send random", err) }()

	page, err := p.storage.PickRandom(userName)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.newMessageSender(chatID, msgNoSavedPages)
	}

	if err := p.newMessageSender(chatID, page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p Processor) sendHelp(chatID int) error {
	return p.newMessageSender(chatID, msgHelp)
}

func (p Processor) sendHello(chatID int) error {
	return p.newMessageSender(chatID, msgHello)
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
