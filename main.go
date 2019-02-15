package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

/*<?xml version="1.0" encoding="utf8" standalone="yes"?>
<dictionary version="0.8" revision="403605">
    <grammemes>
        <grammeme parent="">POST</grammeme>
        <grammeme parent="POST">NOUN</grammeme>
        ...
    </grammemes>
    <lemmata>
        <lemma id="1" rev="402007">
            <l t="абажур"><g v="NOUN"/><g v="inan"/><g v="masc"/></l>
            <f t="абажур"><g v="sing"/><g v="nomn"/></f>
            <f t="абажура"><g v="sing"/><g v="gent"/></f>
            ...
        </lemma>
        ...
    </lemmata>
    <link_types>
        <type id="1">VERB_GERUND</type>
        ...
    </link_types>
    <links>
        <link id="1" from="104" to="106" type="1"/>
        ...
    </links>
</dictionary>`
*/

type Link struct {
	Id   int64 `xml:"id,attr"`
	From int64 `xml:"from,attr"`
	To   int64 `xml:"to,attr"`
	Type int64 `xml:"type,attr"`
}

type Lemma struct {
	Id  int64    `xml:"id,attr"`
	Rev int64    `xml:"rev,attr"`
	L   LemmaL   `xml:"l"`
	F   []LemmaF `xml:"f"`
}

type LemmaL struct {
	T string   `xml:"t,attr"`
	G []LemmaG `xml:"g"`
}

type LemmaF struct {
	T string   `xml:"t,attr"`
	G []LemmaG `xml:"g"`
}

type LemmaG struct {
	V string `xml:"v,attr"`
}

var index map[int64]string = make(map[int64]string, 500000)

func main() {

	br := bufio.NewReader(os.Stdin)

	for {

		str, err := br.ReadString('\n')
		if err != nil && str == "" {
			break
		}

		if strings.Index(str, "<lemma ") != -1 {
			onLemma(str)
			continue
		}

		if strings.Index(str, "<link ") != -1 {
			onLink(str)
			continue
		}
	}

	for _, v := range index {
		fmt.Println(v)
	}
}

func onLink(str string) {

	l := new(Link)

	if e := xml.Unmarshal([]byte(str), l); e != nil {
		return
	}

	if l.Type != 3 {
		return
	}

	if from, h1 := index[l.From]; h1 {
		if to, h2 := index[l.To]; h2 {
			fmt.Println(from + " " + to)
			delete(index, l.From)
			delete(index, l.To)
		}
	}
}

func onLemma(str string) {

	l := new(Lemma)

	if e := xml.Unmarshal([]byte(str), l); e != nil {
		return
	}

	bldr := strings.Builder{}

	data := map[string]bool{}

	data[l.L.T] = true

	bldr.WriteString(l.L.T)

	for _, f := range l.F {

		if _, h := data[f.T]; !h {
			data[f.T] = true
			bldr.WriteRune(' ')
			bldr.WriteString(f.T)
		}
	}

	s := strings.ToLower(bldr.String())
	s = strings.Replace(s, "ё", "е", -1)

	index[l.Id] = s
}
