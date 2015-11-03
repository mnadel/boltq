package boltq

import (
	"fmt"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type UpdateStatement struct {
	Fields     map[string]interface{}
	BucketPath []string
}

func (p *Parser) ParseUpdate() (*UpdateStatement, error) {
	stmt := &UpdateStatement{
		Fields:     make(map[string]interface{}),
		BucketPath: make([]string, 0),
	}

	if tok, lit := p.scanNextNonWhitespaceToken(); tok != UPDATE {
		return nil, fmt.Errorf("found %q, expected update", lit)
	}

	for {
		tok, lit := p.scanNextNonWhitespaceToken()
		if tok == IDENT || tok == BUCKETSEP {
			stmt.BucketPath = append(stmt.BucketPath, lit)
		} else {
			return nil, fmt.Errorf("found %q, expected bucket path", lit)
		}

		if tok, _ := p.scanNextNonWhitespaceToken(); tok != IDENT && tok != BUCKETSEP {
			p.unscan()
			break
		}
	}

	log.Debugf("bucket = %s", strings.Join(stmt.BucketPath, "/"))

	if tok, lit := p.scanNextNonWhitespaceToken(); tok != SET {
		return nil, fmt.Errorf("found %q, expected set", lit)
	}

	for {
		tok, key := p.scanNextNonWhitespaceToken()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected identifier", key)
		}

		if tok, lit := p.scanNextNonWhitespaceToken(); tok != EQUALS {
			return nil, fmt.Errorf("found %q, expected =", lit)
		}

		tok, val := p.scanNextNonWhitespaceToken()
		if tok == QUOTE {
			t2, v2 := p.scanNextNonWhitespaceToken()
			if t2 != IDENT {
				return nil, fmt.Errorf("found %q, expected constant", v2)
			}

			val = v2

			if tok, lit := p.scanNextNonWhitespaceToken(); tok != QUOTE {
				return nil, fmt.Errorf("found %q, expected '", lit)
			}
		}

		log.Debugf("%v -> %v", key, val)

		if i, err := strconv.ParseInt(val, 10, 64); err != nil {
			if f, err2 := strconv.ParseFloat(val, 64); err2 != nil {
				log.Debugf("%v -string-> %v", key, val)
				stmt.Fields[key] = val
			} else {
				log.Debugf("%v -float-> %v", key, val)
				stmt.Fields[key] = f
			}
		} else {
			log.Debugf("%v -int-> %v", key, val)
			stmt.Fields[key] = i
		}

		if tok, _ := p.scanNextNonWhitespaceToken(); tok != COMMA {
			p.unscan()
			break
		}
	}

	return stmt, nil
}
