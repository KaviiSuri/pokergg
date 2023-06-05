package deck

import (
	"fmt"
	"math/rand"
)

type Suit int

// The Suit type is an integer type that represents the suit of a card.
const (
	Spades Suit = iota
	Hearts
	Diamonds
	Clubs
	NumSuits
)

func (s Suit) String() string {
	switch s {
	case Spades:
		return "Spades"
	case Hearts:
		return "Hearts"
	case Diamonds:
		return "Diamonds"
	case Clubs:
		return "Clubs"
	default:
		panic("invalid suit")
	}
}

// The Denom type is an integer type that represents the denomination of a card.
type Denom int

const (
	Ace Denom = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	NumCards
)

func (d Denom) String() string {
	switch d {
	case Ace:
		return "Ace"
	case Jack:
		return "Jack"
	case Queen:
		return "Queen"
	case King:
		return "King"
	default:
		return fmt.Sprintf("%d", d)
	}
}

type Card struct {
	suit  Suit
	value Denom
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %s %s", c.value, c.suit, suitToUnicode(c.suit))
}

func NewCard(s Suit, v Denom) Card {
	if s < Spades || s > Clubs {
		panic("invalid suit")
	}
	if v < Ace || v > King {
		panic("invalid denom")

	}
	return Card{
		suit:  s,
		value: v,
	}
}

type Deck [52]Card

func New() Deck {
	d := Deck{}
	nCards := int(NumCards)
	nSuits := int(NumSuits)

	x := 0
	for i := 0; i < nSuits; i++ {
		for j := 1; j < nCards; j++ {
			d[x] = NewCard(Suit(i), Denom(j))
			x++
		}
	}

	return shuffle(d)
}

func shuffle(d Deck) Deck {
	for i := 0; i < len(d); i++ {
		j := rand.Intn(len(d))
		d[i], d[j] = d[j], d[i]
	}

	return d
}

func suitToUnicode(s Suit) string {
	switch s {
	case Spades:
		return "♠"
	case Hearts:
		return "♥"
	case Diamonds:
		return "♦"
	case Clubs:
		return "♣"
	default:
		panic("invalid suit")
	}
}
