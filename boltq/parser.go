package boltq

import (
	"fmt"
	"io"
)

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token
		lit string
		n   int
	}
}

type SelectStatement struct {
	Fields     []string
	BucketPath []string
}

func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

func (p *Parser) ParseSelect() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	if tok, lit := p.scanNextNonWhitespaceToken(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected %s", lit, SELECT)
	}

	for {
		tok, lit := p.scanNextNonWhitespaceToken()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field name or %s", lit, ASTERISK)
		}
		stmt.Fields = append(stmt.Fields, lit)

		if tok, _ := p.scanNextNonWhitespaceToken(); tok != COMMA {
			p.unscan()
			break
		}
	}

	if tok, lit := p.scanNextNonWhitespaceToken(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected %s", lit, FROM)
	}

	for {
		tok, lit := p.scanNextNonWhitespaceToken()
		if tok != IDENT && tok != BUCKETSEP {
			return nil, fmt.Errorf("found %q, expected bucket name or %q", lit, BUCKETSEP)
		}
		stmt.BucketPath = append(stmt.BucketPath, lit)

		if tok, _ := p.scanNextNonWhitespaceToken(); tok != BUCKETSEP {
			p.unscan()
			break
		}
	}

	return stmt, nil
}

func (p *Parser) scanNextNonWhitespaceToken() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

func (p *Parser) scan() (tok Token, lit string) {
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	tok, lit = p.s.Scan()

	p.buf.tok, p.buf.lit = tok, lit

	return
}

func (p *Parser) unscan() {
	p.buf.n = 1
}
