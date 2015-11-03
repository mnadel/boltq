package boltq

import (
	"fmt"
)

type SelectStatement struct {
	Fields     []string
	BucketPath []string
}

func (p *Parser) ParseSelect() (*SelectStatement, error) {
	stmt := &SelectStatement{
		Fields:     make([]string, 0),
		BucketPath: make([]string, 0),
	}

	if tok, lit := p.scanNextNonWhitespaceToken(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected select", lit)
	}

	for {
		tok, lit := p.scanNextNonWhitespaceToken()
		if tok == IDENT || tok == ASTERISK {
			stmt.Fields = append(stmt.Fields, lit)
		} else {
			return nil, fmt.Errorf("found %q, expected field name or *", lit)
		}

		if tok, _ := p.scanNextNonWhitespaceToken(); tok != COMMA {
			p.unscan()
			break
		}
	}

	if tok, lit := p.scanNextNonWhitespaceToken(); tok != FROM {
		return nil, fmt.Errorf("found %q, expected from", lit)
	}

	for {
		tok, lit := p.scanNextNonWhitespaceToken()
		if tok == IDENT || tok == BUCKETSEP {
			stmt.BucketPath = append(stmt.BucketPath, lit)
		} else {
			return nil, fmt.Errorf("found %q, expected bucket path", lit)
		}

		if tok, _ := p.scanNextNonWhitespaceToken(); tok != BUCKETSEP {
			p.unscan()
			break
		}
	}

	return stmt, nil
}
