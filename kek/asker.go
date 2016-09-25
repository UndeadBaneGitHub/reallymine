// 24 september 2016
package kek

import (
	"fmt"
	"os"
	"encoding/hex"

	"github.com/hashicorp/vault/helper/password"
)

type Asker struct {
	cmdstr	string
	kek		[]byte
	count	uint
	err		error
}

func NewAsker(cmdstr string) *Asker {
	return &Asker{
		cmdstr:	cmdstr,
	}
}

const (
	noteNeedsPassword = "You need the WD password to decrypt this drive."
	notePasswordWrong = "Wrong WD password."
)

func (a *Asker) realAsk(note string) bool {
	if note != "" {
		fmt.Printf("%s\n", note)
	}
	fmt.Print("Enter WD password: ")
	pw, err := password.Read(os.Stdin)
	if err != nil {		// including cancelled
		a.err = err
		return false
	}
	a.kek = FromPassword(pw)
	return true
}

// TODO clean this up somehow
func (a *Asker) Ask() bool {
	defer func() {
		a.count++
	}()
	switch a.cmdstr {
	case "-real":
		switch a.count {
		case 0:		// first time, return default
			a.kek = Default
			return true
		case 1:		// second time, say that one is needed
			return a.realAsk(noteNeedsPassword)
		}
		// all other times, note password is wrong
		return a.realAsk(notePasswordWrong)
	case "-askonce":
		// only ask once, then return no more
		// note not needed since we explicitly asked
		if a.count != 0 {
			return false
		}
		return a.realAsk("")
	case "-onlyask":
		if a.count == 0 {
			return a.realAsk("")
		}
		return a.realAsk(notePasswordWrong)
	case "-default":
		if a.count != 0 {
			return false
		}
		a.kek = Default
		return true
	}
	// otherwise treat as a hex string
	if a.count != 0 {
		return false
	}
	a.kek, a.err = hex.DecodeString(a.cmdstr)
	return a.err == nil
}

func (a *Asker) KEK() []byte {
	return a.kek
}

func (a *Asker) Err() error {
	return a.err
}

const AskerDescription = "" +
	"This specifies a KEK to decrypt the key sector with. " +
	"This argument can be one of the following strings:\n" +
	"- -real to use the default KEK and then ask for a password until the correct one is used, just like the main decrypt command\n" +
	"- -askonce to ask for a password once and only use the resultant KEK\n" +
	"- -onlyask to only ask for passwords until the correct one is used\n" +
	"- -default to only use the default KEK\n" +
	"Any other string is taken as a hexadecimal string to use as the KEK."
